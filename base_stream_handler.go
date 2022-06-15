package kafka_toolkit

import (
	"context"
)

type streamHandler struct {
	producer MessageProducer
	process  StreamProcessor
}

// NewStreamHandler constructor para  message modifier handler
func NewStreamHandler(producer MessageProducer, process StreamProcessor) MessageHandler {
	return &streamHandler{
		producer: producer,
		process:  process,
	}
}

func (h *streamHandler) HandleMessage(ctx context.Context, inMsg *ConsumerMessage) error {

	outMsg, err := h.process(ctx, inMsg)
	if err != nil {
		return err
	}

	if err := h.producer.SendMessage(ctx, outMsg); err != nil {
		Log.Error(
			"errorMessage", "Error enviando mensaje",
			"message", err.Error(),
			"error", err)
		return err
	}

	return nil
}
