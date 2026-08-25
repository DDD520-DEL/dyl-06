package api

import (
	"encoding/json"
	"net/http"

	"github.com/dyl-06/telemetry/internal/idgen"
	"github.com/dyl-06/telemetry/internal/model"
	"github.com/dyl-06/telemetry/internal/notify"
)

// HandleSubscribe 注册告警通知订阅者。
func HandleSubscribe(registry *notify.Registry, gen *idgen.IDGen) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var sub model.Subscription
		if err := json.NewDecoder(r.Body).Decode(&sub); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if sub.ID == "" {
			sub.ID = gen.Next()
		}
		registry.Subscribe(&sub)
		w.WriteHeader(http.StatusOK)
	}
}

// HandleUnsubscribe 移除订阅者。
func HandleUnsubscribe(registry *notify.Registry) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			ID string `json:"id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		registry.Unsubscribe(req.ID)
		w.WriteHeader(http.StatusOK)
	}
}
