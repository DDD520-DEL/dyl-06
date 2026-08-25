package agg

import "github.com/dyl-06/telemetry/internal/model"

// MergeResults 合并多个桶的聚合结果。
func MergeResults(results []model.AggResult) model.AggResult {
	if len(results) == 0 {
		return model.AggResult{}
	}
	out := model.AggResult{Min: results[0].Min}
	for _, result := range results {
		out.Count += result.Count
		out.Sum += result.Sum
		if result.Max > out.Max {
			out.Max = result.Max
		}
		if result.Min < out.Min {
			out.Min = result.Min
		}
	}
	return out
}
