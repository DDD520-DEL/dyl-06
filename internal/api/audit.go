package api

import (
	"encoding/json"
	"net/http"

	"github.com/dyl-06/telemetry/internal/audit"
)

// HandleAudit 返回审计记录。
func HandleAudit(a *audit.Audit) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(a.Recent())
	}
}
