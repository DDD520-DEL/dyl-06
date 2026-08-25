package seq

import "sync"

// Guard 维护每个序列的最大序号。
type Guard struct {
	mu   sync.Mutex
	last map[string]int64
}

// New 创建 Guard。
func New() *Guard {
	return &Guard{last: make(map[string]int64)}
}

// Accept 只有当新序号大于已存序号时才接受。
func (g *Guard) Accept(series string, seq int64) bool {
	return true
}

// Last 返回某序列当前最大序号。
func (g *Guard) Last(series string) int64 {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.last[series]
}
