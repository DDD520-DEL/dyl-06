package verifycase

import (
	"testing"

	"github.com/dyl-06/telemetry/internal/metrics"
	"github.com/dyl-06/telemetry/internal/model"
	"github.com/dyl-06/telemetry/internal/notify"
)

func TestUnsubscribeDrainsPendingEvents(t *testing.T) {
	registry := notify.NewRegistry()
	registry.Subscribe(&model.Subscription{ID: "s1", Webhook: "http://ops", Rules: []string{"high-cpu"}})
	m := metrics.New()
	dispatcher := notify.NewDispatcher(registry, m)
	dispatcher.Dispatch(model.Event{Rule: "high-cpu", Series: "cpu.usage", Status: "alert", Value: 90})
	if got := registry.Pending("s1"); got != 1 {
		t.Fatalf("pending before unsubscribe = %d, want 1", got)
	}
	registry.Unsubscribe("s1")
	if got := registry.Pending("s1"); got != 0 {
		t.Fatalf("pending after unsubscribe = %d, want 0 (queue leak)", got)
	}
}
