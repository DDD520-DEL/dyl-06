package ingest

import (
	"github.com/dyl-06/telemetry/internal/decode"
	"github.com/dyl-06/telemetry/internal/metrics"
	"github.com/dyl-06/telemetry/internal/model"
	"github.com/dyl-06/telemetry/internal/seq"
	"github.com/dyl-06/telemetry/internal/store"
)

// Pipeline 处理一条上报：解码、守卫、写入。
type Pipeline struct {
	store   *store.Store
	metrics *metrics.Metrics
	guard   *seq.Guard
}

// NewPipeline 创建摄入管道。
func NewPipeline(st *store.Store, guard *seq.Guard, m *metrics.Metrics) *Pipeline {
	return &Pipeline{store: st, metrics: m, guard: guard}
}

// HandleRaw 解析并写入一条原始点位；解析失败返回错误。
func (p *Pipeline) HandleRaw(raw []byte) (model.Point, error) {
	point, err := decode.Decode(raw)
	if err != nil {
		p.metrics.Rejected.Add(1)
		return model.Point{}, err
	}
	if !p.store.WritePointWithGuard(p.guard, point) {
		p.metrics.Dropped.Add(1)
		return model.Point{}, nil
	}
	p.metrics.Ingested.Add(1)
	return point, nil
}
