package commons

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/tracing/opentracing"
	httptransport "github.com/go-kit/kit/transport/http"
	stdopentracing "github.com/opentracing/opentracing-go"
)

// HTTPClientBuilder builder de http client
type HTTPClientBuilder interface {
	Build() endpoint.Endpoint
	WithTracer(tracer stdopentracing.Tracer) HTTPClientBuilder
}

type httpClientBuilder struct {
	method         string
	uri            string
	timeout        time.Duration
	tracer         stdopentracing.Tracer
	encodeRequest  httptransport.EncodeRequestFunc
	decodeResponse httptransport.DecodeResponseFunc
	logger         log.Logger
}

func (h *httpClientBuilder) WithTracer(tracer stdopentracing.Tracer) HTTPClientBuilder {
	h.tracer = tracer
	return h
}

func (h *httpClientBuilder) Build() endpoint.Endpoint {
	url, _ := url.Parse(h.uri)

	opts := []httptransport.ClientOption{
		httptransport.SetClient(&http.Client{Timeout: h.timeout}),
	}

	if h.tracer != nil {
		opts = append(opts, httptransport.ClientBefore(opentracing.ContextToHTTP(h.tracer, h.logger)))
	}

	return httptransport.NewClient(
		h.method,
		url,
		h.encodeRequest,
		h.decodeResponse,
		opts...,
	).Endpoint()
}

// MakeHTTPClientBuilder crea un nuevo http client builder
func MakeHTTPClientBuilder(method, uri string,
	timeout time.Duration,
	encodeRequest httptransport.EncodeRequestFunc,
	decodeResponse httptransport.DecodeResponseFunc,
	logger log.Logger) HTTPClientBuilder {
	return &httpClientBuilder{
		method:         method,
		uri:            uri,
		timeout:        timeout,
		encodeRequest:  encodeRequest,
		decodeResponse: decodeResponse,
		logger:         logger,
	}
}

// DefaultRequestEncode encode de request de client
func DefaultRequestEncode(ctx context.Context, r *http.Request, request interface{}) error {
	r.Header.Add("Content-Type", "application/json")

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return err
	}

	r.Body = ioutil.NopCloser(&buf)
	return nil
}

// DefaultDecodeResponse decode de response por defecto, para utilizar referenciar dependencias
func DefaultDecodeResponse(_ context.Context, r *http.Response, response interface{}) (interface{}, error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)

	if r.StatusCode != http.StatusOK {
		responseBody := buf.String()
		return nil,
			&GatewayError{
				Message: fmt.Sprintf("Error llamando a service o api status: %d: %v", r.StatusCode, responseBody),
			}
	}

	if err := json.Unmarshal(buf.Bytes(), response); err != nil {
		return nil, err
	}

	return response, nil
}
