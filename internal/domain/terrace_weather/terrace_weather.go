package terrace_weather

import "time"

type TerraceWeather struct {
	temperature float64
	isRaining   bool
	measuredAt  time.Time
	humidity    float64
}

func (t TerraceWeather) Humidity() float64 {
	return t.humidity
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

func New(temperature, humidity float64, isRaining bool, measuredAt time.Time) *TerraceWeather {
	return &TerraceWeather{
		temperature: temperature,
		isRaining:   isRaining,
		measuredAt:  measuredAt,
		humidity:    humidity,
	}
}
