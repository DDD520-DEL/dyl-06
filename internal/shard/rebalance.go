package shard

// Rebalance 按新分片列表重建路由表。
func (r *Router) Rebalance(shards []string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.shards = append([]string(nil), shards...)
}
