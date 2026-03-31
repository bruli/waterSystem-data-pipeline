package influxdb2

import (
	"context"
	"fmt"

	"github.com/bruli/waterSystem-data-pipeline/internal/domain/terrace_weather"
	influxdb "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type TerraceWeatherRepository struct {
	client influxdb.Client
	org    string
	bucket string
	tracer trace.Tracer
}

func (t *TerraceWeatherRepository) Save(ctx context.Context, terraceWeather *terrace_weather.TerraceWeather) error {
	ctx, span := t.tracer.Start(ctx, "TerraceWeatherRepository.Save")
	defer span.End()
	writeAPI := t.client.WriteAPIBlocking(t.org, t.bucket)

	point := write.NewPoint(
		"weather",
		map[string]string{
			"location": "terrace",
		},
		map[string]interface{}{
			"temperature": terraceWeather.Temperature(),
			"is_raining":  terraceWeather.IsRaining(),
			"humidity":    terraceWeather.Humidity(),
		},
		terraceWeather.MeasuredAt(),
	)

	if err := writeAPI.WritePoint(ctx, point); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return fmt.Errorf("write point to influxdb: %w", err)
	}
	span.SetStatus(codes.Ok, "terrace weather saved")
	return nil
}

func (t *TerraceWeatherRepository) Close() {
	t.client.Close()
}

func NewTerraceWeatherRepository(url, token, org, bucket string, tracer trace.Tracer) *TerraceWeatherRepository {
	client := influxdb.NewClient(url, token)
	return &TerraceWeatherRepository{client: client, org: org, bucket: bucket, tracer: tracer}
}
