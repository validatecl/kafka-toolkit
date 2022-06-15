package kafka_toolkit

import (
	"context"
)

type loggingMiddlewareMessageHandler struct {
	next MessageHandler
}

// MakeLoggingMessageHandlerMiddleware constructor para logging message handler
func MakeLoggingMessageHandlerMiddleware() MessageHandlerMiddleware {
	return func(next MessageHandler) MessageHandler {
		return &loggingMiddlewareMessageHandler{
			next: next,
		}
	}
}

func (h *loggingMiddlewareMessageHandler) HandleMessage(ctx context.Context, msg *ConsumerMessage) (err error) {
	defer func() {

		Log.Info("Mensaje recibido", string(msg.Msg), "key", string(msg.Key), "timestamp", msg.Timestamp, "topic", msg.Topic, "partition", msg.Partition, "offset", msg.Offset)

		if err != nil {
			Log.Error("Error", err.Error(), "topic", msg.Topic, "partition", msg.Partition, "offset", msg.Offset)
		}

	}()

	return h.next.HandleMessage(ctx, msg)
}
