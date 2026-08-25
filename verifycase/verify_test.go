package verifycase

import (
	"testing"
	"time"

	"github.com/dyl-06/telemetry/internal/alert"
	"github.com/dyl-06/telemetry/internal/harness"
	"github.com/dyl-06/telemetry/internal/metrics"
	"github.com/dyl-06/telemetry/internal/model"
	"github.com/dyl-06/telemetry/internal/notify"
)

func TestAlertRecoveryNotifiesSubscribers(t *testing.T) {
	rules := alert.NewRuleSet()
	rules.Upsert(model.Rule{Name: "high-cpu", Metric: "cpu.usage", Op: "gt", Threshold: 80, Window: 60})
	registry := notify.NewRegistry()
	registry.Subscribe(&model.Subscription{ID: "s1", Webhook: "http://ops", Rules: []string{"high-cpu"}})
	m := metrics.New()
	dispatcher := notify.NewDispatcher(registry, m)
	engine := alert.NewEngine(rules, dispatcher, harness.NewFakeClock(time.Unix(1000, 0)), m)

	engine.Evaluate("cpu.usage", model.AggResult{Sum: 90})
	if got := registry.Pending("s1"); got != 1 {
		t.Fatalf("pending after alert = %d, want 1", got)
	}
	engine.Evaluate("cpu.usage", model.AggResult{Sum: 30})
	if got := registry.Pending("s1"); got != 2 {
		t.Fatalf("pending after recovery = %d, want 2 (recovery event missing)", got)
	}
}
