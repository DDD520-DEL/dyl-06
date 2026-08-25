package alert

import (
	"fmt"
	"sync"
)

// History 记录规则状态迁移历史。
type History struct {
	mu      sync.Mutex
	entries []string
	limit   int
}

// NewHistory 创建容量为 limit 的状态迁移历史。
func NewHistory(limit int) *History {
	return &History{limit: limit}
}

// Record 记录一次状态迁移。
func (h *History) Record(rule, status string, value float64) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.entries = append(h.entries, fmt.Sprintf("%s -> %s (%.2f)", rule, status, value))
	if len(h.entries) > h.limit {
		h.entries = h.entries[len(h.entries)-h.limit:]
	}
}

// Recent 返回最近的状态迁移。
func (h *History) Recent() []string {
	h.mu.Lock()
	defer h.mu.Unlock()
	out := make([]string, len(h.entries))
	copy(out, h.entries)
	return out
}
