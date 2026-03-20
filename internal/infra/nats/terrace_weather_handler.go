package nats

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/bruli/waterSystem-data-pipeline/internal/domain/terrace_weather"
	"github.com/nats-io/nats.go/jetstream"
)

type TerraceWeatherHandler struct {
	log *slog.Logger
	svc *terrace_weather.Create
}

func (h TerraceWeatherHandler) Handle(msg jetstream.Msg) error {
	ctx := context.Background()
	var data TerraceWeather
	if err := json.Unmarshal(msg.Data(), &data); err != nil {
		return err
	}
	tw := terrace_weather.New(data.Temperature, data.IsRaining, data.ExecutedAt)
	if err := h.svc.Execute(ctx, tw); err != nil {
		h.log.ErrorContext(ctx, "error saving terrace weather", slog.String("error", err.Error()))
		return err
	}
	return nil
}

func NewTerraceWeatherHandler(log *slog.Logger, svc *terrace_weather.Create) *TerraceWeatherHandler {
	return &TerraceWeatherHandler{log: log, svc: svc}
}
