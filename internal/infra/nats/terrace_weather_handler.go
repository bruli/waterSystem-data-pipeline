package nats

import (
	"encoding/json"
	"log/slog"

	"github.com/nats-io/nats.go/jetstream"
)

type TerraceWeatherHandler struct {
	log *slog.Logger
}

func (h TerraceWeatherHandler) Handle(msg jetstream.Msg) error {
	slog.Info("[terrace_weather]", slog.String("data", string(msg.Data())))
	var data TerraceWeather
	if err := json.Unmarshal(msg.Data(), &data); err != nil {
		return err
	}
	slog.Info(
		"received data",
		slog.Float64("temperature", data.Temperature),
		slog.Bool("is_raining", data.IsRaining),
		slog.Time("executed_at", data.ExecutedAt),
	)
	return nil
}

func NewTerraceWeatherHandler(log *slog.Logger) *TerraceWeatherHandler {
	return &TerraceWeatherHandler{log: log}
}
