package api

import (
	"log"
	"net/http"
	"time"

	"github.com/dyl-06/telemetry/internal/audit"
)

// LogRequests 包装 handler：记录审计与耗时。
func LogRequests(a *audit.Audit, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		elapsed := time.Since(start)
		a.Record(r.Method + " " + r.URL.Path + " " + elapsed.Round(time.Millisecond).String())
		log.Printf("%s %s %s", r.Method, r.URL.Path, elapsed)
	})
}
