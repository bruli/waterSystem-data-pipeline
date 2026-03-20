package nats

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/nats-io/nats.go/jetstream"
)

type ExecutionLogsHandler struct {
	log *slog.Logger
}

func (h ExecutionLogsHandler) Handle(msg jetstream.Msg) error {
	slog.Info("[execution_logs]", slog.String("data", string(msg.Data())))
	var data ExecutionLog
	if err := json.Unmarshal(msg.Data(), &data); err != nil {
		return fmt.Errorf("failed to unmarshal execution log: %w", err)
	}
	slog.Info(
		"execution data:",
		slog.String("zone", data.Zone),
		slog.Int("seconds", data.Seconds),
		slog.Time("executed_at", data.ExecutedAt),
	)
	return nil
}

func NewExecutionLogsHandler(log *slog.Logger) *ExecutionLogsHandler {
	return &ExecutionLogsHandler{log: log}
}
