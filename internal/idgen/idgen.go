package idgen

import (
	"fmt"
	"sync/atomic"
)

// IDGen 生成单调递增 ID。
type IDGen struct {
	counter atomic.Int64
}

// New 创建 IDGen。
func New() *IDGen {
	return &IDGen{}
}

// Next 返回下一个 ID。
func (g *IDGen) Next() string {
	return fmt.Sprintf("id-%010d", g.counter.Add(1))
}
