package httpapi

import (
	"encoding/json"
	"net/http"

	"perx-deps-monitor/internal/monitor"
)

type API struct {
	Mon *monitor.Monitor
}

func (a *API) Mux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", a.health)
	mux.HandleFunc("/ready", a.ready)
	mux.HandleFunc("/status", a.status)
	return mux
}

func (a *API) health(w http.ResponseWriter, r *http.Request) {
	WriteJSON(w, http.StatusOK, map[string]any{"status": "ok"})
}

func (a *API) ready(w http.ResponseWriter, r *http.Request) {
	if a.Mon == nil || !a.Mon.HasCheckedAtLeastOnce() {
		WriteJSON(w, http.StatusServiceUnavailable, map[string]any{"status": "starting"})
		return
	}
	WriteJSON(w, http.StatusOK, map[string]any{"status": "ready"})
}

func (a *API) status(w http.ResponseWriter, r *http.Request) {
	if a.Mon == nil {
		WriteJSON(w, http.StatusOK, map[string]any{
			"mongodb":     "unknown",
			"nats":        "unknown",
			"checked_at":  nil,
			"started_at":  nil,
			"check_count": 0,
		})
		return
	}
	WriteJSON(w, http.StatusOK, a.Mon.Snapshot())
}

func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
