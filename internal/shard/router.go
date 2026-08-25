package shard

import (
	"sync"

	"github.com/cespare/xxhash/v2"
)

// Router 负责序列到分片的哈希路由。
type Router struct {
	mu     sync.RWMutex
	shards []string
}

// New 创建路由表。
func New(shards []string) *Router {
	return &Router{shards: append([]string(nil), shards...)}
}

// Route 返回序列当前应写入的分片。
func (r *Router) Route(series string) string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if len(r.shards) == 0 {
		return ""
	}
	hash := xxhash.Sum64String(series)
	return r.shards[hash%uint64(len(r.shards))]
}

// Shards 返回当前分片列表。
func (r *Router) Shards() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return append([]string(nil), r.shards...)
}
