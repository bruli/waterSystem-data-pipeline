package influxdb2

import (
	"context"
	"fmt"

	"github.com/bruli/waterSystem-data-pipeline/internal/domain/terrace_weather"
	influxdb "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
)

type TerraceWeatherRepository struct {
	client influxdb.Client
	org    string
	bucket string
}

func (t *TerraceWeatherRepository) Save(ctx context.Context, terraceWeather *terrace_weather.TerraceWeather) error {
	writeAPI := t.client.WriteAPIBlocking(t.org, t.bucket)

	point := write.NewPoint(
		"weather",
		map[string]string{
			"location": "terrace",
		},
		map[string]interface{}{
			"temperature": terraceWeather.Temperature(),
			"is_raining":  terraceWeather.IsRaining(),
		},
		terraceWeather.MeasuredAt(),
	)

	if err := writeAPI.WritePoint(ctx, point); err != nil {
		return fmt.Errorf("write point to influxdb: %w", err)
	}

	return nil
}

func (t *TerraceWeatherRepository) Close() {
	t.client.Close()
}

func NewTerraceWeatherRepository(url, token, org, bucket string) *TerraceWeatherRepository {
	client := influxdb.NewClient(url, token)
	return &TerraceWeatherRepository{client: client, org: org, bucket: bucket}
}
