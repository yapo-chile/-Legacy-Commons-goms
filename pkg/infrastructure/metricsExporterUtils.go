package infrastructure

import (
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

func getEventType() string {
	return toSnakeCase(getFuncName(2))
}

func getEntityName() string {
	return toSnakeCase(getEntity(4))
}

var loggerReplacer = regexp.MustCompile(`(_?log(ger)?_?)`)

func getEventName() string {
	loggerName := toSnakeCase(getFuncName(3))
	return loggerReplacer.ReplaceAllString(loggerName, "")
}

func getFuncName(deepness int) string {
	pc, _, _, _ := runtime.Caller(deepness)
	nameFull := runtime.FuncForPC(pc).Name()
	nameEnd := filepath.Ext(nameFull)
	name := strings.TrimPrefix(nameEnd, ".")
	return name
}

var entityRgx = regexp.MustCompile(`\(\*\w+\)`)
var wordsOnly = regexp.MustCompile(`\w+`)

func getEntity(deepness int) string {
	pc, _, _, _ := runtime.Caller(deepness)
	nameFull := runtime.FuncForPC(pc).Name()
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
