package nats

import "time"

type TerraceWeather struct {
	Temperature float64   `json:"temperature"`
	IsRaining   bool      `json:"is_raining"`
	ExecutedAt  time.Time `json:"executed_at"`
}

type ExecutionLog struct {
	Zone       string    `json:"zone_name"`
	Seconds    int       `json:"seconds"`
	ExecutedAt time.Time `json:"executed_at"`
}
