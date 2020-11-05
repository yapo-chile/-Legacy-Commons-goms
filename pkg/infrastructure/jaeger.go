package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"

	opentracing "github.com/opentracing/opentracing-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	jaegermetrics "github.com/uber/jaeger-lib/metrics"
	//	"github.com/uber/jaeger-client-go"

	"github.mpi-internal.com/Yapo/goms/pkg/interfaces/repository"
)

// InitJaegerTracing sets up the global tracer and returns it along with an io.Closer which must
// be shut down before the application finishes
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
	tracer, closer, err := cfg.NewTracer(
		jaegercfg.Logger(jLogger),
		jaegercfg.Metrics(jMetricsFactory))

	if err != nil {
		return nil, nil, fmt.Errorf("could not initialize Jaeger tracer: %s", err.Error())
	}

	opentracing.SetGlobalTracer(tracer)
	return tracer, closer, nil
}

type statusRecorder struct {
	http.ResponseWriter
	Status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.Status = status
	r.ResponseWriter.WriteHeader(status)
}

// JaegerWrapperFunc wraps a http.HandlerFunc with Jaeger Tracing
func JaegerWrapperFunc(pattern string, handler http.HandlerFunc) http.HandlerFunc {
	tracer := opentracing.GlobalTracer()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract the context from the headers
		spanCtx, _ := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
		serverSpan := tracer.StartSpan(pattern, ext.RPCServerOption(spanCtx))
		defer serverSpan.Finish()

		recorder := &statusRecorder{
			ResponseWriter: w,
			Status:         200,
		}

		ctx := opentracing.ContextWithSpan(context.Background(), serverSpan)
		handler(recorder, r.WithContext(ctx))
		serverSpan.SetTag("http.status_code", recorder.Status)
	})
}

// TraceRequest generates a new traced HTTP request with opentracing headers injected into it
func TraceRequest(req repository.HTTPRequest) opentracing.Span {
	var name = req.GetMethod() + " " + req.GetPath()
	span, _ := opentracing.StartSpanFromContext(req.Context(), name)
	span.SetTag("hello-to", "+Estables")
	span.SetBaggageItem("le_request", "is also here")
	span.LogFields(
		log.String("le log", "string type"),
		log.Int("another log", 42),
	)

	ext.SpanKindRPCClient.Set(span)
	ext.HTTPUrl.Set(span, req.GetPath())
	ext.HTTPMethod.Set(span, req.GetMethod())
	span.Tracer().Inject(span.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(req.GetHeaders()),
	)
	return span
}
