package registry

import (
	"sort"
	"sync"

	"github.com/dyl-06/telemetry/internal/model"
)

// Registry 维护摄取 worker 注册表。
type Registry struct {
	mu      sync.RWMutex
	workers map[string]model.Worker
}

// New 创建 worker 注册表。
func New() *Registry {
	return &Registry{workers: make(map[string]model.Worker)}
}

// Register 注册或刷新 worker。
func (r *Registry) Register(worker model.Worker) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.workers[worker.ID] = worker
}

// Heartbeat 刷新 worker 存活时间。
func (r *Registry) Heartbeat(id string, now int64) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if worker, ok := r.workers[id]; ok {
		worker.LastSeen = now
		r.workers[id] = worker
	}
}

// List 返回全部 worker。
func (r *Registry) List() []model.Worker {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]model.Worker, 0, len(r.workers))
	for _, worker := range r.workers {
		out = append(out, worker)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}
