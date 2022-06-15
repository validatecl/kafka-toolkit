package kafka_toolkit

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type baseHandler struct {
	exec   endpoint.Endpoint
	decode DecodeConsumerMessageFunc
}

// NewBaseHandler constructor para base message handler
func NewBaseHandler(ep endpoint.Endpoint, decode DecodeConsumerMessageFunc) MessageHandler {
	return &baseHandler{
		exec:   ep,
		decode: decode,
	}
}

func (h *baseHandler) HandleMessage(ctx context.Context, inMsg *ConsumerMessage) error {

	traceId, spanId := GetDatadogTraceAndSpanFromContext(ctx)

	in, err := h.decode(ctx, inMsg)
	if err != nil {
		Log.Error(
			"errorMessage", "Error encodeando mensaje",
			"error", err,
			"message", err.Error(),
			"dd.trace_id", traceId,
			"dd.span_id", spanId,
		)

		return err
	}

	_, err = h.exec(ctx, in)
	if err != nil {
		Log.Error(
			"errorMessage", "Error procesando mensaje",
			"error", err,
			"message", err.Error(),
			"dd.trace_id", traceId,
			"dd.span_id", spanId,
		)
		return err
	}

	return nil
}
