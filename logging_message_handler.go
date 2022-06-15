package kafka_toolkit

import (
	"context"
	"fmt"
)

type loggingMessageHandler struct{}

// NewLoggingMessageHandler constructor para logging message handler
func NewLoggingMessageHandler() MessageHandler {
	return &loggingMessageHandler{}
}

func (h *loggingMessageHandler) HandleMessage(ctx context.Context, msg *ConsumerMessage) error {
	Log.Info(fmt.Sprintf("Mensaje recibido: llave = %q, valor = %q, timestamp = %v, topic = %s, partition=%d, offset=%d",
		string(msg.Key),
		(msg.Msg),
		msg.Timestamp,
		msg.Topic,
		msg.Partition,
		msg.Offset))
	return nil
}
