package http

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/bruli/waterSystem-data-pipeline/internal/domain/forecast"
)

func NewServer(host string, predSvc *forecast.Prediction, log *slog.Logger) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	mux.HandleFunc("/prediction", func(w http.ResponseWriter, r *http.Request) {
		err := predSvc.Get(r.Context(), forecast.Tomorrow())
		if err != nil {
			log.ErrorContext(r.Context(), "Prediction HTTP Handler: error getting forecast", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	return &http.Server{
		Addr:         host,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}
