package notify

import (
	"sync"

	"github.com/dyl-06/telemetry/internal/model"
)

// Registry 维护告警通知订阅者。
type Registry struct {
	mu      sync.RWMutex
	subs    map[string]*model.Subscription
	pending map[string][]model.Event
	last    map[string]string
}

// NewRegistry 创建订阅注册表。
func NewRegistry() *Registry {
	return &Registry{
		subs:    make(map[string]*model.Subscription),
		pending: make(map[string][]model.Event),
		last:    make(map[string]string),
	}
}

// Subscribe 注册订阅者。
func (r *Registry) Subscribe(sub *model.Subscription) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.subs[sub.ID] = sub
}

// Unsubscribe 移除订阅者并清空其待推送事件。
func (r *Registry) Unsubscribe(id string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.subs, id)
}

// Get 返回订阅者。
func (r *Registry) Get(id string) (*model.Subscription, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	sub, ok := r.subs[id]
	return sub, ok
}

// List 返回全部订阅者。
func (r *Registry) List() []*model.Subscription {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]*model.Subscription, 0, len(r.subs))
	for _, sub := range r.subs {
		out = append(out, sub)
	}
	return out
}

// Pending 返回订阅者待推送事件数。
func (r *Registry) Pending(id string) int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.pending[id])
}
