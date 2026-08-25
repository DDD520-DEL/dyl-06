package downsample

import (
	"github.com/dyl-06/telemetry/internal/store"
)

// Manager 负责关闭到期时间桶。
type Manager struct {
	store      *store.Store
	bucketSize int64
}

// New 创建降采样管理器。
func New(st *store.Store, bucketSize int64) *Manager {
	return &Manager{store: st, bucketSize: bucketSize}
}

// CloseBucket 关闭指定序列的指定时间桶。
func (m *Manager) CloseBucket(series string, slot int64) {
	m.store.CloseBucket(series, slot)
}

// CloseExpired 关闭所有早于 cutoff 的桶。
func (m *Manager) CloseExpired(cutoff int64) {
	for _, series := range m.store.Series() {
		m.store.CloseBucket(series, cutoff/m.bucketSize)
	}
}
