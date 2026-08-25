package ingest

import (
	"github.com/dyl-06/telemetry/internal/clock"
	"github.com/dyl-06/telemetry/internal/model"
	"github.com/dyl-06/telemetry/internal/registry"
)

// RegisterWorker 注册一个摄取 worker。
func RegisterWorker(reg *registry.Registry, id string, clk clock.Clock) {
	reg.Register(model.Worker{ID: id, LastSeen: clk.Now()})
}

// Heartbeat 刷新 worker 心跳。
func Heartbeat(reg *registry.Registry, id string, clk clock.Clock) {
	reg.Heartbeat(id, clk.Now())
}
