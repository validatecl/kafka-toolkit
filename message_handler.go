package kafka_toolkit

import (
	"context"
)

//MessageHandler interfaz para manejo de mensajes
type MessageHandler interface {
	HandleMessage(ctx context.Context, inMsg *ConsumerMessage) error
}
