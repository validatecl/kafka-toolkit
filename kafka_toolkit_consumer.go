package kafka_toolkit

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/Shopify/sarama"
)

// KafkaConsumer interfaz para iniciar kafka consumer
type KafkaConsumer interface {
	Start() error
}

// SaramaConsumerBuilder builder de sarama consumer
type SaramaConsumerBuilder interface {
	WithErrorHandler(ConsumerErrorHandler) SaramaConsumerBuilder
	Build() (KafkaConsumer, error)
}

type saramaConsumerBuilder struct {
	consumerCfg  ConsumerGroupInput
	msgHandler   MessageHandler
	errorHandler ConsumerErrorHandler
}

// MakeSaramaConsumerBuilder consumer builder
func MakeSaramaConsumerBuilder(cfg ConsumerGroupInput, handler MessageHandler) SaramaConsumerBuilder {
	return &saramaConsumerBuilder{
		consumerCfg:  cfg,
		msgHandler:   handler,
		errorHandler: NewLoggingConsumerErrorHandler(),
	}
}

func (b *saramaConsumerBuilder) WithErrorHandler(errorHandler ConsumerErrorHandler) SaramaConsumerBuilder {
	b.errorHandler = errorHandler
	return b
}

func (b *saramaConsumerBuilder) Build() (KafkaConsumer, error) {
	conf, consumer, err := createBaseConsumer(b.consumerCfg, b.msgHandler, b.errorHandler)
	if err != nil {
		Log.Error("Error generando configuracion:", err)
		return nil, err
	}

	client, err := sarama.NewConsumerGroup(conf.Brokers, conf.Group, conf.SaramaConfig)
	if err != nil {
		Log.Error("Error creando cliente para consumer group:", err)
		return nil, err
	}

	return &saramaKafkaConsumer{
		conf:     conf,
		consumer: consumer,
		client:   client,
	}, nil
}

type saramaKafkaConsumer struct {
	conf     *ConsumerGroupConfig
	consumer BaseConsumer
	client   sarama.ConsumerGroup
}

//StartConsumer Inicializa consumo de topico Kafka
func (s *saramaKafkaConsumer) Start() error {

	ctx, cancel := context.WithCancel(context.Background())

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if err := s.client.Consume(ctx, []string{s.conf.Topic}, &s.consumer); err != nil {
				Log.Error(
					"errorMessage", "Error de consumer",
					"error", err.Error())
			}
			// Checkeando si el context ha sido cancelado
			if ctx.Err() != nil {
				return
			}
			s.consumer.Ready = make(chan bool)
		}
	}()

	<-s.consumer.Ready // Esperamos a que el consumer este configurado
	Log.Info("message", "Sarama Consumer inicializado!")

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-ctx.Done():
		Log.Info("message", "Terminando: Contexto cancelado")
	case <-sigterm:
		Log.Warn("message", "Terminando: por signal")
	}

	cancel()
	wg.Wait()
	if err := s.client.Close(); err != nil {
		Log.Error(
			"errorMessage", "Error cerrando cliente",
			"error", err)
	}

	return nil
}

func createBaseConsumer(consumerCfg ConsumerGroupInput, msgHandler MessageHandler, errorHandler ConsumerErrorHandler) (*ConsumerGroupConfig, BaseConsumer, error) {
	balanceStrategyResolver := NewBalanceStrategyResolver()
	configurer := NewSaramaConsumerConfigurer(balanceStrategyResolver)

	Log.Info(
		"message", "Inicializando consumer",
		"topic", consumerCfg.Topic,
		"group", consumerCfg.Group,
		"client_id", consumerCfg.ClientID)

	conf, err := configurer.GenerateConfig(consumerCfg)

	if err != nil {
		Log.Error(
			"errorMessage", "Error generando configuracion:",
			"error", err)

		return nil, BaseConsumer{}, err
	}

	consumer := NewBaseConsumer(msgHandler, errorHandler)

	return conf, consumer, err
}
