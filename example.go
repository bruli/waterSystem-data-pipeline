package main

import (
	"context"
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

func main() {
	ctx := context.Background()

	nc, err := nats.Connect("nats://localhost:4222")
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Drain()

	js, err := jetstream.New(nc)
	if err != nil {
		log.Fatal(err)
	}

	// Consumer només per execution.logs
	logsConsumer, err := js.CreateOrUpdateConsumer(ctx, "EVENTS", jetstream.ConsumerConfig{
		Durable:       "execution-logs-consumer",
		AckPolicy:     jetstream.AckExplicitPolicy,
		FilterSubject: "execution.logs",
	})
	if err != nil {
		log.Fatal(err)
	}

	// Consumer només per terrace.weather
	weatherConsumer, err := js.CreateOrUpdateConsumer(ctx, "EVENTS", jetstream.ConsumerConfig{
		Durable:       "terrace-weather-consumer",
		AckPolicy:     jetstream.AckExplicitPolicy,
		FilterSubject: "terrace.weather",
	})
	if err != nil {
		log.Fatal(err)
	}

	go consume(ctx, logsConsumer, handleExecutionLog)
	go consume(ctx, weatherConsumer, handleTerraceWeather)

	select {}
}

func consume(ctx context.Context, cons jetstream.Consumer, handler func(jetstream.Msg) error) {
	cc, err := cons.Consume(func(msg jetstream.Msg) {
		if err := handler(msg); err != nil {
			log.Printf("error processant missatge: %v", err)
			return // sense Ack -> es redelivera
		}

		if err := msg.Ack(); err != nil {
			log.Printf("error fent ack: %v", err)
		}
	})
	if err != nil {
		log.Printf("error iniciant consumer: %v", err)
		return
	}
	defer cc.Stop()

	<-ctx.Done()
}

func handleExecutionLog(msg jetstream.Msg) error {
	fmt.Printf("[execution.logs] %s\n", string(msg.Data()))
	return nil
}

func handleTerraceWeather(msg jetstream.Msg) error {
	fmt.Printf("[terrace.weather] %s\n", string(msg.Data()))
	return nil
}
