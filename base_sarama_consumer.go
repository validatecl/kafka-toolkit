package kafka_toolkit

import (
	"context"
	"fmt"

	"github.com/Shopify/sarama"
)

// Implementando interfaz consumer de Sarama (tomado desde ejemplos de Sarama)
// https://github.com/Shopify/sarama/blob/master/examples/consumergroup/main.go
// TODO: Identificar como implementar correctamente metodos de Sarama consumer

// BaseConsumer estructura que representa el consumer de Kafka
type BaseConsumer struct {
	Ready          chan bool
	MessageHandler MessageHandler
	ErrorHandler   ConsumerErrorHandler
}

const errorSaramaMessage = "sarama message is nil"

// NewBaseConsumer construye un nuevo consumer base
func NewBaseConsumer(handler MessageHandler, errorHandler ConsumerErrorHandler) BaseConsumer {
	return BaseConsumer{MessageHandler: handler, ErrorHandler: errorHandler, Ready: make(chan bool)}
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (consumer *BaseConsumer) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(consumer.Ready)
	return nil
}

// Cleanup Realiza clean up de sarama.
func (consumer *BaseConsumer) Cleanup(sess sarama.ConsumerGroupSession) error {
	//TODO: verificar que se debe hacer en cleanup
	return nil
}

// ConsumeClaim inicia loop para cobrar mensajes
func (consumer *BaseConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		go logClaims(session)

		select {
		case message := <-claim.Messages():
			consumer.HandleSaramaMessage(session, message)
		case <-session.Context().Done():
			return nil
		}
	}
}

func logClaims(session sarama.ConsumerGroupSession) {
	for _, v := range session.Claims() {
		Log.Info(fmt.Sprintf("Claims Partition Ids: %v \n", v))
	}
}

//HandleSaramaMessage maneja el mensaje Sarama
func (consumer *BaseConsumer) HandleSaramaMessage(session sarama.ConsumerGroupSession, saramaMessage *sarama.ConsumerMessage) {
	if saramaMessage == nil {
		Log.Error("error", errorSaramaMessage)
		return
	}

	msg := saramaToGenericMessage(saramaMessage)

	if err := consumer.MessageHandler.HandleMessage(context.Background(), msg); err != nil {
		consumer.ErrorHandler.HandleError(msg.Msg, err)
	}

	session.MarkMessage(saramaMessage, "")
}

func saramaToGenericMessage(msg *sarama.ConsumerMessage) *ConsumerMessage {
	return &ConsumerMessage{
		Headers:   saramaHeaderToMap(msg.Headers),
		Msg:       msg.Value,
		Key:       msg.Key,
		Offset:    msg.Offset,
		Partition: msg.Partition,
		Timestamp: msg.Timestamp,
		Topic:     msg.Topic,
	}
}

func saramaHeaderToMap(saramaHeaders []*sarama.RecordHeader) map[string]string {
	headers := make(map[string]string, 0)

	for _, header := range saramaHeaders {
		headers[string(header.Key)] = string(header.Value)
	}

	return headers
}

// Close implementa metodo close de sarama.Consumer
func (consumer *BaseConsumer) Close() error {
	return nil
}
