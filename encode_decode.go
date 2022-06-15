package kafka_toolkit

import "context"

// DecodeConsumerMessageFunc extrae y convierte mensaje kafka a una estructura
type DecodeConsumerMessageFunc func(context.Context, *ConsumerMessage) (interface{}, error)

// DecodeConsumerMessageFuncMiddleware middleware para decode de consumer message
type DecodeConsumerMessageFuncMiddleware func(DecodeConsumerMessageFunc) DecodeConsumerMessageFunc

// EncodeProducerMessageFunc convierte estructura en mensaje kafka
type EncodeProducerMessageFunc func(context.Context, interface{}) (*ProducerMessage, error)

// EncodeProducerMessageFuncMiddleware middleware para encode producer message
type EncodeProducerMessageFuncMiddleware func(EncodeProducerMessageFunc) EncodeProducerMessageFunc
