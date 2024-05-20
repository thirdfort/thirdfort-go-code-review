package internal

import (
	"context"
	"net/http"

	"github.com/thirdfort/go-slogctx"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// Returns first value from map['key'][]values or empty
func SafeGetValueFromMap(valueMap http.Header, key string) string {
	values, ok := valueMap[key]
	if !ok || len(values) == 0 {
		return ""
	}
	return values[0]
}

func GetContextHeaders(ctx context.Context) map[string][]string {
	headerMap := ctx.Value(HeaderKey)
	hp, ok := headerMap.(http.Header)
	if !ok {
		slogctx.Warn(ctx, "GetContextHeaders: missing headers")
		return make(map[string][]string)
	}

	return hp
}

// To replace strings.Title() without having to use cases etc
func Title(str string) string {
	return cases.Title(language.English).String(str)
}

func StrPtr(str string) *string { return &str }

func PtrStr(str *string) string {
	if str == nil {
		return ""
	}

	return *str
}
