package tracing

import (
	"github.com/bruli/raspberryWaterSystem/pkg/cqs"
	"go.opentelemetry.io/otel/trace"
)

type Event struct {
	SpanContext trace.SpanContext
	Event       cqs.Event
}
