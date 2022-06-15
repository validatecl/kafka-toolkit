package kafka_toolkit

// MakeStreamer crea un consumer de tipo streamer
func MakeStreamerBuilder(consumerCfg ConsumerGroupInput, processor StreamProcessor, producerCfg BaseProducerConfigInput, asyncProduce bool) (builder SaramaConsumerBuilder, err error) {
	var producer MessageProducer

	if asyncProduce {
		producer, err = NewSimpleAsyncProducer(producerCfg)
		if err != nil {
			return nil, err
		}
	} else {
		producer, err = NewSimpleSyncProducer(producerCfg)
		if err != nil {
			return nil, err
		}
	}

	handler := NewStreamHandler(producer, processor)

	return MakeSaramaConsumerBuilder(consumerCfg, handler), nil
}
