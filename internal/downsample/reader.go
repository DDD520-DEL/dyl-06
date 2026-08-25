package downsample

import (
	"github.com/dyl-06/telemetry/internal/model"
	"github.com/dyl-06/telemetry/internal/store"
)

// Reader 读取降采样桶结果。
type Reader struct {
	store *store.Store
}

// NewReader 创建读取器。
func NewReader(st *store.Store) *Reader {
	return &Reader{store: st}
}

// ReadBucket 返回序列时间桶聚合结果；桶不存在时返回空结果。
func (r *Reader) ReadBucket(series string, slot int64) (model.AggResult, error) {
	bucket := r.store.GetBucket(series, slot)
	return bucket.Read(), nil
}
