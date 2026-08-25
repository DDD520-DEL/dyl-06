package registry

import (
	"github.com/dyl-06/telemetry/internal/model"
)

// PruneStale 移除超过 ttl 秒未心跳的 worker，返回被移除的 ID。
func (r *Registry) PruneStale(now, ttl int64) []string {
	r.mu.Lock()
	defer r.mu.Unlock()
	removed := make([]string, 0)
	for id, worker := range r.workers {
		if now-worker.LastSeen > ttl {
			delete(r.workers, id)
			removed = append(removed, id)
		}
	}
	return removed
}

// Count 返回注册的 worker 数量。
func (r *Registry) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.workers)
}

// Get 返回指定 worker。
func (r *Registry) Get(id string) (model.Worker, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	worker, ok := r.workers[id]
	return worker, ok
}
