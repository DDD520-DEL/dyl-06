package downsample

import (
	"sort"

	"github.com/dyl-06/telemetry/internal/agg"
	"github.com/dyl-06/telemetry/internal/model"
	"github.com/dyl-06/telemetry/internal/store"
)

// Rollup 对序列在 [from, to] 范围内的桶做粗粒度聚合。
func Rollup(st *store.Store, series string, from, to int64) []model.RollupResult {
	results := st.QuerySlots(series, from, to)
	bySlot := make(map[int64][]model.AggResult)
	slots := make([]int64, 0)
	for _, result := range results {
		key := result.Slot / 60
		if _, ok := bySlot[key]; !ok {
			slots = append(slots, key)
		}
		bySlot[key] = append(bySlot[key], result.Agg)
	}
	sort.Slice(slots, func(i, j int) bool { return slots[i] < slots[j] })
	out := make([]model.RollupResult, 0, len(slots))
	for _, slot := range slots {
		merged := agg.MergeResults(bySlot[slot])
		out = append(out, model.RollupResult{
			Slot:  slot,
			Count: merged.Count,
			Sum:   merged.Sum,
			Max:   merged.Max,
			Min:   merged.Min,
		})
	}
	return out
}

// Run 关闭截止桶并返回降采样结果。
func (m *Manager) Run(cutoff int64) []model.RollupResult {
	m.CloseExpired(cutoff)
	out := make([]model.RollupResult, 0)
	for _, series := range m.store.Series() {
		out = append(out, Rollup(m.store, series, cutoff, cutoff+59)...)
	}
	return out
}
