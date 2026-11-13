package handler

import (
	"context"
	"reflect"
	"testing"

	"github.com/nawaphonOHM/resume-assets/internal/domain"
)

type assetHeaderWarning struct {
	message string
	args    []any
}

type assetHeaderWarningRecorder struct {
	warnings []assetHeaderWarning
	contexts []context.Context
}

func (r *assetHeaderWarningRecorder) log(ctx context.Context, message string, args ...any) {
	r.warnings = append(r.warnings, assetHeaderWarning{message: message, args: args})
	r.contexts = append(r.contexts, ctx)
}

func assetResponseHeadersForTest(
	filename string,
	asset *domain.StaticAsset,
	document *domain.AssetDocument,
	warn func(context.Context, string, ...any),
) map[string]string {
	responseContext := assetHeaderResponseContext{
		ctx: context.Background(), filename: filename, warn: warn,
	}
	return assetResponseHeaders(responseContext, asset, document)
}

func assertResponseHeaders(t *testing.T, got, want map[string]string) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("response headers: got %#v, want %#v", got, want)
	}
}

func assertResponseHeaderAbsent(t *testing.T, headers map[string]string, name string) {
	t.Helper()
	if _, exists := headers[name]; exists {
		t.Fatalf("response header %s should be absent: %#v", name, headers)
	}
}

func assertWarningCount(t *testing.T, warnings []assetHeaderWarning, expected int) {
	t.Helper()
	if len(warnings) != expected {
		t.Fatalf("warning count: got %d, want %d (%#v)", len(warnings), expected, warnings)
	}
}

func assertConflictWarning(t *testing.T, warning assetHeaderWarning, header, source string) {
	t.Helper()
	if warning.message != assetHeaderConflictWarning {
		t.Fatalf("unexpected conflict warning: %q", warning.message)
	}
	assertWarningAttributes(t, warning, header, source)
}

func assertWarningAttributes(t *testing.T, warning assetHeaderWarning, header, source string) {
	t.Helper()
	wantArgs := []any{"filename", assetResponseTestFilename, "header", header, "source", source}
	if !reflect.DeepEqual(warning.args, wantArgs) {
		t.Fatalf("warning attributes: got %#v, want %#v", warning.args, wantArgs)
	}
}
