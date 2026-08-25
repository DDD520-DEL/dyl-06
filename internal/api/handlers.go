package api

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/dyl-06/telemetry/internal/alert"
	"github.com/dyl-06/telemetry/internal/clock"
	"github.com/dyl-06/telemetry/internal/ingest"
	"github.com/dyl-06/telemetry/internal/metrics"
	"github.com/dyl-06/telemetry/internal/model"
	"github.com/dyl-06/telemetry/internal/policy"
	"github.com/dyl-06/telemetry/internal/registry"
	"github.com/dyl-06/telemetry/internal/shard"
	"github.com/dyl-06/telemetry/internal/store"
)

// Server 聚合 HTTP 处理器依赖。
type Server struct {
	store   *store.Store
	router  *shard.Router
	ingest  *ingest.Pipeline
	metrics *metrics.Metrics
	engine  *alert.Engine
	workers *registry.Registry
	clock   clock.Clock
}

// NewServer 创建 API 服务器。
func NewServer(
	st *store.Store,
	router *shard.Router,
	pipeline *ingest.Pipeline,
	m *metrics.Metrics,
	engine *alert.Engine,
	workers *registry.Registry,
	clk clock.Clock,
) *Server {
	return &Server{
		store:   st,
		router:  router,
		ingest:  pipeline,
		metrics: m,
		engine:  engine,
		workers: workers,
		clock:   clk,
	}
}

// HandleIngest 接收单条上报。
func (s *Server) HandleIngest(w http.ResponseWriter, r *http.Request) {
	raw, err := io.ReadAll(io.LimitReader(r.Body, 64*1024))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	point, err := s.ingest.HandleRaw(raw)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if point.Series != "" {
		if bucket := s.store.GetBucket(point.Series, point.Ts/60); bucket != nil {
			s.engine.Evaluate(point.Series, bucket.Read())
		}
	}
	w.WriteHeader(http.StatusAccepted)
}

// HandleRebalance 触发分片重平衡并刷新存储路由索引。
func (s *Server) HandleRebalance(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Shards []string `json:"shards"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := RebalanceRouter(s.router, s.store, req.Shards); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// RebalanceRouter 执行分片重平衡并刷新存储路由索引。
func RebalanceRouter(r *shard.Router, st *store.Store, shards []string) error {
	if err := policy.ValidateShards(shards); err != nil {
		return err
	}
	r.Rebalance(shards)
	// 路由表已变更：旧的路由缓存按 hash%N 计算已失效，必须清空，
	// 否则既有序列仍按旧分片写入，新分片空转。
	st.Reindex()
	return nil
}

// HandleHeartbeat 接收 worker 心跳。
func (s *Server) HandleHeartbeat(w http.ResponseWriter, r *http.Request) {
	var req model.Worker
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	s.workers.Register(model.Worker{ID: req.ID, LastSeen: s.clock.Now()})
	w.WriteHeader(http.StatusOK)
}
