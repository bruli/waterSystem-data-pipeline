package nats

import (
	"log/slog"

	"github.com/nats-io/nats.go/jetstream"
)

type ExecutionLogsHandler struct {
	log *slog.Logger
}

func (h ExecutionLogsHandler) Handle(msg jetstream.Msg) error {
	slog.Info("[execution.logs]", slog.String("data", string(msg.Data())))
	return nil
}

func NewExecutionLogsHandler(log *slog.Logger) *ExecutionLogsHandler {
	return &ExecutionLogsHandler{log: log}
}
