package store

// Stats 汇总存储规模。
type Stats struct {
	Series  int `json:"series"`
	Buckets int `json:"buckets"`
	Closed  int `json:"closed"`
}

// Stats 返回存储规模统计。
func (s *Store) Stats() Stats {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return Stats{
		Series:  len(s.points),
		Buckets: len(s.buckets),
		Closed:  len(s.closed),
	}
}
