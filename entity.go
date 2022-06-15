package kafka_toolkit

import "time"

// ConsumerMessage representa un mensaje consumido desde un topico
type ConsumerMessage struct {
	Headers   map[string]string
	Timestamp time.Time
	Key       []byte
	Msg       []byte
	Topic     string
	Partition int32
	Offset    int64
}

// ProducerMessage representa un mensaje a enviar desde un topico
type ProducerMessage struct {
	Headers map[string]string
	Key     []byte
	Msg     []byte
}
