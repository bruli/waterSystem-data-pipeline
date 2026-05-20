package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/bruli/waterSystem-data-pipeline/internal/domain/forecast"
)

const (
	OpenMeteoAPI = "https://api.open-meteo.com/v1/forecast"
	Latitude     = 41.5451707
	Longitude    = 2.1032168
)

var items = []string{"temperature_2m", "relative_humidity_2m", "precipitation", "cloudcover", "shortwave_radiation"}

type OpenMeteoResponse struct {
	Hourly HourlyData `json:"hourly"`
}
type HourlyData struct {
	Time               []string  `json:"time"`
	Temperature2m      []float64 `json:"temperature_2m"`
	RelativeHumidity2m []int     `json:"relative_humidity_2m"`
	Precipitation      []float64 `json:"precipitation"`
	CloudCover         []int     `json:"cloudcover"`
	ShortwaveRadiation []float64 `json:"shortwave_radiation"`
}

type OpenMeteoReader struct {
	client *http.Client
}

func (o OpenMeteoReader) Read(ctx context.Context, slot *forecast.Slot) ([]*forecast.Weather, error) {
	it := strings.Join(items, ",")
	url := fmt.Sprintf(
		"%s?latitude=%v&longitude=%v&hourly=%s&start_date=%s&end_date=%s&timezone=auto",
		OpenMeteoAPI,
		Latitude,
		Longitude,
		it,
		slot.From().Format("2006-01-02"),
		slot.To().Format("2006-01-02"),
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("http.NewRequestWithContext: %w", err)
	}

	resp, err := o.client.Do(req)
	if err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			return nil, fmt.Errorf("open-meteo request timeout: %w", err)
		case errors.Is(err, context.Canceled):
			return nil, fmt.Errorf("open-meteo request canceled: %w", err)
		default:
			return nil, fmt.Errorf("open-meteo request failed: %w", err)
		}
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
		return nil, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var results OpenMeteoResponse
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	n := len(results.Hourly.Time)
	if n == 0 {
		return []*forecast.Weather{}, nil
	}

	if len(results.Hourly.Temperature2m) != n ||
		len(results.Hourly.RelativeHumidity2m) != n ||
		len(results.Hourly.Precipitation) != n ||
		len(results.Hourly.CloudCover) != n ||
		len(results.Hourly.ShortwaveRadiation) != n {
		return nil, fmt.Errorf(
			"inconsistent hourly arrays length: time=%d temp=%d humidity=%d precipitation=%d cloudcover=%d radiation=%d",
			n,
			len(results.Hourly.Temperature2m),
			len(results.Hourly.RelativeHumidity2m),
			len(results.Hourly.Precipitation),
			len(results.Hourly.CloudCover),
			len(results.Hourly.ShortwaveRadiation),
		)
	}

	generatedAt := time.Now()
	weather := make([]*forecast.Weather, n)

	for i := 0; i < n; i++ {
		parsedTime, err := time.Parse("2006-01-02T15:04", results.Hourly.Time[i])
		if err != nil {
			return nil, fmt.Errorf("failed to parse time at index %d (%q): %w", i, results.Hourly.Time[i], err)
		}

		weather[i] = forecast.NewWeather(
			parsedTime,
			results.Hourly.Temperature2m[i],
			results.Hourly.RelativeHumidity2m[i],
			results.Hourly.Precipitation[i],
			results.Hourly.CloudCover[i],
			results.Hourly.ShortwaveRadiation[i],
			generatedAt,
		)
	}

	return weather, nil
}

func NewOpenMeteoReader() *OpenMeteoReader {
	return &OpenMeteoReader{client: &http.Client{Timeout: 10 * time.Second}}
}
