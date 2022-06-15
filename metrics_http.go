package kafka_toolkit

import (
	"context"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	kitlog "github.com/go-kit/log"
	commons "github.com/validatecl/go-microservices-commons"
)

// MakeHealthHandlerBuilder Crea Handler builder para health check
func MakeHealthHandlerBuilder(logger kitlog.Logger, service HealthCheck) commons.HTTPHandlerBuilder {
	decoder := func(context.Context, *http.Request) (request interface{}, err error) { return nil, nil }

	endpoint := MakeServiceHealthCheckEndpoint(service)

	endpointCfgs := []commons.EndpointConfig{
		commons.GET("/healthz", "HEALTHZ", endpoint, decoder, EncodeResponse),
	}

	return commons.MakeHTTPHandlerBuilder(logger, endpointCfgs)
}

// MakeServiceHealthCheckEndpoint Registra la URI /healthz para el healthckeck
func MakeServiceHealthCheckEndpoint(service HealthCheck) endpoint.Endpoint {
	return func(_ context.Context, _ interface{}) (interface{}, error) {
		return nil, service.Health()
	}
}

// EncodeResponse hola
func EncodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	if response == nil {
		w.WriteHeader(http.StatusOK)
	} else if f, ok := response.(endpoint.Failer); ok && f.Failed() != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		return nil
	}

	return nil
}
