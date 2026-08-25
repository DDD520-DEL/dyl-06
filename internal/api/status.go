package api

import (
	"encoding/json"
	"net/http"

	"github.com/dyl-06/telemetry/internal/alert"
	"github.com/dyl-06/telemetry/internal/metrics"
	"github.com/dyl-06/telemetry/internal/model"
	"github.com/dyl-06/telemetry/internal/notify"
	"github.com/dyl-06/telemetry/internal/registry"
	"github.com/dyl-06/telemetry/internal/store"
)

// HandleStatus 返回服务运行状态汇总。
func HandleStatus(
	st *store.Store,
	m *metrics.Metrics,
	rules *alert.RuleSet,
	subs *notify.Registry,
	workers *registry.Registry,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pending := 0
		for _, sub := range subs.List() {
			pending += subs.Pending(sub.ID)
		}
		report := model.StatusReport{
			Workers:      workers.Count(),
			Rules:        len(rules.List()),
			Subscribers:  len(subs.List()),
			Pending:      pending,
			Ingested:     m.Snapshot()["ingested"],
			Dropped:      m.Snapshot()["dropped"],
			Rejected:     m.Snapshot()["rejected"],
			Alerts:       m.Snapshot()["alerts"],
			Events:       m.Snapshot()["events"],
			KnownSeries:  st.SeriesCount(),
			RouteVersion: len(st.Routes()),
		}
		stats := st.Stats()
		report.Buckets = stats.Buckets
		report.ClosedBuckets = stats.Closed
		_ = json.NewEncoder(w).Encode(report)
	}
}

// HandleSnapshot 导出存储快照。
func HandleSnapshot(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := st.Snapshot()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(data)
	}
}
