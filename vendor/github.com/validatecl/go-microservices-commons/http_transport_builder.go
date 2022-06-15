package commons

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/go-kit/kit/tracing/opentracing"
	"github.com/go-kit/kit/transport"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	httptransport "github.com/go-kit/kit/transport/http"
	stdopentracing "github.com/opentracing/opentracing-go"
)

//EndpointConfig entrada de endpoint
type EndpointConfig struct {
	Method          string
	Path            string
	Service         string
	Endpoint        endpoint.Endpoint
	RequestDecoder  httptransport.DecodeRequestFunc
	ResponseEncoder httptransport.EncodeResponseFunc
	QueryParams     []string
}

// Crea un endpoint config, si response encoder no es provisto utiliza el valor por defecto
func newEndpointConfig(method, path, service string, ep endpoint.Endpoint, reqDecoder httptransport.DecodeRequestFunc, respEncoder httptransport.EncodeResponseFunc) EndpointConfig {
	if respEncoder == nil {
		respEncoder = MakeDefaultEncodeHTTPResponseFunc()
	}

	return EndpointConfig{
		Method:          method,
		Path:            path,
		Service:         service,
		Endpoint:        ep,
		RequestDecoder:  reqDecoder,
		ResponseEncoder: respEncoder,
	}
}

//GET crea un GET endpoint config - Parametro Response Encoder es opcional
func GET(path, service string, ep endpoint.Endpoint, reqDecoder httptransport.DecodeRequestFunc, respEncoder httptransport.EncodeResponseFunc) EndpointConfig {
	return newEndpointConfig("GET", path, service, ep, reqDecoder, respEncoder)
}

//POST crea un POST endpoint config - Parametro Response Encoder es opcional
func POST(path, service string, ep endpoint.Endpoint, reqDecoder httptransport.DecodeRequestFunc, respEncoder httptransport.EncodeResponseFunc) EndpointConfig {
	return newEndpointConfig("POST", path, service, ep, reqDecoder, respEncoder)
}

//PUT crea un PUT endpoint config - Parametro Response Encoder es opcional
func PUT(path, service string, ep endpoint.Endpoint, reqDecoder httptransport.DecodeRequestFunc, respEncoder httptransport.EncodeResponseFunc) EndpointConfig {
	return newEndpointConfig("PUT", path, service, ep, reqDecoder, respEncoder)
}

//PATCH crea un PATCH endpoint config - Parametro Response Encoder es opcional
func PATCH(path, service string, ep endpoint.Endpoint, reqDecoder httptransport.DecodeRequestFunc, respEncoder httptransport.EncodeResponseFunc) EndpointConfig {
	return newEndpointConfig("PATCH", path, service, ep, reqDecoder, respEncoder)
}

//DELETE crea un DELETE endpoint config
func DELETE(path, service string, ep endpoint.Endpoint, reqDecoder httptransport.DecodeRequestFunc, respEncoder httptransport.EncodeResponseFunc) EndpointConfig {
	return newEndpointConfig("DELETE", path, service, ep, reqDecoder, respEncoder)
}

// MakeDefaultEncodeHTTPResponseFunc default http response encoder
func MakeDefaultEncodeHTTPResponseFunc() func(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return MakeEncodeHTTPResponseFunc(Error2Wrapper)
}

// WithQueryParams setea query params
func (c EndpointConfig) WithQueryParams(queryParams []string) EndpointConfig {
	c.QueryParams = queryParams
	return c
}

//HTTPHandlerBuilder builder para httpServer
type HTTPHandlerBuilder interface {
	WithCustomErrorWrapper(err2wFunc Error2WrapperFunc) HTTPHandlerBuilder
	WithTracer(otTracer stdopentracing.Tracer) HTTPHandlerBuilder
	WithMetrics(metricsCfg *MetricsConfig) HTTPHandlerBuilder
	WithCustomRequestFuncMiddleware(customRequestFuncMiddleware RequestFuncMiddleware) HTTPHandlerBuilder
	Build() http.Handler
}

type httpHandlerBuilder struct {
	logger                      log.Logger
	endpoints                   []EndpointConfig
	err2wFunc                   Error2WrapperFunc
	openTracer                  stdopentracing.Tracer
	metricsCfg                  *MetricsConfig
	customRequestFuncMiddleware RequestFuncMiddleware
}

// RequestFuncMiddleware Middleware de funcion para actualizar context a partir de una request http
type RequestFuncMiddleware func(httptransport.RequestFunc) httptransport.RequestFunc

//MakeHTTPHandlerBuilder crea un http server builder
func MakeHTTPHandlerBuilder(logger log.Logger, endpoints []EndpointConfig) HTTPHandlerBuilder {
	return &httpHandlerBuilder{
		logger:    logger,
		endpoints: endpoints,
	}
}

func (b *httpHandlerBuilder) WithCustomErrorWrapper(err2wFunc Error2WrapperFunc) HTTPHandlerBuilder {
	b.err2wFunc = err2wFunc
	return b
}

func (b *httpHandlerBuilder) WithTracer(otTracer stdopentracing.Tracer) HTTPHandlerBuilder {
	b.openTracer = otTracer
	return b
}

func (b *httpHandlerBuilder) WithMetrics(metricsCfg *MetricsConfig) HTTPHandlerBuilder {
	b.metricsCfg = metricsCfg
	return b
}

func (b *httpHandlerBuilder) WithCustomRequestFuncMiddleware(customRequestFuncMiddleware RequestFuncMiddleware) HTTPHandlerBuilder {
	b.customRequestFuncMiddleware = customRequestFuncMiddleware
	return b
}

//InitHTTPServer inicializa un http server
func (b *httpHandlerBuilder) Build() http.Handler {

	var errorEncoder httptransport.ErrorEncoder

	//Si hay un error2Wrapper function personalizado se utiliza ese en lugar del por defecto
	if b.err2wFunc != nil {
		errorEncoder = MakeServerErrorEncoderFunc(b.err2wFunc)
	} else {
		errorEncoder = MakeDefaultServerErrorEncoderFunc()
	}

	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(errorEncoder),
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(b.logger)),
	}

	m := mux.NewRouter()

	b.initEndpoints(m, options)

	if b.metricsCfg != nil {
		m.Methods("GET").Path("/metrics").Handler(promhttp.Handler())
	}

	return m
}

func (b *httpHandlerBuilder) initEndpoints(m *mux.Router, options []httptransport.ServerOption) {
	for _, epCfg := range b.endpoints {

		var ep endpoint.Endpoint
		ep = MakeDefaultEntryEndpoint(epCfg.Service, epCfg.Path, epCfg.Method, b.logger)(epCfg.Endpoint)

		if b.openTracer != nil {
			ep = opentracing.TraceServer(b.openTracer, epCfg.Service)(ep)

		}

		if b.metricsCfg != nil {
			ep = EndpointSuccesAndFailureMetricsMiddleware(b.metricsCfg.RequestCount, epCfg.Service, true)(ep)
			ep = EndpointHTTPTimeTakenMetricsMiddleware(epCfg.Service, b.metricsCfg.RequestDuration)(ep)
		}

		tracingOpts := b.resolveTracingOptions(epCfg, options)

		handler := httptransport.NewServer(
			ep,
			epCfg.RequestDecoder,
			epCfg.ResponseEncoder,
			tracingOpts...,
		)

		m.Methods(epCfg.Method).Path(epCfg.Path).Handler(handler).Queries(epCfg.QueryParams...)
		m.Methods(epCfg.Method).Path(epCfg.Path).Handler(handler)

		level.Info(b.logger).Log(epCfg.Path, " incializado...")
	}
}

func (b *httpHandlerBuilder) resolveTracingOptions(epCfg EndpointConfig, options []httptransport.ServerOption) []httptransport.ServerOption {
	if b.openTracer != nil {
		operationName := fmt.Sprintf("%s-http", epCfg.Service)
		beforRequestFunc := HTTPToContextOnce(b.openTracer, operationName, b.logger)
		if b.customRequestFuncMiddleware != nil {
			beforRequestFunc = b.customRequestFuncMiddleware(beforRequestFunc)
		}
		return append(options, httptransport.ServerBefore(beforRequestFunc))
	}

	return options
}
