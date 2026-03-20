package executedlogs

import "context"

type Create struct {
	repo Repository
}

func (c *Create) Execute(ctx context.Context, executedLog *ExecutedLog) error {
	return c.repo.Save(ctx, executedLog)
}

func NewCreate(repo Repository) *Create {
	return &Create{repo: repo}
}
