package kafka_toolkit

import (
	"errors"

	"github.com/Shopify/sarama"
)

var (
	errInvalidParameters = "Error parametros ingresados son invalidos"
	errAvailableBrokers  = "Error no es posible conectarse a este brokers"
	errBrokerHealth      = "Error brokers activos son menores a la cantidad de particiones definidas"
)

type HealthCheck interface {
	Health() error
}

type healthCheck struct {
	brokers []string
	config  *sarama.Config
}

// NewHealthCheck constructor que implementa dependencias de HealthCheck
func NewHealthCheck(brokers []string, config *sarama.Config) HealthCheck {
	return &healthCheck{
		brokers: brokers,
		config:  config,
	}
}

// Health retorna error si no es posible crear un cliente Kafka
func (h *healthCheck) Health() error {
	if len(h.brokers) < 1 || h.config == nil {
		return errors.New(errInvalidParameters)
	}

	client, err := sarama.NewClient(h.brokers, h.config)
	if err != nil {
		return errors.New(errAvailableBrokers)
	}

	defer func() {
		client.Close()
	}()

	return nil
}
