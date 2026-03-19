package config

import "github.com/caarlos0/env/v11"

type Config struct {
	ServerHost     string `env:"SERVER_HOST,required"`
	NatsServerURL  string `env:"NATS_SERVER_URL,required"`
	InfluxDBURL    string `env:"INFLUXDB_URL,required"`
	InfluxDBToken  string `env:"INFLUXDB_TOKEN,required"`
	InfluxDBOrg    string `env:"INFLUXDB_ORG,required"`
	InfluxDBBucket string `env:"INFLUXDB_BUCKET,required"`
}

func New() (*Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
