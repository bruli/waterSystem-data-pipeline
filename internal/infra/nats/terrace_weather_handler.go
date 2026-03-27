package nats

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/bruli/waterSystem-data-pipeline/internal/domain/terrace_weather"
	"github.com/nats-io/nats.go/jetstream"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type TerraceWeatherHandler struct {
	log    *slog.Logger
	svc    *terrace_weather.Create
	tracer trace.Tracer
}

func (h TerraceWeatherHandler) Handle(msg jetstream.Msg) error {
	ctx := buildTracerContext(context.Background(), msg)
	ctx, span := h.tracer.Start(ctx, "TerraceWeatherHandler.Handle")
	defer span.End()
	var data TerraceWeather
	if err := json.Unmarshal(msg.Data(), &data); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	tw := terrace_weather.New(data.Temperature, data.IsRaining, data.ExecutedAt)
	if err := h.svc.Execute(ctx, tw); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		h.log.ErrorContext(ctx, "error saving terrace weather", slog.String("error", err.Error()))
		return err
	}
	span.SetStatus(codes.Ok, "terrace weather handled")
	return nil
}

func NewTerraceWeatherHandler(log *slog.Logger, svc *terrace_weather.Create, tracer trace.Tracer) *TerraceWeatherHandler {
	return &TerraceWeatherHandler{log: log, svc: svc, tracer: tracer}
}
