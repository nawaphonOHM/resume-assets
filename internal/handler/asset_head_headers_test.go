package handler

import (
	"maps"
	"strconv"
	"testing"
)

const (
	headTestContentLength      = 42
	headTestIntegerHeaderValue = 7
	headTestFloatValue         = 12
)

func TestValidatedAssetHeadersFiltersUnsafeValuesAndHeaders(t *testing.T) {
	headers := map[string]any{
		handlerETagHeader:          "unchanged",
		handlerContentLengthHeader: int64(headTestContentLength),
		"X-Integer":                int32(headTestIntegerHeaderValue),
		"Bad Header":               "invalid-name",
		"X-Invalid-Value":          "safe\r\nInjected: yes",
		"X-Unsupported":            true,
		"Connection":               "close",
		"Keep-Alive":               "timeout=5",
		"Proxy-Authenticate":       "Basic",
		"Proxy-Authorization":      "Basic token",
		"Proxy-Connection":         "keep-alive",
		"TE":                       "trailers",
		"Trailer":                  "X-Trailer",
		"Transfer-Encoding":        "chunked",
		"Upgrade":                  "websocket",
	}
	want := map[string]string{
		handlerETagHeader:          "unchanged",
		handlerContentLengthHeader: strconv.Itoa(headTestContentLength),
		"X-Integer":                strconv.Itoa(headTestIntegerHeaderValue),
	}
	if got := validatedAssetHeaders(headers); !maps.Equal(got, want) {
		t.Fatalf("validated headers: got %#v, want %#v", got, want)
	}
}

func TestValidatedAssetHeadersRejectsInvalidContentLength(t *testing.T) {
	tests := []struct {
		value any
		name  string
	}{
		{name: "negative integer", value: int64(-1)},
		{name: "negative string", value: "-1"},
		{name: "non numeric string", value: "12 bytes"},
		{name: "unsupported type", value: float64(headTestFloatValue)},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assertContentLengthIsOmitted(t, test.value)
		})
	}
}

func assertContentLengthIsOmitted(t *testing.T, value any) {
	t.Helper()
	got := validatedAssetHeaders(map[string]any{handlerContentLengthHeader: value})
	if len(got) != 0 {
		t.Fatalf("invalid Content-Length was forwarded: %#v", got)
	}
}
