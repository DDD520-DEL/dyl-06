package agg

import (
	"math"
	"sync"

	"github.com/dyl-06/telemetry/internal/model"
)

// Bucket 是一个时间桶的聚合状态。
type Bucket struct {
	mu    sync.Mutex
	Count int64
	Sum   float64
	Max   float64
	Min   float64
}

// NewBucket 创建空桶。
func NewBucket() *Bucket {
	return &Bucket{Max: math.Inf(-1), Min: math.Inf(1)}
}

// Merge 把点位合并进桶，线程安全。
func (b *Bucket) Merge(p model.Point) {
	b.Count++
	b.Sum += p.Value
	if p.Value > b.Max {
		b.Max = p.Value
	}
	if p.Value < b.Min {
		b.Min = p.Value
	}
}

// Read 返回聚合结果快照。
func (b *Bucket) Read() model.AggResult {
	b.mu.Lock()
	defer b.mu.Unlock()
	return model.AggResult{Count: b.Count, Sum: b.Sum, Max: b.Max, Min: b.Min}
}
