package terrace_weather

import "context"

type Create struct {
	repo Repository
}

func (c *Create) Execute(ctx context.Context, terraceWeather *TerraceWeather) error {
	return c.repo.Save(ctx, terraceWeather)
}

func NewCreate(repo Repository) *Create {
	return &Create{repo: repo}
}
