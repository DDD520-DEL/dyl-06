package api

import (
	"encoding/json"
	"net/http"

	"github.com/dyl-06/telemetry/internal/store"
)

// HandleShards 返回分片承载统计。
func HandleShards(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeRoutes := make(map[string]int)
		for _, series := range st.Series() {
			writeRoutes[st.ShardFor(series)]++
		}
		_ = json.NewEncoder(w).Encode(map[string]int(writeRoutes))
	}
}
