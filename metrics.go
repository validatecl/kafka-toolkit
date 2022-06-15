package kafka_toolkit

import (
	"fmt"

	"github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	commons "github.com/validatecl/go-microservices-commons"
)

const kafkaHandlerSubsystem = "kafka_handler"

// MakeDefaultKafkaMetrics metricas de kafka por defecto
func MakeDefaultKafkaMetrics(serviceName string) *commons.MetricsConfig {
	return &commons.MetricsConfig{
		RequestDuration: prometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: serviceName,
			Subsystem: kafkaHandlerSubsystem,
			Name:      "handle_duration_seconds",
			Help:      "Duracion de request en segundos.",
		}, []string{"operation", "status"}),
		RequestCount: prometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: serviceName,
			Subsystem: kafkaHandlerSubsystem,
			Name:      "handle_count",
			Help:      "Contador de request",
		}, []string{"operation", "status"}),
	}
}

func MakeKafkaConsumerMetrics(serviceName string, consumerName string) *commons.MetricsConfig {
	return &commons.MetricsConfig{
		RequestDuration: prometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: serviceName,
			Subsystem: kafkaHandlerSubsystem,
			Name:      fmt.Sprintf("consumer_%s_duration_seconds", consumerName),
			Help:      "Duracion de en la fase de consumo en segundos.",
		}, []string{"operation", "status"}),
		RequestCount: prometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: serviceName,
			Subsystem: kafkaHandlerSubsystem,
			Name:      fmt.Sprintf("consumer_%s_handle_count", consumerName),
			Help:      "Contador de request, consumer",
		}, []string{"operation", "status"}),
	}
}

func MakeKafkaProducerMetrics(serviceName string, producerName string) *commons.MetricsConfig {
	return &commons.MetricsConfig{
		RequestDuration: prometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: serviceName,
			Subsystem: kafkaHandlerSubsystem,
			Name:      fmt.Sprintf("producer_%s_duration_seconds", producerName),
			Help:      "Duracion de en la fase de producir mensajes en segundos.",
		}, []string{"operation", "status"}),
		RequestCount: prometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: serviceName,
			Subsystem: kafkaHandlerSubsystem,
			Name:      fmt.Sprintf("producer_%s_handle_count", producerName),
			Help:      "Contador de request, producer",
		}, []string{"operation", "status"}),
	}
}
