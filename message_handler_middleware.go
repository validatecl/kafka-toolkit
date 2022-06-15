package kafka_toolkit

import (
	"context"
	"time"

	commons "github.com/validatecl/go-microservices-commons"
)

type (
	messageOffsetPartitionHandlerMdw struct {
		next MessageHandler
	}

	messageMetricsMiddleware struct {
		next   MessageHandler
		config *commons.MetricsConfig
	}
)

// MakeOffsetPartitionCtxMiddleware ...
func MakeOffsetPartitionCtxMiddleware() MessageHandlerMiddleware {
	return func(next MessageHandler) MessageHandler {
		return newOffsetPartitionWithContextMiddleware(next)
	}
}

// MakeMessageHandlerMetricsMiddleware message handler metrics middleware
func MakeMessageHandlerMetricsMiddleware(config *commons.MetricsConfig) MessageHandlerMiddleware {
	return func(next MessageHandler) MessageHandler {
		return newMessageHandlerMetricsMiddleware(next, config)
	}
}

func newOffsetPartitionWithContextMiddleware(next MessageHandler) MessageHandler {
	return &messageOffsetPartitionHandlerMdw{next}
}

func (m *messageOffsetPartitionHandlerMdw) HandleMessage(ctx context.Context, inMsg *ConsumerMessage) error {

	ctx = ContextWithOffset(ctx, inMsg.Offset)
	ctx = ContextWithPartition(ctx, inMsg.Partition)

	return m.next.HandleMessage(ctx, inMsg)
}

func newMessageHandlerMetricsMiddleware(next MessageHandler, config *commons.MetricsConfig) MessageHandler {

	return &messageMetricsMiddleware{next, config}
}

func (mw *messageMetricsMiddleware) HandleMessage(ctx context.Context, inMsg *ConsumerMessage) error {

	var err error
	defer func(begin time.Time) {
		if err != nil {
			mw.config.RequestDuration.With("operation", "consume", "status", "ERROR").Observe(time.Since(begin).Seconds())
			mw.config.RequestCount.With("operation", "handle message", "status", "ERROR").Add(1)
		} else {
			mw.config.RequestDuration.With("operation", "consume", "status", "STATUS OK").Observe(time.Since(begin).Seconds())
			mw.config.RequestCount.With("operation", "handle message", "status", "STATUS OK").Add(1)
		}

	}(time.Now())

	err = mw.next.HandleMessage(ctx, inMsg)

	return err
}
