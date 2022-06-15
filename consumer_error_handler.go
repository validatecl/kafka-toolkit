package kafka_toolkit

//ConsumerErrorHandler Interface
type ConsumerErrorHandler interface {
	HandleError(messageVal []byte, err error) error
}
