package ingest

import (
	"github.com/dyl-06/telemetry/internal/model"
	"github.com/dyl-06/telemetry/internal/store"
	"github.com/dyl-06/telemetry/internal/validation"
)

// Batch 提交一批点位：先整体校验，通过后整批写入。
func Batch(st *store.Store, points []model.Point) error {
	if err := validation.ValidateBatch(points); err != nil {
		return err
	}
	for _, p := range points {
		st.WritePoint(p)
	}
	return nil
}
