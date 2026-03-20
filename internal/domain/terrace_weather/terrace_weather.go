package terrace_weather

import "time"

type TerraceWeather struct {
	temperature float64
	isRaining   bool
	measuredAt  time.Time
}

func (t TerraceWeather) Temperature() float64 {
	return t.temperature
}

func (t TerraceWeather) IsRaining() bool {
	return t.isRaining
}

func (t TerraceWeather) MeasuredAt() time.Time {
	return t.measuredAt
}

func New(temperature float64, isRaining bool, measuredAt time.Time) *TerraceWeather {
	return &TerraceWeather{temperature: temperature, isRaining: isRaining, measuredAt: measuredAt}
}
