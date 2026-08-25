package alert

import (
	"github.com/dyl-06/telemetry/internal/store"
)

// EvaluateAll 用所有已知序列的最新聚合值评估全部规则。
func (e *Engine) EvaluateAll(st *store.Store) {
	for _, series := range st.Series() {
		point, ok := st.LastPoint(series)
		if !ok {
			continue
		}
		bucket := st.GetBucket(series, point.Ts/60)
		if bucket == nil {
			continue
		}
		e.Evaluate(series, bucket.Read())
	}
}
