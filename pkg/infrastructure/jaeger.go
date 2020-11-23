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
			Status:         http.StatusOK,
		}

		ctx := opentracing.ContextWithSpan(r.Context(), serverSpan)
		handler(recorder, r.WithContext(ctx))
		serverSpan.SetTag("http.status_code", recorder.Status)
	})
}

type JaegerTracedHTTPHandler struct {
	httpHandler repository.HTTPHandler
}

// NewJaegerTracedHTTPHandler will create a new instance of a custom http request handler
func NewJaegerTracedHTTPHandler(h repository.HTTPHandler) repository.HTTPHandler {
	return &JaegerTracedHTTPHandler{
		httpHandler: h,
	}
}

// Send will execute the sending of a http request
// but in this case it will wait until it obtains a successful response
// in order to continue it's execution
func (h *JaegerTracedHTTPHandler) Send(req repository.HTTPRequest) (repository.HTTPResponse, error) {
	var name = req.GetMethod() + " " + req.GetPath()
	span, _ := opentracing.StartSpanFromContext(req.Context(), name)
	defer span.Finish()

	ext.SpanKindRPCClient.Set(span)
	ext.HTTPUrl.Set(span, req.GetPath())
	ext.HTTPMethod.Set(span, req.GetMethod())
	span.Tracer().Inject(span.Context(), // nolint:errcheck
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(req.GetHeaders()),
	)

	response, err := h.httpHandler.Send(req)

	span.SetTag("http.status_code", response.GetStatusCode())
	return response, err
}

// NewRequest returns an initialized struct that can be used to make a http request
func (h *JaegerTracedHTTPHandler) NewRequest(ctx context.Context) repository.HTTPRequest {
	return h.httpHandler.NewRequest(ctx)
}
