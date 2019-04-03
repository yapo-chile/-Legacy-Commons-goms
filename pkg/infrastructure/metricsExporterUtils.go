package infrastructure

import (
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

const (
	eventTypeFuncCallDeepness  int = 4
	entityNameFuncCallDeepness int = 6
	eventNameFuncCallDeepness  int = 5
)

func getEventType() string {
	return toSnakeCase(getFuncName(eventTypeFuncCallDeepness))
}

func getEntityName() string {
	return toSnakeCase(getEntity(entityNameFuncCallDeepness))
}

var loggerReplacer = regexp.MustCompile(`(_?log(ger)?_?)`)

func getEventName() string {
	loggerName := toSnakeCase(getFuncName(eventNameFuncCallDeepness))
	return loggerReplacer.ReplaceAllString(loggerName, "")
}

func getFuncName(deepness int) string {
	nameFull := funcName(deepness, 1)
	nameEnd := filepath.Ext(nameFull)
	name := strings.TrimPrefix(nameEnd, ".")
	return name
}

var entityRgx = regexp.MustCompile(`\(\*\w+\)`)
var wordsOnly = regexp.MustCompile(`\w+`)

func getEntity(deepness int) string {
	nameFull := funcName(deepness, 1)
	entityName := wordsOnly.FindString(entityRgx.FindString(nameFull))
	if entityName == "" {
		return wordsOnly.FindString(filepath.Ext(nameFull))
	}
	return entityName
}

var camel = regexp.MustCompile("(^[^A-Z]*|[A-Z]*)([A-Z][^A-Z]+|$)")

func toSnakeCase(s string) string {
	var a []string
	for _, sub := range camel.FindAllStringSubmatch(s, -1) {
		if sub[1] != "" {
			a = append(a, sub[1])
		}
		if sub[2] != "" {
			a = append(a, sub[2])
		}
	}
	return strings.ToLower(strings.Join(a, "_"))
}

// funcName returns the last function name of invocations on the calling goroutine's stack.
// The stack trace can be skipped using deepness parameter. BufferInitialCapacity gives
// the capacity to record in trace stack slice.
func funcName(deepness, bufferInitialCapacity int) string {
	pc := make([]uintptr, bufferInitialCapacity)
	n := runtime.Callers(deepness, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	return frame.Function
}
