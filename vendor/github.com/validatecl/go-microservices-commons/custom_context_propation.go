package commons

import (
	"context"
	"net/http"

	"github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

// HTTPToContextOnce http context or span
func HTTPToContextOnce(tracer opentracing.Tracer, operationName string, logger log.Logger) kithttp.RequestFunc {
	return func(ctx context.Context, req *http.Request) context.Context {
		var span opentracing.Span
		wireContext, err := tracer.Extract(
			opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(req.Header),
		)

		if err != nil && err != opentracing.ErrSpanContextNotFound {
			logger.Log("err", err)
		}
		if err != nil && err == opentracing.ErrSpanContextNotFound {
			_, ctx := opentracing.StartSpanFromContext(ctx, operationName)
			return ctx
		}

		span = tracer.StartSpan(operationName, ext.RPCServerOption(wireContext))
		ext.HTTPMethod.Set(span, req.Method)
		ext.HTTPUrl.Set(span, req.URL.String())

		return opentracing.ContextWithSpan(ctx, span)

	}
}
