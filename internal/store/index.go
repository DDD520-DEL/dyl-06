package store

// SeriesCount 返回已知序列数量。
func (s *Store) SeriesCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.points)
}

// ShardFor 返回序列当前写入分片。
func (s *Store) ShardFor(series string) string {
	return s.Target(series)
}
