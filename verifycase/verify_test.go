package verifycase

import (
	"fmt"
	"sync"
	"testing"

	"github.com/dyl-06/telemetry/internal/alert"
	"github.com/dyl-06/telemetry/internal/model"
)

func TestRuleSetConcurrentSafe(t *testing.T) {
	rules := alert.NewRuleSet()
	const writers = 8
	const perWriter = 50
	var wg sync.WaitGroup
	for i := 0; i < writers; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			for j := 0; j < perWriter; j++ {
				rules.Upsert(model.Rule{
					Name:      fmt.Sprintf("rule-%d", j%5),
					Metric:    "cpu.usage",
					Op:        "gt",
					Threshold: float64(n),
					Window:    60,
				})
			}
		}(i)
	}
	wg.Wait()
	if got := len(rules.List()); got != 5 {
		t.Fatalf("rules = %d, want 5 (lost updates or race)", got)
	}
}
