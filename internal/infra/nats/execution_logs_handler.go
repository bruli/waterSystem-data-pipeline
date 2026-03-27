package nats

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/bruli/waterSystem-data-pipeline/internal/domain/executed_logs"
	"github.com/nats-io/nats.go/jetstream"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type ExecutionLogsHandler struct {
	log    *slog.Logger
	svc    *executedlogs.Create
	tracer trace.Tracer
}

func (h ExecutionLogsHandler) Handle(msg jetstream.Msg) error {
	ctx := buildTracerContext(context.Background(), msg)
	ctx, span := h.tracer.Start(ctx, "ExecutionLogsHandler.Handle")
	defer span.End()
	var data ExecutionLog
	if err := json.Unmarshal(msg.Data(), &data); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return fmt.Errorf("failed to unmarshal execution log: %w", err)
	}
	el := executedlogs.New(data.Zone, data.Seconds, data.ExecutedAt)
	if err := h.svc.Execute(ctx, el); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		h.log.ErrorContext(ctx, "error saving execution log", slog.String("error", err.Error()))
		return err
	}
	span.SetStatus(codes.Ok, "execution log handled")
	return nil
}

func NewExecutionLogsHandler(log *slog.Logger, svc *executedlogs.Create, tracer trace.Tracer) *ExecutionLogsHandler {
	return &ExecutionLogsHandler{log: log, svc: svc, tracer: tracer}
}
