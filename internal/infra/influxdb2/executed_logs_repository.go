package influxdb2

import (
	"context"
	"fmt"
	"strings"

	"github.com/bruli/waterSystem-data-pipeline/internal/domain/executed_logs"
	influxdb "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type ExecutedLogsRepository struct {
	client influxdb.Client
	org    string
	bucket string
	tracer trace.Tracer
}

func (e *ExecutedLogsRepository) Save(ctx context.Context, executedLog *executedlogs.ExecutedLog) error {
	ctx, span := e.tracer.Start(ctx, "ExecutedLogsRepository.Save")
	defer span.End()
	writeAPI := e.client.WriteAPIBlocking(e.org, e.bucket)

	point := write.NewPoint(
		"logs",
		map[string]string{
			"zone": formatZoneName(executedLog.Zone()),
		},
		map[string]interface{}{
			"seconds": executedLog.Seconds(),
		},
		executedLog.ExecutedAt(),
	)

	if err := writeAPI.WritePoint(ctx, point); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return fmt.Errorf("write point to influxdb: %w", err)
	}
	span.SetStatus(codes.Ok, "executed log saved")
	return nil
}

func formatZoneName(zone string) string {
	return strings.ReplaceAll(zone, " with fertilizer", "")
}

func (e *ExecutedLogsRepository) Close() {
	e.client.Close()
}

func NewExecutedLogsRepository(url, token, org, bucket string, tracer trace.Tracer) *ExecutedLogsRepository {
	client := influxdb.NewClient(url, token)
	return &ExecutedLogsRepository{client: client, org: org, bucket: bucket, tracer: tracer}
}
