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
	executedlogs "github.com/bruli/waterSystem-data-pipeline/internal/domain/executed_logs"
	"github.com/bruli/waterSystem-data-pipeline/internal/domain/forecast"
	"github.com/bruli/waterSystem-data-pipeline/internal/domain/terrace_weather"
	apiinfra "github.com/bruli/waterSystem-data-pipeline/internal/infra/api"
	httpinfra "github.com/bruli/waterSystem-data-pipeline/internal/infra/http"
	"github.com/bruli/waterSystem-data-pipeline/internal/infra/influxdb2"
	"github.com/bruli/waterSystem-data-pipeline/internal/infra/nats"
	"github.com/bruli/waterSystem-data-pipeline/internal/infra/tracing"
	"github.com/robfig/cron/v3"
	"go.opentelemetry.io/otel"
)

const serviceName = "waterSystem-data-pipeline"

func main() {
	ctx := context.Background()
	log := buildLog()
	conf, err := config.New()
	if err != nil {
		log.ErrorContext(ctx, "error loading config", "error", err)
		os.Exit(1)
	}

	tracingProv, err := tracing.InitTracing(ctx, serviceName)
	if err != nil {
		log.ErrorContext(ctx, "Error initializing tracing", "err", err)
		os.Exit(1)
	}
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err = tracingProv.Shutdown(shutdownCtx); err != nil {
			log.ErrorContext(ctx, "Error shutting down tracing", "err", err)
		}
	}()

	tracer := otel.Tracer(serviceName)

	serverListener, err := net.Listen("tcp", conf.ServerHost)
	log.InfoContext(ctx, "Starting server", "host", conf.ServerHost)
	if err != nil {
		log.ErrorContext(ctx, "Error starting server", "err", err)
		os.Exit(1)
	}
	defer func() {
		_ = serverListener.Close()
	}()

	forecastPred := buildForecastPrediction(conf, log)

	cron, err := buildCron()
	if err != nil {
		log.ErrorContext(ctx, "Error creating cron", "err", err)
		os.Exit(1)
	}
	forecastCh := make(chan struct{})
	defer close(forecastCh)

	go recurrentForecast(ctx, forecastCh, forecastPred, log)
	go forecastCron(ctx, cron, forecastPred, log, forecastCh)

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
	execLogRepo := influxdb2.NewExecutedLogsRepository(conf.InfluxDBURL, conf.InfluxDBToken, conf.InfluxDBOrg, conf.InfluxDBBucket, tracer)
	defer execLogRepo.Close()
	elCreate := executedlogs.NewCreate(execLogRepo)
	executionLogHandler := nats.NewExecutionLogsHandler(log, elCreate, tracer)

	terraceWeatherConsumer, err := consumer.Create(ctx, nats.TerraceWeatherSubject)
	if err != nil {
		slog.ErrorContext(ctx, "Error starting terrace weather consumer", "err", err.Error())
		os.Exit(1)
	}
	twRepo := influxdb2.NewTerraceWeatherRepository(conf.InfluxDBURL, conf.InfluxDBToken, conf.InfluxDBOrg, conf.InfluxDBBucket, tracer)
	defer twRepo.Close()
	twCreate := terrace_weather.NewCreate(twRepo)
	terraceWeatherHandler := nats.NewTerraceWeatherHandler(log, twCreate, tracer)

	go nats.Consume(ctx, executionLogsConsumer, log, executionLogHandler.Handle)
	go nats.Consume(ctx, terraceWeatherConsumer, log, terraceWeatherHandler.Handle)

	srv := httpinfra.NewServer(conf.ServerHost, forecastPred, log)
	defer func() {
		log.InfoContext(ctx, "Closing server")
		_ = srv.Shutdown(ctx)
	}()

	runHTTPServer(ctx, srv, log, serverListener)
}

func recurrentForecast(ctx context.Context, ch <-chan struct{}, pred *forecast.Prediction, log *slog.Logger) {
	for {
		select {
		case <-ctx.Done():
			log.InfoContext(ctx, "Recurrent forecast stopped")
			return

		case <-ch:
			log.InfoContext(ctx, "Recurrent forecast started...")
			ticker := time.NewTicker(30 * time.Minute)
			defer ticker.Stop()

			for {
				select {
				case <-ctx.Done():
					log.InfoContext(ctx, "Recurrent forecast context done")
					return

				case <-ticker.C:
					err := pred.Get(ctx, forecast.Tomorrow())
					if err != nil {
						log.ErrorContext(ctx, "Error getting forecast on recurrent process", "err", err)
						continue
					}

					log.InfoContext(ctx, "Recurrent forecast prediction finished")
				}
			}
		}
	}
}

func forecastCron(ctx context.Context, c *cron.Cron, pred *forecast.Prediction, log *slog.Logger, ch chan<- struct{}) {
	defer c.Stop()
	_, err := c.AddFunc("* 7 * * *", func() {
		if err := pred.Get(ctx, forecast.Tomorrow()); err != nil {
			log.ErrorContext(ctx, "Error getting forecast", "err", err)
			ch <- struct{}{}
		}
		log.InfoContext(ctx, "Forecast prediction finished")
	})
	if err != nil {
		log.ErrorContext(ctx, "Error adding cron job", "err", err)
		os.Exit(1)
	}
	log.InfoContext(ctx, "Forecast cron prediction started")
	c.Start()
	<-ctx.Done()
	log.InfoContext(ctx, "Forecast cron prediction stopped")
}

func buildForecastPrediction(conf *config.Config, log *slog.Logger) *forecast.Prediction {
	foreCastWeathRepo := influxdb2.NewForecastWeatherRepository(conf.InfluxDBURL, conf.InfluxDBToken, conf.InfluxDBOrg, conf.InfluxDBBucket)
	opMetReader := apiinfra.NewOpenMeteoReader()
	return forecast.NewPrediction(opMetReader, foreCastWeathRepo, log)
}

func buildLog() *slog.Logger {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	log := slog.New(handler)
	log.With("service", serviceName)
	return log
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

func buildCron() (*cron.Cron, error) {
	loc, err := time.LoadLocation("Europe/Madrid")
	if err != nil {
		return nil, err
	}
	c := cron.New(cron.WithLocation(loc))
	return c, nil
}
