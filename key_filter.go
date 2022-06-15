package kafka_toolkit

import (
	"context"
	"fmt"
)

type keyFilterMessageMiddleware struct {
	allowedKeys []string
	next        MessageHandler
}

// MakeKeyFilterMessageHandlerMiddleware key filter middleware
func MakeKeyFilterMessageHandlerMiddleware(allowedKeys []string) MessageHandlerMiddleware {
	return func(next MessageHandler) MessageHandler {
		return &keyFilterMessageMiddleware{
			allowedKeys: allowedKeys,
			next:        next,
		}
	}
}

func (k *keyFilterMessageMiddleware) HandleMessage(ctx context.Context, inMsg *ConsumerMessage) error {
	msgKey := string(inMsg.Key)

	for _, key := range k.allowedKeys {
		if key == msgKey {
			return k.next.HandleMessage(ctx, inMsg)
		}
	}

	return fmt.Errorf("Key no soportada %s - debe ser una de: %v", msgKey, k.allowedKeys)
}
