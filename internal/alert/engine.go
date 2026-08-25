package alert

import (
	"sync"

	"github.com/dyl-06/telemetry/internal/clock"
	"github.com/dyl-06/telemetry/internal/metrics"
	"github.com/dyl-06/telemetry/internal/model"
	"github.com/dyl-06/telemetry/internal/notify"
)

// Engine 按规则评估指标并推进状态机。
type Engine struct {
	rules      *RuleSet
	notifier   *notify.Dispatcher
	clock      clock.Clock
	metrics    *metrics.Metrics
	mu         sync.Mutex
	states     map[string]model.RuleState
	lastValue  map[string]float64
	history    *History
}

// SetHistory 挂接状态迁移历史。
func (e *Engine) SetHistory(history *History) {
	e.history = history
}

// NewEngine 创建告警引擎。
func NewEngine(rules *RuleSet, notifier *notify.Dispatcher, clk clock.Clock, m *metrics.Metrics) *Engine {
	return &Engine{
		rules:     rules,
		notifier:  notifier,
		clock:     clk,
		metrics:   m,
		states:    make(map[string]model.RuleState),
		lastValue: make(map[string]float64),
	}
}

// Evaluate 用最新聚合值评估所有规则并推送状态变化事件。
func (e *Engine) Evaluate(series string, result model.AggResult) {
	now := e.clock.Now()
	for _, rule := range e.rules.List() {
		if rule.Metric != series {
			continue
		}
		status := "ok"
		if triggered(rule.Op, result.Sum, rule.Threshold) {
			status = "alert"
		}
		e.mu.Lock()
		prev := e.states[rule.Name].Status
		e.states[rule.Name] = model.RuleState{Status: status, Since: now}
		e.lastValue[rule.Name] = result.Sum
		e.mu.Unlock()
		if status != prev {
			if e.history != nil {
				e.history.Record(rule.Name, status, result.Sum)
			}
			e.notifier.Dispatch(model.Event{
				Rule:   rule.Name,
				Series: series,
				Status: status,
				Value:  result.Sum,
			})
			e.metrics.Alerts.Add(1)
		}
	}
}

// State 返回规则当前状态。
func (e *Engine) State(name string) model.RuleState {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.states[name]
}

func triggered(op string, value, threshold float64) bool {
	switch op {
	case "gt":
		return value > threshold
	case "lt":
		return value < threshold
	default:
		return false
	}
}
