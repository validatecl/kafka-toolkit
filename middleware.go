package kafka_toolkit

// MessageHandlerMiddleware Middleware function para message handler
type MessageHandlerMiddleware func(MessageHandler) MessageHandler

// MessageProducerMiddleware Middleware function para message producer
type MessageProducerMiddleware func(MessageProducer) MessageProducer
