package api

import (
	"encoding/json"
	"net/http"

	"github.com/dyl-06/telemetry/internal/ingest"
	"github.com/dyl-06/telemetry/internal/model"
	"github.com/dyl-06/telemetry/internal/store"
)

// HandleBatchIngest 批量上报点位，整批原子写入。
func HandleBatchIngest(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Points []model.Point `json:"points"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := ingest.Batch(st, req.Points); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusAccepted)
	}
}
