package handler

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/nawaphonOHM/resume-assets/internal/domain"
)

func TestServeAssetHeadMapsErrorsWithoutMetadataHeaders(t *testing.T) {
	tests := []handlerErrorCase{
		{status: http.StatusNotFound, err: domain.ErrAssetNotFound, name: "missing"},
		{status: http.StatusForbidden, err: domain.ErrAssetAccessDenied, name: "denied"},
		{status: http.StatusBadRequest, err: domain.ErrInvalidAssetCriteria, name: "invalid criteria"},
		{status: http.StatusConflict, err: domain.ErrAmbiguousAssetMatch, name: "ambiguous match"},
		{status: http.StatusInternalServerError, err: errors.New("private backend detail"), name: "backend failure"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assertHeadErrorResponse(t, test)
		})
	}
}

func assertHeadErrorResponse(t *testing.T, test handlerErrorCase) {
	t.Helper()
	service := &assetServiceFake{
		err: test.err,
		document: &domain.AssetDocument{
			Filename: handlerTestFilename,
			Access:   handlerAllowedAccess,
			Metadata: domain.AssetMetadata{ForCaching: map[string]any{handlerETagHeader: "private-tag"}},
		},
	}
	requestContext := newHandlerTestContext(context.Background(), handlerTestFilename)
	response := NewAssetHandler(service).ServeAssetHead(requestContext)
	if response.StatusCode() != test.status {
		t.Fatalf("HEAD error response: got status=%d, want %d", response.StatusCode(), test.status)
	}
	assertHandlerServiceCall(t, service, handlerTestFilename)
}

func TestServeAssetHeadUsesGETRequestValidation(t *testing.T) {
	assertInvalidHeadRequest(t, "nested/file.js", nil)
	assertInvalidHeadRequest(
		t,
		handlerTestFilename,
		map[string][]string{"criteria": {handlerTestAvailableCriteria}},
	)
}

func assertInvalidHeadRequest(t *testing.T, filename string, query map[string][]string) {
	t.Helper()
	ctx := newHandlerTestContext(context.Background(), filename)
	ctx.request.query = query
	service := &assetServiceFake{}
	response := NewAssetHandler(service).ServeAssetHead(ctx)
	if response.StatusCode() != http.StatusBadRequest {
		t.Fatalf("invalid HEAD request: got status=%d, want %d", response.StatusCode(), http.StatusBadRequest)
	}
	if service.calls != 0 {
		t.Fatalf("invalid request reached service %d times", service.calls)
	}
}

func TestServeAssetHeadPassesCriteria(t *testing.T) {
	criteria := handlerTestAvailableCriteria
	ctx := newHandlerTestContext(context.Background(), handlerTestFilename)
	ctx.request.query = map[string][]string{criteria: {""}}
	service := &assetServiceFake{document: &domain.AssetDocument{
		Filename: handlerTestFilename,
		Access:   handlerAllowedAccess,
	}}
	response := NewAssetHandler(service).ServeAssetHead(ctx)
	if response.StatusCode() != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.StatusCode())
	}
	assertHandlerServiceCall(t, service, handlerTestFilename)
	assertHandlerCriteria(t, service.call.criteria, criteria)
}
