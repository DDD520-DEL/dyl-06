package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dyl-06/telemetry/internal/alert"
	"github.com/dyl-06/telemetry/internal/api"
	"github.com/dyl-06/telemetry/internal/audit"
	"github.com/dyl-06/telemetry/internal/clock"
	"github.com/dyl-06/telemetry/internal/config"
	"github.com/dyl-06/telemetry/internal/downsample"
	"github.com/dyl-06/telemetry/internal/idgen"
	"github.com/dyl-06/telemetry/internal/ingest"
	"github.com/dyl-06/telemetry/internal/metrics"
	"github.com/dyl-06/telemetry/internal/notify"
	"github.com/dyl-06/telemetry/internal/policy"
	"github.com/dyl-06/telemetry/internal/registry"
	"github.com/dyl-06/telemetry/internal/seq"
	"github.com/dyl-06/telemetry/internal/shard"
	"github.com/dyl-06/telemetry/internal/store"
)

func main() {
	cfg := config.Load()
	clk := clock.SystemClock{}
	m := metrics.New()
	auditLog := audit.New(1000)
	router := shard.New([]string{"shard-a", "shard-b"})
	st := store.New(router)
	pipeline := ingest.NewPipeline(st, seq.New(), m)
	rules := alert.NewRuleSet()
	reg := notify.NewRegistry()
	dispatcher := notify.NewDispatcher(reg, m)
	engine := alert.NewEngine(rules, dispatcher, clk, m)
	history := alert.NewHistory(200)
	engine.SetHistory(history)
	manager := downsample.New(st, cfg.BucketSize)
	reader := downsample.NewReader(st)
	workers := registry.New()
	delivery := notify.NewDeliveryWorker(reg, m)
	gen := idgen.New()
	for i := 0; i < cfg.Workers; i++ {
		ingest.RegisterWorker(workers, "worker-"+string(rune('a'+i)), clk)
	}

	go func() {
		retentionSlots := cfg.RetentionSeconds / cfg.BucketSize
		if err := (policy.RetentionPolicy{Slots: retentionSlots}).Validate(); err != nil {
			log.Printf("retention policy invalid: %v", err)
			return
		}
		for {
			now := clk.Now()
			manager.Run(now - cfg.BucketSize*10)
			workers.PruneStale(now, cfg.PruneTTL)
			delivery.Deliver()
			st.PurgeBefore(now/cfg.BucketSize - retentionSlots)
			engine.EvaluateAll(st)
			time.Sleep(30 * time.Second)
		}
	}()

	server := api.NewServer(st, router, pipeline, m, engine, workers, clk)
	mux := http.NewServeMux()
	mux.HandleFunc("POST /ingest", server.HandleIngest)
	mux.HandleFunc("GET /query", server.HandleQuery(reader))
	mux.HandleFunc("POST /ingest/batch", api.HandleBatchIngest(st))
	mux.HandleFunc("POST /subscribe", api.HandleSubscribe(reg, gen))
	mux.HandleFunc("POST /unsubscribe", api.HandleUnsubscribe(reg))
	mux.HandleFunc("POST /rules", api.HandleRuleUpdate(rules))
	mux.HandleFunc("POST /rebalance", server.HandleRebalance)
	mux.HandleFunc("POST /heartbeat", server.HandleHeartbeat)
	mux.HandleFunc("GET /status", api.HandleStatus(st, m, rules, reg, workers))
	mux.HandleFunc("GET /snapshot", api.HandleSnapshot(st))
	mux.HandleFunc("GET /audit", api.HandleAudit(auditLog))
	mux.HandleFunc("GET /shards", api.HandleShards(st))
	mux.HandleFunc("GET /top", api.HandleTop(st))
	handler := api.LogRequests(auditLog, mux)

	log.Printf("telemetry listening on %s", cfg.Addr)
	if snapshotPath := os.Getenv("TELEMETRY_SNAPSHOT"); snapshotPath != "" {
		if data, err := os.ReadFile(snapshotPath); err == nil {
			if err := st.Restore(data); err != nil {
				log.Printf("snapshot restore failed: %v", err)
			}
		}
	}
	if err := http.ListenAndServe(cfg.Addr, handler); err != nil {
		log.Fatal(err)
	}
}
