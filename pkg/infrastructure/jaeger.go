package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/opentracing/opentracing-go/ext"

	opentracing "github.com/opentracing/opentracing-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	jaegermetrics "github.com/uber/jaeger-lib/metrics"
	//	"github.com/uber/jaeger-client-go"

	"github.mpi-internal.com/Yapo/goms/pkg/interfaces/repository"
)

func InitJaegerTracing() (opentracing.Tracer, io.Closer, error) {
	cfg, err := jaegercfg.FromEnv()
	if err != nil {
		// parsing errors might happen here, such as when we get a string where we expect a number
		return nil, nil, fmt.Errorf("could not parse Jaeger env vars: %s", err.Error())
	}

	if jconf, err := json.MarshalIndent(cfg, "", "    "); err == nil {
		fmt.Printf("JaegerConfig: \n%s\n", jconf)
	}

	jLogger := jaegerlog.StdLogger
	jMetricsFactory := jaegermetrics.NullFactory

	// Initialize tracer with a logger and a metrics factory
	tracer, closer, err := cfg.NewTracer(
		jaegercfg.Logger(jLogger),
		jaegercfg.Metrics(jMetricsFactory),
	)

	if err != nil {
		return nil, nil, fmt.Errorf("could not initialize jaeger tracer: %s", err.Error())
	}
	opentracing.SetGlobalTracer(tracer)

	return tracer, closer, nil
}

func JaegerWrapperFunc(pattern string, handler http.HandlerFunc) http.HandlerFunc {
	tracer := opentracing.GlobalTracer()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract the context from the headers
		spanCtx, _ := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
		serverSpan := tracer.StartSpan(pattern, ext.RPCServerOption(spanCtx))
		defer serverSpan.Finish()

		ctx := opentracing.ContextWithSpan(context.Background(), serverSpan)
		handler(w, r.WithContext(ctx))
	})
}

/*
func TrackOutgoingCall(req repository.HTTPRequest) {
	h.logger.Debug("Http - %s - Sending HTTP request to: %+v", req.GetMethod(), req.GetPath())
	span := tracer.StartSpan("format", ext.RPCServerOption(spanCtx))

}
*/
// NewTracedRequest generates a new traced HTTP request with opentracing headers injected into it
func TraceRequest(req repository.HTTPRequest) opentracing.Span {
	tracer := opentracing.GlobalTracer()
	span := tracer.StartSpan("call " + req.GetMethod() + " " + req.GetPath())
	span.SetBaggageItem("le_request", "is also here")

	ext.SpanKindRPCClient.Set(span)
	ext.HTTPUrl.Set(span, req.GetPath())
	ext.HTTPMethod.Set(span, req.GetMethod())
	span.Tracer().Inject(span.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(req.GetHeaders()),
	)
	return span
}
