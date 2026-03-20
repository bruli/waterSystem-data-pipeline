package influxdb2

import (
	"context"
	"fmt"

	"github.com/bruli/waterSystem-data-pipeline/internal/domain/executed_logs"
	influxdb "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
)

type ExecutedLogsRepository struct {
	client influxdb.Client
	org    string
	bucket string
}

func (e *ExecutedLogsRepository) Save(ctx context.Context, executedLog *executedlogs.ExecutedLog) error {
	writeAPI := e.client.WriteAPIBlocking(e.org, e.bucket)

	point := write.NewPoint(
		"logs",
		map[string]string{
			"location": "terrace",
		},
		map[string]interface{}{
			"zone":    executedLog.Zone(),
			"seconds": executedLog.Seconds(),
		},
		executedLog.ExecutedAt(),
	)

	if err := writeAPI.WritePoint(ctx, point); err != nil {
		return fmt.Errorf("write point to influxdb: %w", err)
	}

	return nil
}

func (e *ExecutedLogsRepository) Close() {
	e.client.Close()
}

func NewExecutedLogsRepository(url, token, org, bucket string) *ExecutedLogsRepository {
	client := influxdb.NewClient(url, token)
	return &ExecutedLogsRepository{client: client, org: org, bucket: bucket}
}
