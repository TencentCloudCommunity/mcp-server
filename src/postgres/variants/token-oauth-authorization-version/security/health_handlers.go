package security

import (
	"net/http"
	"time"
)

type HealthStatus struct {
	Service   string
	Version   string
	StartedAt time.Time
	Now       func() time.Time
}

func RegisterHealthRoutes(mux *http.ServeMux, status HealthStatus) {
	if mux == nil {
		return
	}
	if status.Service == "" {
		status.Service = "postgres-mcp-server"
	}
	if status.Version == "" {
		status.Version = "unknown"
	}
	if status.Now == nil {
		status.Now = func() time.Time { return time.Now().UTC() }
	}
	mux.Handle("/healthz", healthHandler(status, "live"))
	mux.Handle("/readyz", healthHandler(status, "ready"))
}

func healthHandler(status HealthStatus, probe string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet, http.MethodHead:
		default:
			writeJSON(w, http.StatusMethodNotAllowed, map[string]any{"error": "method not allowed"})
			return
		}
		now := status.Now().UTC()
		payload := map[string]any{
			"status":    "ok",
			"probe":     probe,
			"service":   status.Service,
			"version":   status.Version,
			"timestamp": now.Format(time.RFC3339),
		}
		if !status.StartedAt.IsZero() {
			startedAt := status.StartedAt.UTC()
			payload["started_at"] = startedAt.Format(time.RFC3339)
			payload["uptime_seconds"] = int64(now.Sub(startedAt).Seconds())
		}
		writeJSON(w, http.StatusOK, payload)
	})
}
