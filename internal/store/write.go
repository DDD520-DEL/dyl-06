package store

import (
	"errors"

	"github.com/dyl-06/telemetry/internal/model"
	"github.com/dyl-06/telemetry/internal/validation"
)

// WritePoints 整批写入；任一点位非法时整批不落库。
func (s *Store) WritePoints(points []model.Point) error {
	if err := validation.ValidateBatch(points); err != nil {
		return err
	}
	for _, p := range points {
		s.WritePoint(p)
	}
	return nil
}

// ErrBatchRejected 表示整批被拒绝。
var ErrBatchRejected = errors.New("batch rejected")
