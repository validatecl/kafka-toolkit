package kafka_toolkit

import (
	"context"
	"time"

	"github.com/go-kit/kit/metrics"
)

func StreamSuccesAndFailureMetricsMiddleware(requests metrics.Counter, operation string) StreamProcessorMiddleware {
	return func(process StreamProcessor) StreamProcessor {
		return func(ctx context.Context, inMsg *ConsumerMessage) (outMsg *ProducerMessage, err error) {
			defer func() {
				status := "SUCCESS"
				if err != nil {
					status = "ERROR"
				}
				requests.With("operation", operation, "status", status).Add(1)
			}()

			return process(ctx, inMsg)
		}
	}
}

func StreamTimeTakenMetricsMiddleware(duration metrics.Histogram, operation string) StreamProcessorMiddleware {
	return func(process StreamProcessor) StreamProcessor {
		return func(ctx context.Context, inMsg *ConsumerMessage) (outMsg *ProducerMessage, err error) {
			defer func(begin time.Time) {
				status := "SUCCESS"
				if err != nil {
					status = "ERROR"
				}

				duration.With("operation", operation, "status", status).Observe(time.Since(begin).Seconds())
			}(time.Now())
			return process(ctx, inMsg)
		}
	}
}
