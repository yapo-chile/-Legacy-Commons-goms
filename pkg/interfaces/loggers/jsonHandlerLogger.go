package loggers

import (
	"net/http"

	"github.com/Yapo/goutils"

	"github.mpi-internal.com/Yapo/goms/pkg/interfaces/handlers"
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

func (l *jsonHandlerDefaultLogger) LogResponseFromCache(r *http.Request, response *goutils.Response) {
	l.logger.Info("< %s %s %s (%d) (cache)", r.RemoteAddr, r.Method, r.URL, response.Code)
}

func (l *jsonHandlerDefaultLogger) LogErrorSettingCache(r *http.Request, err error) {
	l.logger.Info("> %s %s %s (Error setting cache): %+v", r.RemoteAddr, r.Method, r.URL, err)
}

// MakeJSONHandlerLogger sets up a JsonHandlerLogger instrumented
// via the provided logger
//func MakeJSONHandlerLogger(logger Logger) *jsonHandlerDefaultLogger {
func MakeJSONHandlerLogger(logger Logger) handlers.JSONHandlerLogger {
	return &jsonHandlerDefaultLogger{
		logger: logger,
	}
}
