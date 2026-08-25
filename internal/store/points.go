package store

import (
	"github.com/dyl-06/telemetry/internal/model"
	"github.com/dyl-06/telemetry/internal/seq"
)

// WritePointWithGuard 写入点位前先做序列号单调守卫。
func (s *Store) WritePointWithGuard(guard *seq.Guard, p model.Point) bool {
	if !guard.Accept(p.Series, p.Seq) {
		return false
	}
	s.WritePoint(p)
	return true
}
