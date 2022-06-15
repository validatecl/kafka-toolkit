package kafka_toolkit

import "fmt"

type loggingConsumerErrorHandler struct {
}

//NewLoggingConsumerErrorHandler Constructor
func NewLoggingConsumerErrorHandler() ConsumerErrorHandler {
	return &loggingConsumerErrorHandler{}
}

func (s *loggingConsumerErrorHandler) HandleError(msg []byte, err error) error {
	Log.Error(fmt.Sprintf("Se produjo el siguiente error: %v procesando el mensaje %v ", err, string(msg)))
	return nil
}
