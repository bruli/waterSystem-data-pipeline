package http

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/bruli/waterSystem-data-pipeline/internal/domain/forecast"
)

const (
	readTimeout  = 10 * time.Second
	writeTimeout = 10 * time.Second
)

func NewServer(host string, predSvc *forecast.Prediction, log *slog.Logger) *http.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", healthzHandler)
	mux.HandleFunc("/prediction", predictionHandler(predSvc, log))

	return &http.Server{
		Addr:         host,
		Handler:      mux,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}
}

func healthzHandler(w http.ResponseWriter, _ *http.Request) {
	writeText(w, http.StatusOK, "ok")
}

func predictionHandler(predSvc *forecast.Prediction, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeText(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		if err := predSvc.Get(r.Context(), forecast.Tomorrow()); err != nil {
			log.ErrorContext(
				r.Context(),
				"Prediction HTTP Handler: error getting forecast",
				slog.String("error", err.Error()),
			)
			writeText(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func writeText(w http.ResponseWriter, statusCode int, body string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(statusCode)
	_, _ = w.Write([]byte(body))
}
