package store

import (
	"sort"
)

// TopSeries 是聚合值最高的序列。
type TopSeries struct {
	Series string  `json:"series"`
	Sum    float64 `json:"sum"`
	Count  int64   `json:"count"`
}

// QueryTopN 返回聚合总和最大的 N 个序列。
func (s *Store) QueryTopN(limit int) []TopSeries {
	s.mu.RLock()
	defer s.mu.RUnlock()
	aggregated := make(map[string]TopSeries)
	for key, bucket := range s.buckets {
		result := bucket.Read()
		entry := aggregated[key.Series]
		entry.Series = key.Series
		entry.Sum += result.Sum
		entry.Count += result.Count
		aggregated[key.Series] = entry
	}
	out := make([]TopSeries, 0, len(aggregated))
	for _, entry := range aggregated {
		out = append(out, entry)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Sum > out[j].Sum })
	if len(out) > limit {
		out = out[:limit]
	}
	return out
}
