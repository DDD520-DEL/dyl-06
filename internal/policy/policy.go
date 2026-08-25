package policy

import (
	"errors"
	"regexp"
)

var seriesNameRE = regexp.MustCompile(`^[a-zA-Z0-9_.\-]{1,128}$`)

// MaxBatchSize 单批最大点位数量。
const MaxBatchSize = 5000

// ValidateSeriesName 校验序列名。
func ValidateSeriesName(name string) error {
	if !seriesNameRE.MatchString(name) {
		return errors.New("invalid series name")
	}
	return nil
}

// ValidateShards 校验分片列表合法。
func ValidateShards(shards []string) error {
	if len(shards) == 0 {
		return errors.New("at least one shard required")
	}
	seen := make(map[string]bool, len(shards))
	for _, shard := range shards {
		if shard == "" {
			return errors.New("shard name must not be empty")
		}
		if seen[shard] {
			return errors.New("duplicate shard name")
		}
		seen[shard] = true
	}
	return nil
}

// RetentionPolicy 描述保留窗口。
type RetentionPolicy struct {
	Slots int64
}

// Validate 校验保留策略。
func (p RetentionPolicy) Validate() error {
	if p.Slots <= 0 {
		return errors.New("retention slots must be positive")
	}
	return nil
}
