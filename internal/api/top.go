package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/dyl-06/telemetry/internal/store"
)

// HandleTop 返回聚合值最高的序列。
func HandleTop(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limit := 10
		if raw := r.URL.Query().Get("limit"); raw != "" {
			if parsed, err := strconv.Atoi(raw); err == nil && parsed > 0 {
				limit = parsed
			}
		}
		_ = json.NewEncoder(w).Encode(st.QueryTopN(limit))
	}
}
