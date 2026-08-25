package metrics

import "sync/atomic"

// Metrics 汇总运行时计数。
type Metrics struct {
	Ingested atomic.Int64
	Dropped  atomic.Int64
	Rejected atomic.Int64
	Alerts   atomic.Int64
	Events   atomic.Int64
}

// New 创建 Metrics。
func New() *Metrics {
	return &Metrics{}
}
