package commons

import (
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
)

//MetricsConfig configuracion de metricas para endpoints
type MetricsConfig struct {
	RequestDuration metrics.Histogram
	RequestCount    metrics.Counter
}

//MakeDefaultEndpointMetrics crea metricas por defecto para endpoint middleware de metricas
func MakeDefaultEndpointMetrics() *MetricsConfig {
	return &MetricsConfig{
		RequestDuration: prometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "api",
			Subsystem: "metrics",
			Name:      "request_duration_seconds",
			Help:      "Duracion de request en segundos.",
		}, []string{"operation", "status"}),
		RequestCount: prometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "api",
			Subsystem: "metrics",
			Name:      "requests",
			Help:      "Total de peticiones.",
		}, []string{"operation", "status"}),
	}
}
