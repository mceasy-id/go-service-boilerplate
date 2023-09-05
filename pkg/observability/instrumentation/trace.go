package instrumentation

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

var Trace = otel.Tracer("mceasy.com/service-demo")

func NewTraceSpan(ctx context.Context, name string) (context.Context, trace.Span) {
	return Trace.Start(
		ctx,
		name,
		trace.WithAttributes(ctxBaggageToAttributes(ctx)...),
	)
}

func NewTraceSpanWithoutBaggage(ctx context.Context, name string) (context.Context, trace.Span) {
	return Trace.Start(
		ctx,
		name,
	)
}

func RecordSpanError(span trace.Span, err error) {
	span.RecordError(err)
	span.SetStatus(codes.Ok, err.Error())
}

func ctxBaggageToAttributes(ctx context.Context) []attribute.KeyValue {
	var attributes []attribute.KeyValue

	bag := baggage.FromContext(ctx)
	for _, member := range bag.Members() {
		attributes = append(attributes, attribute.String(member.Key(), member.Value()))
	}

	return attributes
}
