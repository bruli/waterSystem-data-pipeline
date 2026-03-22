package api

import (
	"context"
	"encoding/json"
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
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("http.NewRequestWithContext: %w", err)
	}
	resp, err := o.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http.DefaultClient.Do: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	var results OpenMeteoResponse
	if err := json.Unmarshal(body, &results); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}
	weather := make([]*forecast.Weather, len(results.Hourly.Time))
	for i := range results.Hourly.Time {
		time, err := time.Parse("2006-01-02T15:04", results.Hourly.Time[i])
		if err != nil {
			return nil, fmt.Errorf("failed to parse time: %w", err)
		}
		weather[i] = forecast.NewWeather(
			time,
			results.Hourly.Temperature2m[i],
			results.Hourly.RelativeHumidity2m[i],
			results.Hourly.Precipitation[i],
			results.Hourly.CloudCover[i],
			results.Hourly.ShortwaveRadiation[i],
		)
	}
	return weather, nil
}

func NewOpenMeteoReader() *OpenMeteoReader {
	return &OpenMeteoReader{client: &http.Client{Timeout: 10 * time.Second}}
}
