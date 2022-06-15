package kafka_toolkit

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Shopify/sarama"
)

const (
	// SendMessageOperation nombre de operacion SendMessage
	SendMessageOperation = "SendMessage"
)

//MessageProducer interface para envio mensaje
type MessageProducer interface {
	SendMessage(ctx context.Context, msg *ProducerMessage) error
}

type baseProducer struct {
	topic    string
	producer sarama.SyncProducer
}

// NewBaseSaramaProducer constructor para enviar mensaje
func NewBaseSaramaProducer(topic string, producer sarama.SyncProducer) MessageProducer {
	return &baseProducer{topic: topic, producer: producer}

}

func (b *baseProducer) SendMessage(ctx context.Context, msg *ProducerMessage) error {
	cant := len(msg.Msg)
	traceId, spanId := GetDatadogTraceAndSpanFromContext(ctx)

	if cant < 1 {
		return errors.New(InvalidInputProducerErrorKind)
	}

	producerMsg := &sarama.ProducerMessage{
		Topic:     b.topic,
		Value:     sarama.StringEncoder(msg.Msg),
		Key:       sarama.StringEncoder(msg.Key),
		Timestamp: time.Now(),
		Headers:   encodeHeaders(msg.Headers),
	}

	partition, offset, err := b.producer.SendMessage(producerMsg)
	if err != nil {
		return fmt.Errorf("%s: %q", SaramaProducerErrorKind, err.Error())
	}

	Log.Info(
		"message", "Mensaje enviado",
		"outMessage", string(msg.Msg),
		"topic", b.topic,
		"partition", partition,
		"offset", offset,
		"dd.trace_id", traceId,
		"dd.span_id", spanId,
	)

	return nil
}

type baseProducerAsync struct {
	topic               string
	producer            sarama.AsyncProducer
	timeoutMilliseconds int
}

// NewBaseSaramaProducerAsync constructor para enviar mensaje
func NewBaseSaramaProducerAsync(topic string, producer sarama.AsyncProducer, timeoutMilliseconds int) MessageProducer {
	return &baseProducerAsync{topic: topic, producer: producer, timeoutMilliseconds: timeoutMilliseconds}

}

func (b *baseProducerAsync) SendMessage(ctx context.Context, msg *ProducerMessage) (err error) {
	cant := len(msg.Msg)

	if cant < 1 {
		return errors.New(InvalidInputProducerErrorKind)
	}

	traceId, spanId := GetDatadogTraceAndSpanFromContext(ctx)

	producerMsg := &sarama.ProducerMessage{
		Topic:     b.topic,
		Value:     sarama.StringEncoder(msg.Msg),
		Key:       sarama.StringEncoder(msg.Key),
		Timestamp: time.Now(),
		Headers:   encodeHeaders(msg.Headers),
	}

	b.producer.Input() <- producerMsg

	go func() {
		for {
			select {
			case producerMsg = <-b.producer.Successes():
				Log.Info(
					"message", "Produce message success",
					"outMessage:", string(msg.Msg),
					"topic", b.topic,
					"partition", producerMsg.Partition,
					"offset", producerMsg.Offset,
					"dd.trace_id", traceId,
					"dd.span_id", spanId,
				)
				return
			case err = <-b.producer.Errors():
				Log.Error(
					"message", err.Error(),
					"errorMessage", SaramaProducerErrorKind,
					"error", err,
					"dd.trace_id", traceId,
					"dd.span_id", spanId,
				)
				return
			case <-time.After(time.Duration(b.timeoutMilliseconds) * time.Millisecond):
				Log.Warn(
					"message", "Timeout de Producer message",
					"dd.trace_id", traceId,
					"dd.span_id", spanId,
				)
				return
			}
		}
	}()

	return nil
}

func encodeHeaders(headers map[string]string) []sarama.RecordHeader {
	var recordHeaders []sarama.RecordHeader = make([]sarama.RecordHeader, 0)

	for k, v := range headers {
		rHeader := sarama.RecordHeader{
			Key:   []byte(k),
			Value: []byte(v),
		}
		recordHeaders = append(recordHeaders, rHeader)
	}

	return recordHeaders
}
