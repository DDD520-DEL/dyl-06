package harness

import (
	"time"
)

// FakeClock 是测试用固定时钟。
type FakeClock struct {
	now time.Time
}

// NewFakeClock 创建固定时钟。
func NewFakeClock(now time.Time) *FakeClock {
	return &FakeClock{now: now}
}

// Advance 推进时间。
func (f *FakeClock) Advance(d time.Duration) {
	f.now = f.now.Add(d)
}

// Now 返回当前时间。
func (f *FakeClock) Now() int64 {
	return f.now.Unix()
}
