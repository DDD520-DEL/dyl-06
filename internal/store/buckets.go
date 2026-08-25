package store

import (
	"github.com/dyl-06/telemetry/internal/model"
)

// WritePoint 写入单个点位；已关闭的时间桶被忽略。
func (s *Store) WritePoint(p model.Point) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.bucket(p).Merge(p)
	s.points[p.Series] = p
}
