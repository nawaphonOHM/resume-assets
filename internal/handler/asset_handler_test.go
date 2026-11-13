package handler

import (
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/nawaphonOHM/resume-assets/internal/domain"
)

const (
	handlerTestFilename          = "app.js"
	handlerTestContentLength     = 12
	handlerTestAvailableCriteria = `{"status":"available"}`
)

func TestServeAssetValidatesFilename(t *testing.T) {
	for _, filename := range []string{"", ".", "..", "nested/file.js", "bad\x00name"} {
		t.Run(filename, func(t *testing.T) {
			assertInvalidFilename(t, filename)
		})
	}
}

func assertInvalidFilename(t *testing.T, filename string) {
	t.Helper()
	service := &assetServiceFake{}
	response := NewAssetHandler(service).ServeAsset(newHandlerTestContext(context.Background(), filename))
	if response.StatusCode() != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid filename, got %d", response.StatusCode())
	}
	if service.calls != 0 {
		t.Fatalf("invalid filename reached service %d times", service.calls)
	}
}

func TestServeAssetReturnsStream(t *testing.T) {
	service := &assetServiceFake{
		asset: &domain.StaticAsset{
			Body:          handlerBody{},
			ContentType:   "application/javascript",
			ContentLength: handlerTestContentLength,
			ETag:          "asset-tag",
		},
		document: &domain.AssetDocument{},
	}
	response := NewAssetHandler(service).ServeAsset(newHandlerTestContext(context.Background(), handlerTestFilename))
	if response.StatusCode() != http.StatusOK {
		t.Fatalf("expected status 200, got %d", response.StatusCode())
	}
	assertHandlerServiceCall(t, service, handlerTestFilename)
	if service.call.criteria != nil {
		t.Fatalf("expected absent criteria, got %q", *service.call.criteria)
	}
}

func TestServeAssetClosesStreamWhenDocumentIsMissing(t *testing.T) {
	body := &trackedHandlerBody{}
	service := &assetServiceFake{asset: &domain.StaticAsset{Body: body}}
	response := NewAssetHandler(service).ServeAsset(newHandlerTestContext(context.Background(), handlerTestFilename))
	if response.StatusCode() != http.StatusInternalServerError {
		t.Fatalf("expected status 500 for missing metadata document, got %d", response.StatusCode())
	}
	if !body.closed {
		t.Fatal("asset stream was not closed after missing metadata document")
	}
}

func TestServeAssetRejectsQueryBindingErrors(t *testing.T) {
	ctx := newHandlerTestContext(context.Background(), handlerTestFilename)
	ctx.request.queryBindErr = context.Canceled
	service := &assetServiceFake{}
	response := NewAssetHandler(service).ServeAsset(ctx)
	if response.StatusCode() != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", response.StatusCode())
	}
	if service.calls != 0 {
		t.Fatalf("query binding error reached service %d times", service.calls)
	}
}

func assertHandlerServiceCall(t *testing.T, service *assetServiceFake, filename string) {
	t.Helper()
	if service.calls != 1 {
		t.Fatalf("expected one service call, got %d", service.calls)
	}
	assertHandlerFilename(t, service.call, filename)
}

func assertHandlerFilename(t *testing.T, call *handlerServiceCall, filename string) {
	t.Helper()
	if call == nil {
		t.Fatal("asset service request was not recorded")
	}
	if call.filename != filename {
		t.Fatalf("expected filename %q, got %q", filename, call.filename)
	}
}

func newHandlerTestContext(ctx context.Context, filename string) *handlerTestContext {
	return &handlerTestContext{request: &handlerTestRequest{filename: filename, ctx: ctx}}
}

type trackedHandlerBody struct {
	closed bool
}

func (*trackedHandlerBody) Read([]byte) (int, error) { return 0, io.EOF }

func (b *trackedHandlerBody) Close() error {
	b.closed = true
	return nil
}
