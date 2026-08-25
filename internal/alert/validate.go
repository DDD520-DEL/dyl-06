package alert

import (
	"errors"
	"math"

	"github.com/dyl-06/telemetry/internal/model"
)

var allowedOps = map[string]bool{"gt": true, "lt": true, "gte": true, "lte": true}

// ValidateRule 校验规则字段合法性。
func ValidateRule(rule model.Rule) error {
	if rule.Name == "" {
		return errors.New("rule name required")
	}
	if rule.Metric == "" {
		return errors.New("rule metric required")
	}
	if !allowedOps[rule.Op] {
		return errors.New("unsupported operator")
	}
	if math.IsNaN(rule.Threshold) || math.IsInf(rule.Threshold, 0) {
		return errors.New("threshold must be finite")
	}
	if rule.Window <= 0 {
		return errors.New("window must be positive")
	}
	return nil
}
