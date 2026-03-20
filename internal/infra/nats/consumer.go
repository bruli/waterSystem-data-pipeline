package nats

import (
	"context"
	"errors"
	"fmt"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

const StreamName = "WATER_SYSTEM"

type Consumer struct {
	js jetstream.JetStream
}

func (c Consumer) Create(ctx context.Context, subject string) (jetstream.Consumer, error) {
	co, err := c.js.CreateOrUpdateConsumer(ctx, StreamName, jetstream.ConsumerConfig{
		Durable:       fmt.Sprintf("%s-consumer", subject),
		AckPolicy:     jetstream.AckExplicitPolicy,
		FilterSubject: subject,
	})
	if err != nil {
		return nil, err
	}
	return co, nil
}

func NewConsumer(ctx context.Context, url string, subjects []string) (*Consumer, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, fmt.Errorf("nats.Connect: %w", err)
	}
	js, err := jetstream.New(nc)
	if err != nil {
		return nil, fmt.Errorf("jetstream.New: %w", err)
	}
	if err = ensureStream(ctx, js, subjects); err != nil {
		return nil, err
	}
	return &Consumer{js: js}, nil
}

func ensureStream(ctx context.Context, js jetstream.JetStream, subjects []string) error {
	stream, err := js.Stream(ctx, StreamName)
	if err == nil {
		info, err := stream.Info(ctx)
		if err != nil {
			return fmt.Errorf("stream info: %w", err)
		}

		info.Config.Subjects = subjects

		_, err = js.UpdateStream(ctx, info.Config)
		if err != nil {
			return fmt.Errorf("update stream: %w", err)
		}

		return nil
	}

	if !errors.Is(err, jetstream.ErrStreamNotFound) {
		return fmt.Errorf("get stream: %w", err)
	}

	_, err = js.CreateStream(ctx, jetstream.StreamConfig{
		Name:      StreamName,
		Subjects:  subjects,
		Storage:   jetstream.FileStorage,
		Retention: jetstream.LimitsPolicy,
	})
	if err != nil {
		return fmt.Errorf("create stream: %w", err)
	}

	return nil
}
