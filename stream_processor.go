package kafka_toolkit

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/endpoint"
)

// StreamProcessor funcion streamer que procesa mensaje de entrada, llama a endpoint y genera un producer message
type StreamProcessor func(ctx context.Context, inMsg *ConsumerMessage) (*ProducerMessage, error)

// MakeStreamProcessor crea un nuevo streamer
func MakeStreamProcessor(process endpoint.Endpoint, decode DecodeConsumerMessageFunc, encode EncodeProducerMessageFunc) StreamProcessor {
	return func(ctx context.Context, inMsg *ConsumerMessage) (*ProducerMessage, error) {
		in, err := decode(ctx, inMsg)
		if err != nil {
			Log.Error(fmt.Sprintf("Error encodeando mensaje %v", err))
			return nil, err
		}

		out, err := process(ctx, in)
		if err != nil {
			Log.Error(fmt.Sprintf("Error procesando mensaje %v", err))
			return nil, err
		}

		outMsg, err := encode(ctx, out)
		if err != nil {
			Log.Error(fmt.Sprintf("Error decodeando mensaje %v", err))
			return nil, err
		}

		return outMsg, err
	}
}
