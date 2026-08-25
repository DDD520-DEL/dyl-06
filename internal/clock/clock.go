package clock

import "time"

// Clock 抽象时间来源，便于测试注入。
type Clock interface {
	Now() int64
}

// SystemClock 使用真实时间。
type SystemClock struct{}

// Now 返回当前 Unix 秒。
func (SystemClock) Now() int64 {
	return time.Now().Unix()
}
