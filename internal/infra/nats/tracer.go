package nats

import (
	"context"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"go.opentelemetry.io/otel"
)

type propagationNATSHeaderCarrier struct {
	headers nats.Header
}

func (c propagationNATSHeaderCarrier) Get(key string) string {
	return c.headers.Get(key)
}

func (c propagationNATSHeaderCarrier) Set(key, value string) {
	c.headers.Set(key, value)
}

func (c propagationNATSHeaderCarrier) Keys() []string {
	keys := make([]string, 0, len(c.headers))
	for k := range c.headers {
		keys = append(keys, k)
	}
	return keys
}

func buildTracerContext(ctx context.Context, msg jetstream.Msg) context.Context {
	return otel.GetTextMapPropagator().Extract(ctx, propagationNATSHeaderCarrier{headers: msg.Headers()})
}
