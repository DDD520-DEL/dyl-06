package audit

import "sync"

// Audit 保存最近操作记录。
type Audit struct {
	mu      sync.Mutex
	entries []string
	limit   int
}

// New 创建容量为 limit 的审计记录。
func New(limit int) *Audit {
	return &Audit{limit: limit}
}

// Record 追加一条记录。
func (a *Audit) Record(entry string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.entries = append(a.entries, entry)
	if len(a.entries) > a.limit {
		a.entries = a.entries[len(a.entries)-a.limit:]
	}
}

// Recent 返回最近的记录。
func (a *Audit) Recent() []string {
	a.mu.Lock()
	defer a.mu.Unlock()
	out := make([]string, len(a.entries))
	copy(out, a.entries)
	return out
}
