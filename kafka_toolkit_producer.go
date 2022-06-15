package kafka_toolkit

import "github.com/Shopify/sarama"

// NewSimpleSyncProducer Crea un nuevo simple producer
func NewSimpleSyncProducer(configInput BaseProducerConfigInput) (MessageProducer, error) {
	saramaProducer, err := newProducer(configInput)

	if err != nil {
		return nil, err
	}

	return NewBaseSaramaProducer(configInput.Topic, saramaProducer), nil
}

func newProducer(configInput BaseProducerConfigInput) (sarama.SyncProducer, error) {

	conf, err := NewBaseProducerConfigurer().GenerateConfig(configInput)

	if err != nil {
		return nil, err
	}

	producer, err := sarama.NewSyncProducer(conf.Brokers, conf.SaramaConfig)

	return producer, err
}

func NewSimpleAsyncProducer(configInput BaseProducerConfigInput) (MessageProducer, error) {
	saramaAsyncProducer, err := newAsyncProducer(configInput)

	if err != nil {
		return nil, err
	}

	if configInput.TimeoutMs == 0 {
		configInput.TimeoutMs = 1000
	}

	return NewBaseSaramaProducerAsync(configInput.Topic, saramaAsyncProducer, int(configInput.TimeoutMs)), nil
}

func newAsyncProducer(configInput BaseProducerConfigInput) (sarama.AsyncProducer, error) {

	conf, err := NewBaseProducerConfigurer().GenerateConfig(configInput)

	if err != nil {
		return nil, err
	}

	producer, err := sarama.NewAsyncProducer(conf.Brokers, conf.SaramaConfig)

	return producer, err
}
