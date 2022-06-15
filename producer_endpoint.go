package kafka_toolkit

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

// MakeMessageProducerEndpoint crea un endpoint para producir mensajes
func MakeMessageProducerEndpoint(producer MessageProducer, encode EncodeProducerMessageFunc) endpoint.Endpoint {
	return func(ctx context.Context, in interface{}) (interface{}, error) {
		msg, err := encode(ctx, in)
		if err != nil {
			return nil, err
		}

		err = producer.SendMessage(ctx, msg)

		if err != nil {
			return nil, err
		}

		return nil, nil
	}
}
