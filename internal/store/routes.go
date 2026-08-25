package store

// Reindex 使路由缓存失效，下次写入重新计算。
func (s *Store) Reindex() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.routes = make(map[string]string)
}

// Target 返回序列当前写入目标分片，并缓存首次计算结果。
func (s *Store) Target(series string) string {
	s.mu.RLock()
	target, ok := s.routes[series]
	s.mu.RUnlock()
	if ok {
		return target
	}
	target = s.router.Route(series)
	s.mu.Lock()
	if _, exists := s.routes[series]; !exists {
		s.routes[series] = target
	}
	s.mu.Unlock()
	return target
}

// Routes 返回路由索引副本。
func (s *Store) Routes() map[string]string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make(map[string]string, len(s.routes))
	for k, v := range s.routes {
		out[k] = v
	}
	return out
}
