package influxdb2

import (
	"context"
	"fmt"

	"github.com/bruli/waterSystem-data-pipeline/internal/domain/forecast"
	influxdb "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
)

type ForecastWeatherRepository struct {
	client influxdb.Client
	org    string
	bucket string
}

func (f ForecastWeatherRepository) Save(ctx context.Context, weather *forecast.Weather) error {
	writeAPI := f.client.WriteAPIBlocking(f.org, f.bucket)

	point := write.NewPoint(
		"forecast_v2",
		map[string]string{
			"location": "terrace",
		},
		map[string]interface{}{
			"temperature":               weather.Temperature(),
			"relative_humidity":         float64(weather.RelativeHumidity()),
			"precipitation_probability": weather.PrecipitationProbability(),
			"cloud_cover":               float64(weather.CloudCover()),
			"shortwave_radiation":       weather.ShortwaveRadiation(),
			"forecast_generated_at":     weather.GeneratedAt().Unix(),
		},
		weather.Hour(),
	)
	if err := writeAPI.WritePoint(ctx, point); err != nil {
		return fmt.Errorf("write forecast point to influxdb: %w", err)
	}

	return nil
}

func NewForecastWeatherRepository(url, token, org, bucket string) *ForecastWeatherRepository {
	client := influxdb.NewClient(url, token)
	return &ForecastWeatherRepository{client: client, org: org, bucket: bucket}
}
