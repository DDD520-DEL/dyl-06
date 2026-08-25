package store

import (
	"sort"
	"sync"

	"github.com/dyl-06/telemetry/internal/agg"
	"github.com/dyl-06/telemetry/internal/model"
	"github.com/dyl-06/telemetry/internal/shard"
)

// Store 持有时间桶、最新点位与路由缓存。
type Store struct {
	router  *shard.Router
	mu      sync.RWMutex
	buckets map[model.BucketKey]*agg.Bucket
	closed  map[model.BucketKey]bool
	points  map[string]model.Point
	routes  map[string]string
}

// New 创建 Store。
func New(router *shard.Router) *Store {
	s := &Store{
		router:  router,
		buckets: make(map[model.BucketKey]*agg.Bucket),
		closed:  make(map[model.BucketKey]bool),
		points:  make(map[string]model.Point),
		routes:  make(map[string]string),
	}
	s.Reindex()
	return s
}

func (s *Store) bucket(p model.Point) *agg.Bucket {
	key := model.BucketKey{Series: p.Series, Slot: p.Ts / 60}
	b, ok := s.buckets[key]
	if !ok {
		b = agg.NewBucket()
		s.buckets[key] = b
	}
	return b
}

// GetBucket 返回时间桶，不存在时返回 nil。
func (s *Store) GetBucket(series string, slot int64) *agg.Bucket {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.buckets[model.BucketKey{Series: series, Slot: slot}]
}

// CloseBucket 标记时间桶已关闭，不再接受新点。
func (s *Store) CloseBucket(series string, slot int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.closed[model.BucketKey{Series: series, Slot: slot}] = true
}

// IsClosed 查询时间桶是否已关闭。
func (s *Store) IsClosed(series string, slot int64) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.closed[model.BucketKey{Series: series, Slot: slot}]
}

// LastPoint 返回某序列最新点位。
func (s *Store) LastPoint(series string) (model.Point, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	p, ok := s.points[series]
	return p, ok
}

// Query 返回某序列时间范围内的桶聚合结果。
func (s *Store) Query(series string, from, to int64) []model.AggResult {
	results := s.QuerySlots(series, from, to)
	out := make([]model.AggResult, 0, len(results))
	for _, result := range results {
		out = append(out, result.Agg)
	}
	return out
}

// QuerySlots 返回某序列时间范围内带槽位的桶聚合结果。
func (s *Store) QuerySlots(series string, from, to int64) []model.BucketResult {
	s.mu.RLock()
	defer s.mu.RUnlock()
	keys := make([]model.BucketKey, 0)
	for k := range s.buckets {
		if k.Series == series && k.Slot >= from && k.Slot <= to {
			keys = append(keys, k)
		}
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i].Slot < keys[j].Slot })
	out := make([]model.BucketResult, 0, len(keys))
	for _, k := range keys {
		out = append(out, model.BucketResult{Slot: k.Slot, Agg: s.buckets[k].Read()})
	}
	return out
}

// BucketCount 返回时间桶总数。
func (s *Store) BucketCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.buckets)
}

// PurgeBefore 删除早于 cutoff 的所有时间桶与点位。
func (s *Store) PurgeBefore(cutoff int64) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	removed := 0
	for key := range s.buckets {
		if key.Slot < cutoff {
			delete(s.buckets, key)
			delete(s.closed, key)
			removed++
		}
	}
	for name, point := range s.points {
		if point.Ts/60 < cutoff {
			delete(s.points, name)
		}
	}
	return removed
}

// Series 返回已知序列列表。
func (s *Store) Series() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]string, 0, len(s.points))
	for name := range s.points {
		out = append(out, name)
	}
	sort.Strings(out)
	return out
}
