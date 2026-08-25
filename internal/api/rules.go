package api

import (
	"encoding/json"
	"net/http"

	"github.com/dyl-06/telemetry/internal/alert"
	"github.com/dyl-06/telemetry/internal/model"
)

// HandleRuleUpdate 热更新告警规则。
func HandleRuleUpdate(rules *alert.RuleSet) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var rule model.Rule
		if err := json.NewDecoder(r.Body).Decode(&rule); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := alert.ValidateRule(rule); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		rules.Upsert(rule)
		w.WriteHeader(http.StatusOK)
	}
}
