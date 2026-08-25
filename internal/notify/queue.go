package notify

import "github.com/dyl-06/telemetry/internal/model"

// enqueue 为订阅者追加一条待推送事件。
func (r *Registry) enqueue(id string, event model.Event) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.pending[id] = append(r.pending[id], event)
}

// Drain 取出订阅者的待推送事件并清空。
func (r *Registry) Drain(id string) []model.Event {
	r.mu.Lock()
	defer r.mu.Unlock()
	events := r.pending[id]
	delete(r.pending, id)
	return events
}

// LastStatus 返回订阅者某规则最近一次推送的状态。
func (r *Registry) LastStatus(id, rule string) string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.last[id + "/" + rule]
}

// recordStatus 记录订阅者某规则最近一次推送的状态。
func (r *Registry) recordStatus(id, rule, status string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.last[id+"/"+rule] = status
}
