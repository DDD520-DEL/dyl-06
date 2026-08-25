package notify

import (
	"github.com/dyl-06/telemetry/internal/metrics"
)

// DeliveryWorker 周期清空订阅者待推送队列。
type DeliveryWorker struct {
	registry *Registry
	metrics  *metrics.Metrics
}

// NewDeliveryWorker 创建投递 worker。
func NewDeliveryWorker(registry *Registry, m *metrics.Metrics) *DeliveryWorker {
	return &DeliveryWorker{registry: registry, metrics: m}
}

// Deliver 清空所有在线订阅者的待推送队列。
func (d *DeliveryWorker) Deliver() int {
	delivered := 0
	for _, sub := range d.registry.List() {
		payload, ok := d.registry.DeliverTo(sub.ID)
		if !ok {
			continue
		}
		delivered += len(payload)
	}
	return delivered
}
