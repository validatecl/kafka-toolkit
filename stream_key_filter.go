package kafka_toolkit

import (
	"context"
	"fmt"
)

// MakeKeyFilterStreamMiddleware middleware de filtrado de keys para streams
func MakeKeyFilterStreamMiddleware(allowedKeys []string) StreamProcessorMiddleware {
	return func(next StreamProcessor) StreamProcessor {
		return func(ctx context.Context, inMsg *ConsumerMessage) (*ProducerMessage, error) {
			msgKey := string(inMsg.Key)

			for _, key := range allowedKeys {
				if key == msgKey {
					return next(ctx, inMsg)
				}
			}

			return nil, fmt.Errorf("Key no soportada %s - debe ser una de: %v", msgKey, allowedKeys)
		}
	}
}
