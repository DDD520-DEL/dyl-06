package alert

import (
	"sort"
	"sync"

	"github.com/dyl-06/telemetry/internal/model"
)

// RuleSet 维护告警规则集。
type RuleSet struct {
	mu    sync.RWMutex
	rules map[string]model.Rule
}

// NewRuleSet 创建规则集。
func NewRuleSet() *RuleSet {
	return &RuleSet{rules: make(map[string]model.Rule)}
}

// Upsert 新增或替换规则。
func (rs *RuleSet) Upsert(rule model.Rule) {
	rs.rules[rule.Name] = rule
}

// Get 返回规则。
func (rs *RuleSet) Get(name string) (model.Rule, bool) {
	rs.mu.RLock()
	defer rs.mu.RUnlock()
	rule, ok := rs.rules[name]
	return rule, ok
}

// List 返回全部规则。
func (rs *RuleSet) List() []model.Rule {
	rs.mu.RLock()
	defer rs.mu.RUnlock()
	out := make([]model.Rule, 0, len(rs.rules))
	for _, rule := range rs.rules {
		out = append(out, rule)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}
