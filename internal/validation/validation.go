package validation

import (
	"errors"
	"math"

	"github.com/dyl-06/telemetry/internal/model"
	"github.com/dyl-06/telemetry/internal/policy"
)

// ValidatePoint 校验单个点位。
func ValidatePoint(p model.Point) error {
	if err := policy.ValidateSeriesName(p.Series); err != nil {
		return err
	}
	if p.Seq < 0 {
		return errors.New("sequence must be non-negative")
	}
	if p.Ts <= 0 {
		return errors.New("timestamp must be positive")
	}
	if math.IsNaN(p.Value) || math.IsInf(p.Value, 0) {
		return errors.New("value must be finite")
	}
	return nil
}

// ValidateBatch 校验整批点位。
func ValidateBatch(points []model.Point) error {
	if len(points) == 0 {
		return errors.New("empty batch")
	}
	if len(points) > policy.MaxBatchSize {
		return errors.New("batch too large")
	}
	for _, p := range points {
		if err := ValidatePoint(p); err != nil {
			return err
		}
	}
	return nil
}
