package main

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/bruli/waterSystem-data-pipeline/internal/config"
	"github.com/bruli/waterSystem-data-pipeline/internal/domain/terrace_weather"
	httpinfra "github.com/bruli/waterSystem-data-pipeline/internal/infra/http"
	"github.com/bruli/waterSystem-data-pipeline/internal/infra/influxdb2"
	"github.com/bruli/waterSystem-data-pipeline/internal/infra/nats"
)

func main() {
	ctx := context.Background()
	log := buildLog()
	conf, err := config.New()
	if err != nil {
		log.ErrorContext(ctx, "error loading config", "error", err)
		os.Exit(1)
	}
	serverListener, err := net.Listen("tcp", conf.ServerHost)
	log.InfoContext(ctx, "Starting server", "host", conf.ServerHost)
	if err != nil {
		log.ErrorContext(ctx, "Error starting server", "err", err)
		os.Exit(1)
	}
	defer func() {
		_ = serverListener.Close()
	}()

	consumer, err := nats.NewConsumer(ctx, conf.NatsServerURL, nats.Subjects)
	if err != nil {
		slog.ErrorContext(ctx, "Error starting consumer", "err", err.Error())
		os.Exit(1)
	}
	executionLogsConsumer, err := consumer.Create(ctx, nats.ExecutionLogsSubject)
	if err != nil {
		slog.ErrorContext(ctx, "Error starting execution logs consumer", "err", err.Error())
		os.Exit(1)
	}
	executionLogHandler := nats.NewExecutionLogsHandler(log)

	terraceWeatherConsumer, err := consumer.Create(ctx, nats.TerraceWeatherSubject)
	if err != nil {
		slog.ErrorContext(ctx, "Error starting terrace weather consumer", "err", err.Error())
		os.Exit(1)
	}
	twRepo := influxdb2.NewTerraceWeatherRepository(conf.InfluxDBURL, conf.InfluxDBToken, conf.InfluxDBOrg, conf.InfluxDBBucket)
	defer twRepo.Close()
	twCreate := terrace_weather.NewCreate(twRepo)
	terraceWeatherHandler := nats.NewTerraceWeatherHandler(log, twCreate)

	go nats.Consume(ctx, executionLogsConsumer, log, executionLogHandler.Handle)
	go nats.Consume(ctx, terraceWeatherConsumer, log, terraceWeatherHandler.Handle)

	srv := httpinfra.NewServer(conf.ServerHost)
	defer func() {
		log.InfoContext(ctx, "Closing server")
		_ = srv.Shutdown(ctx)
	}()

	runHTTPServer(ctx, srv, log, serverListener)
}

func buildLog() *slog.Logger {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	return slog.New(handler)
}

func runHTTPServer(ctx context.Context, srv *http.Server, log *slog.Logger, serverListener net.Listener) {
	go shutdown(ctx, srv, log)

	if err := srv.Serve(serverListener); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.ErrorContext(ctx, "Error starting server", "err", err)
		os.Exit(1)
	}
}

func shutdown(ctx context.Context, srv *http.Server, log *slog.Logger) {
	<-ctx.Done()
	log.InfoContext(ctx, "Ctrl+C received, shutting down server")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error("error shutting down server", "err", err)
	}
}
