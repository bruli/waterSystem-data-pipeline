package executedlogs

import "context"

type Repository interface {
	Save(ctx context.Context, executedLog *ExecutedLog) error
}
