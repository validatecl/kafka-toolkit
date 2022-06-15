package commons

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/go-kit/kit/metrics"

	opentracing "github.com/opentracing/opentracing-go"
)

//MakeDefaultEntryEndpoint endpoint middleware por defecto
func MakeDefaultEntryEndpoint(service, path, method string, logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return EndpointLogMiddleware(service, path, method, logger)(next)
	}
}

// EndpointLogMiddleware loguea tiempo que tomo servir request y error si es que hubo alguno
func EndpointLogMiddleware(service, path, method string, logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		logger = log.WithPrefix(logger, "method", method)
		logger = log.WithPrefix(logger, "path", path)
		logger = log.WithPrefix(logger, "serviceName", service)
		logger = log.WithPrefix(logger, "caller", log.DefaultCaller)

		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			now := time.Now()

			span := opentracing.SpanFromContext(ctx)

			if span != nil {
				logger = log.WithPrefix(logger, "spanContext", span.Context())
			}

			level.Info(logger).Log("started", now)
			level.Debug(logger).Log("request", request)

			defer func(begin time.Time) {
				level.Info(logger).Log("took", time.Since(begin))

				if err != nil {
					level.Info(logger).Log("Result", "NOK")
					level.Error(logger).Log("endpoint_error", err)
				} else {
					level.Info(logger).Log("Result", "OK")
					level.Debug(logger).Log("response", response)
				}
			}(now)
			return next(ctx, request)
		}
	}
}

// EndpointSuccesAndFailureMetricsMiddleware middleware de metrics de success y failures
func EndpointSuccesAndFailureMetricsMiddleware(requests metrics.Counter, operation string, isHTTP bool) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			defer func() {
				status := "SUCCESS"

				if isHTTP {
					status = "2xx"
				}

				if err != nil && isHTTP {
					status = fmt.Sprintf("%d", Error2Status(err))
				} else if err != nil {
					status = "ERROR"
				}

				requests.With("operation", operation, "status", status).Add(1)
			}()

			return next(ctx, request)
		}
	}
}

// EndpointTimeTakenMetricsMiddleware middleware de metrics para tiempo que tomo la peticion
func EndpointTimeTakenMetricsMiddleware(operation string, duration metrics.Histogram) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			defer func(begin time.Time) {
				status := "SUCCESS"
				if err != nil {
					status = "ERROR"
				}

				duration.With("operation", operation, "status", status).Observe(time.Since(begin).Seconds())
			}(time.Now())
			return next(ctx, request)
		}
	}
}

// EndpointHTTPTimeTakenMetricsMiddleware middleware de metrics para tiempo que tomo la peticion HTTP
func EndpointHTTPTimeTakenMetricsMiddleware(operation string, duration metrics.Histogram) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			defer func(begin time.Time) {
				status := "2xx"
				if err != nil {
					status = fmt.Sprintf("%d", Error2Status(err))
				}

				duration.With("operation", operation, "status", status).Observe(time.Since(begin).Seconds())
			}(time.Now())
			return next(ctx, request)
		}
	}
}
