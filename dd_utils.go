package kafka_toolkit

import (
	"context"
	"strconv"

	opentracing "github.com/opentracing/opentracing-go"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
)

func GetDatadogTraceAndSpanFromContext(ctx context.Context) (ddTraceId string, ddSpanId string) {

	span := opentracing.SpanFromContext(ctx)
	{
		spanCtx := span.Context()
		ddSpanCtx, isDatadogContext := spanCtx.(ddtrace.SpanContext)
		if isDatadogContext {
			ddSpanId = strconv.FormatUint(ddSpanCtx.SpanID(), 10)
			ddTraceId = strconv.FormatUint(ddSpanCtx.TraceID(), 10)
		}
	}

	return ddTraceId, ddSpanId
}
