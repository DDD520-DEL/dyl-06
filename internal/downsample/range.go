package downsample

import (
	"github.com/dyl-06/telemetry/internal/model"
)

// ReadRange 返回序列在 [from, to] 槽位内的降采样结果。
func (r *Reader) ReadRange(series string, from, to int64) []model.RollupResult {
	return Rollup(r.store, series, from, to)
}

// ReadLatest 返回序列最近一个粗粒度槽的降采样结果。
func (r *Reader) ReadLatest(series string, slot int64) (model.RollupResult, bool) {
	results := Rollup(r.store, series, slot, slot)
	if len(results) == 0 {
		return model.RollupResult{}, false
	}
	return results[len(results)-1], true
}
