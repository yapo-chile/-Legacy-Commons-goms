package infrastructure

import (
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

func getEventType() string {
	return toSnakeCase(getFuncName(4))
}

func getEntityName() string {
	return toSnakeCase(getEntity(6))
}

var loggerReplacer = regexp.MustCompile(`(_?log(ger)?_?)`)

func getEventName() string {
	loggerName := toSnakeCase(getFuncName(5))
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

func funcName(deepness, maxDeepness int) string {
	pc := make([]uintptr, maxDeepness)
	n := runtime.Callers(deepness, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	return frame.Function
}
