package store

import (
	"encoding/json"

	"github.com/dyl-06/telemetry/internal/model"
)

type snapshot struct {
	Points map[string]model.Point `json:"points"`
	Closed []model.BucketKey      `json:"closed"`
}

// Snapshot 导出存储状态。
func (s *Store) Snapshot() ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	snap := snapshot{Points: make(map[string]model.Point, len(s.points))}
	for name, point := range s.points {
		snap.Points[name] = point
	}
	for key := range s.closed {
		snap.Closed = append(snap.Closed, key)
	}
	return json.Marshal(snap)
}

// Restore 恢复存储状态。
func (s *Store) Restore(data []byte) error {
	var snap snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.points = snap.Points
	for _, key := range snap.Closed {
		s.closed[key] = true
	}
	return nil
}
