package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/dyl-06/telemetry/internal/downsample"
)

// HandleQuery 查询序列时间桶聚合结果；带 from/to 时返回降采样范围。
func (s *Server) HandleQuery(reader *downsample.Reader) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		series := r.URL.Query().Get("series")
		if fromRaw := r.URL.Query().Get("from"); fromRaw != "" {
			from, errFrom := strconv.ParseInt(fromRaw, 10, 64)
			to, errTo := strconv.ParseInt(r.URL.Query().Get("to"), 10, 64)
			if errFrom != nil || errTo != nil {
				http.Error(w, "from/to required", http.StatusBadRequest)
				return
			}
			_ = json.NewEncoder(w).Encode(reader.ReadRange(series, from, to))
			return
		}
		slot, err := strconv.ParseInt(r.URL.Query().Get("slot"), 10, 64)
		if err != nil {
			http.Error(w, "slot required", http.StatusBadRequest)
			return
		}
		result, err := reader.ReadBucket(series, slot)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		_ = json.NewEncoder(w).Encode(result)
	}
}
