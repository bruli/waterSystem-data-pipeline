package terrace_weather

import "context"

type Repository interface {
	Save(ctx context.Context, terraceWeather *TerraceWeather) error
}
