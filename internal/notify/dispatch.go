package notify

import (
	"github.com/dyl-06/telemetry/internal/metrics"
	"github.com/dyl-06/telemetry/internal/model"
)

// Dispatcher 负责向订阅者分发告警事件。
type Dispatcher struct {
	registry *Registry
	metrics  *metrics.Metrics
}

// NewDispatcher 创建分发器。
func NewDispatcher(registry *Registry, m *metrics.Metrics) *Dispatcher {
	return &Dispatcher{registry: registry, metrics: m}
}

// Dispatch 向订阅了该规则的客户端入队事件。
func (d *Dispatcher) Dispatch(event model.Event) {
	for _, sub := range d.registry.List() {
		if !contains(sub.Rules, event.Rule) {
			continue
		}
		if d.registry.LastStatus(sub.ID, event.Rule) == event.Status {
			continue
		}
		d.registry.enqueue(sub.ID, event)
		d.registry.recordStatus(sub.ID, event.Rule, event.Status)
		d.metrics.Events.Add(1)
	}
}

func contains(items []string, target string) bool {
	for _, item := range items {
		if item == target {
			return true
		}
	}
	return false
}
