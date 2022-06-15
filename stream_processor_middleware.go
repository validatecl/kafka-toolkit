package kafka_toolkit

import (
	"context"
	"fmt"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

const headersKey = "kafka-headers"

// StreamProcessorMiddleware middleware de streamer
type StreamProcessorMiddleware func(StreamProcessor) StreamProcessor

// MakeHeaderPassThruStreamProcessorMiddleware middleware de headers pass thru
func MakeHeaderPassThruStreamProcessorMiddleware() StreamProcessorMiddleware {
	return func(process StreamProcessor) StreamProcessor {
		return func(ctx context.Context, inMsg *ConsumerMessage) (*ProducerMessage, error) {
			ctx = headersToContext(ctx, inMsg)

			outMsg, err := process(ctx, inMsg)
			if err != nil {
				return nil, err
			}

			contextToHeaders(ctx, outMsg)

			return outMsg, nil
		}
	}
}

func headersToContext(ctx context.Context, inMsg *ConsumerMessage) context.Context {
	headers := inMsg.Headers

	if headers == nil {
		headers = make(map[string]string)
	}

	return context.WithValue(ctx, headersKey, headers)
}

func contextToHeaders(ctx context.Context, outMsg *ProducerMessage) {
	if outMsg.Headers == nil {
		outMsg.Headers = make(map[string]string)
	}

	headers := ctx.Value(headersKey).(map[string]string)

	for k, v := range headers {
		outMsg.Headers[k] = v
	}
}

// Open Tracing

// MakeOpenTracingStreamProcessorMiddleware crea un nuevo middleware de open tracing
func MakeOpenTracingStreamProcessorMiddleware(operationName string, tracer opentracing.Tracer) StreamProcessorMiddleware {
	return func(process StreamProcessor) StreamProcessor {
		return func(ctx context.Context, inMsg *ConsumerMessage) (*ProducerMessage, error) {

			ctxWithSpan := kafkaTraceToContext(ctx, inMsg, tracer, operationName)
			// Finish span when the stream processesor reach the end of the scope :D
			if span := opentracing.SpanFromContext(ctxWithSpan); span != nil {
				defer span.Finish()
			}

			outMsg, err := process(ctxWithSpan, inMsg)
			if err != nil {
				return nil, err
			}

			// Propagate the span before defer span.Finish()
			messageWithTraceHeader := contextToKafkaTrace(ctxWithSpan, outMsg, tracer)

			return messageWithTraceHeader, nil
		}
	}
}

func kafkaTraceToContext(ctx context.Context, inMsg *ConsumerMessage, tracer opentracing.Tracer, operationName string) context.Context {
	if inMsg.Headers == nil {
		inMsg.Headers = make(map[string]string)
	}

	var span opentracing.Span
	wireContext, _ := tracer.Extract(
		opentracing.TextMap,
		opentracing.TextMapCarrier(inMsg.Headers),
	)

	span = tracer.StartSpan(operationName, ext.RPCServerOption(wireContext))
	ext.SpanKind.Set(span, ext.SpanKindConsumerEnum)
	span.SetTag("consumer-kind", "kafka")

	return opentracing.ContextWithSpan(ctx, span)
}

func contextToKafkaTrace(ctx context.Context, outMsg *ProducerMessage, tracer opentracing.Tracer) *ProducerMessage {
	if outMsg.Headers == nil {
		outMsg.Headers = make(map[string]string)
	}

	if span := opentracing.SpanFromContext(ctx); span != nil {
		tracer.Inject(
			span.Context(),
			opentracing.TextMap,
			opentracing.TextMapCarrier(outMsg.Headers),
		)
	}

	return outMsg
}

// MakeLoggingStreamProcessorMiddleware logging stream processor middleware
func MakeLoggingStreamProcessorMiddleware() StreamProcessorMiddleware {
	return func(process StreamProcessor) StreamProcessor {
		return func(ctx context.Context, inMsg *ConsumerMessage) (out *ProducerMessage, err error) {

			defer func() {
				Log.Debug(fmt.Sprintf("Mensaje stream de entrada: llave = %q, msg = %q, timestamp = %v, in_topic = %s, partition=%d, offset=%d",
					string(inMsg.Key),
					string(inMsg.Msg),
					inMsg.Timestamp,
					inMsg.Topic,
					inMsg.Partition,
					inMsg.Offset))

				if err != nil {
					Log.Error(fmt.Sprintf("Error de streaming: %v", err))
				}

				if out != nil {
					Log.Debug(fmt.Sprintf("Mensaje stream de salida: llave = %q, msg = %q, headers = %v",
						string(out.Key),
						string(out.Msg),
						out.Headers))
				}

			}()

			return process(ctx, inMsg)
		}
	}
}
