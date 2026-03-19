package nats

import (
	"context"
	"log/slog"

	"github.com/nats-io/nats.go/jetstream"
)

func Consume(ctx context.Context, cons jetstream.Consumer, log *slog.Logger, handler func(msg jetstream.Msg) error) {
	cc, err := cons.Consume(func(msg jetstream.Msg) {
		if err := handler(msg); err != nil {
			slog.ErrorContext(ctx, "error processing message", slog.String("error", err.Error()))
			return
		}
		if err := msg.Ack(); err != nil {
			log.ErrorContext(ctx, "error acknowledging message", slog.String("error", err.Error()))
		}
	})
	if err != nil {
		log.ErrorContext(ctx, "error consuming message", slog.String("error", err.Error()))
		return
	}
	defer cc.Stop()

	<-ctx.Done()
}
