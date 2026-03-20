package nats

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/bruli/waterSystem-data-pipeline/internal/domain/executed_logs"
	"github.com/nats-io/nats.go/jetstream"
)

type ExecutionLogsHandler struct {
	log *slog.Logger
	svc *executedlogs.Create
}

func (h ExecutionLogsHandler) Handle(msg jetstream.Msg) error {
	ctx := context.Background()
	var data ExecutionLog
	if err := json.Unmarshal(msg.Data(), &data); err != nil {
		return fmt.Errorf("failed to unmarshal execution log: %w", err)
	}
	el := executedlogs.New(data.Zone, data.Seconds, data.ExecutedAt)
	if err := h.svc.Execute(ctx, el); err != nil {
		h.log.ErrorContext(ctx, "error saving execution log", slog.String("error", err.Error()))
		return err
	}
	return nil
}

func NewExecutionLogsHandler(log *slog.Logger, svc *executedlogs.Create) *ExecutionLogsHandler {
	return &ExecutionLogsHandler{log: log, svc: svc}
}
