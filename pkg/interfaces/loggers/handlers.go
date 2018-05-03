package loggers

import (
	"net/http"

	"github.com/Yapo/goutils"
)

type jsonHandlerDefaultLogger struct {
	logger Logger
}

func (l *jsonHandlerDefaultLogger) LogRequestStart(r *http.Request) {
	l.logger.Info("< %s %s %s", r.RemoteAddr, r.Method, r.URL)
}

func (l *jsonHandlerDefaultLogger) LogRequestEnd(r *http.Request, response *goutils.Response) {
	l.logger.Info("> %s %s %s (%d)", r.RemoteAddr, r.Method, r.URL, response.Code)
}

func (l *jsonHandlerDefaultLogger) LogRequestPanic(r *http.Request, response *goutils.Response, err interface{}) {
	l.logger.Error("> %s %s %s (%d): %s", r.RemoteAddr, r.Method, r.URL, response.Code, err)
}

// MakeJsonHandlerLogger sets up a JsonHandlerLogger instrumented
// via the provided logger
func MakeJsonHandlerLogger(logger Logger) *jsonHandlerDefaultLogger {
	return &jsonHandlerDefaultLogger{
		logger: logger,
	}
}
