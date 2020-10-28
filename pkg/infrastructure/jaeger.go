package infrastructure

import (
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
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	return tracer, closer, nil
}

func JaegerWrapperFunc(pattern string, handler http.HandlerFunc) http.HandlerFunc {
	tracer := opentracing.GlobalTracer()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract the context from the headers
		spanCtx, _ := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
		serverSpan := tracer.StartSpan("handler", ext.RPCServerOption(spanCtx))
		defer serverSpan.Finish()
		handler(w, r)
	})
}
