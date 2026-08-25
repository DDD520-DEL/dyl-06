package notify

import (
	"encoding/json"

	"github.com/dyl-06/telemetry/internal/model"
)

// Encode 把事件列表编码为 JSON 载荷。
func Encode(events []model.Event) []byte {
	payload, _ := json.Marshal(events)
	return payload
}

// DeliverTo 返回订阅者待推送事件并生成载荷。
func (r *Registry) DeliverTo(id string) ([]byte, bool) {
	events := r.Drain(id)
	if len(events) == 0 {
		return nil, false
	}
	return Encode(events), true
}
