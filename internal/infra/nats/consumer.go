package nats

import (
	"context"
	"fmt"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type Consumer struct {
	js jetstream.JetStream
}

func (c Consumer) Create(ctx context.Context, subject string) (jetstream.Consumer, error) {
	co, err := c.js.CreateOrUpdateConsumer(ctx, "EVENTS", jetstream.ConsumerConfig{
		Durable:       fmt.Sprintf("%s-consumer", subject),
		AckPolicy:     jetstream.AckExplicitPolicy,
		FilterSubject: subject,
	})
	if err != nil {
		return nil, err
	}
	return co, nil
}

func NewConsumer(url string) (*Consumer, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, fmt.Errorf("nats.Connect: %w", err)
	}
	js, err := jetstream.New(nc)
	if err != nil {
		return nil, fmt.Errorf("jetstream.New: %w", err)
	}
	return &Consumer{js: js}, nil
}
