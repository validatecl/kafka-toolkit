package kafka_toolkit

import (
	"context"
	"time"

	opentracing "github.com/opentracing/opentracing-go"
	commons "github.com/validatecl/go-microservices-commons"
)

type (
	// ProducerMessageMiddleware producer message middleware
	ProducerMessageMiddleware func(MessageProducer) MessageProducer

	opentracingProducerMiddleware struct {
		next   MessageProducer
		tracer opentracing.Tracer
	}

	metricsProducerMiddleware struct {
		next   MessageProducer
		config *commons.MetricsConfig
	}
)

// MakeOpentracingProducerMiddleware producer with context propagation
func MakeOpentracingProducerMiddleware(operationName string, tracer opentracing.Tracer) ProducerMessageMiddleware {

	return func(next MessageProducer) MessageProducer {
		return &opentracingProducerMiddleware{next, tracer}
	}
}

func (mw *opentracingProducerMiddleware) SendMessage(ctx context.Context, msg *ProducerMessage) error {

	messageWithContext := contextToKafkaTrace(ctx, msg, mw.tracer)

	return mw.next.SendMessage(ctx, messageWithContext)
}

// MakeMetricsProducerMiddleware producer metrics middleware
func MakeMetricsProducerMiddleware(config *commons.MetricsConfig) ProducerMessageMiddleware {
	return func(next MessageProducer) MessageProducer {
		return &metricsProducerMiddleware{next, config}
	}
}

func (mw *metricsProducerMiddleware) SendMessage(ctx context.Context, msg *ProducerMessage) error {

	var err error
	defer func(begin time.Time) {

		if err != nil {
			mw.config.RequestDuration.With("operation", "produce", "status", "ERROR").Observe(time.Since(begin).Seconds())
			mw.config.RequestCount.With("operation", "send message", "status", "ERROR").Add(1)
		} else {
			mw.config.RequestDuration.With("operation", "produce", "status", "STATUS OK").Observe(time.Since(begin).Seconds())
			mw.config.RequestCount.With("operation", "send message", "status", "STATUS OK").Add(1)
		}

	}(time.Now())

	err = mw.next.SendMessage(ctx, msg)

	return err
}
