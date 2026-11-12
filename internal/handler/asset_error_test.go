package handler

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/nawaphonOHM/resume-assets/internal/domain"
)

type handlerErrorCase struct {
	err    error
	name   string
	status int
}

func TestServeAssetMapsErrorsAndPropagatesContext(t *testing.T) {
	ctx := context.WithValue(context.Background(), handlerTestContextKey("request"), "request")
	tests := []handlerErrorCase{
		{status: http.StatusNotFound, err: domain.ErrAssetNotFound, name: "missing"},
		{status: http.StatusForbidden, err: domain.ErrAssetAccessDenied, name: "denied"},
		{status: http.StatusBadRequest, err: domain.ErrInvalidAssetCriteria, name: "invalid criteria"},
		{status: http.StatusConflict, err: domain.ErrAmbiguousAssetMatch, name: "ambiguous match"},
		{status: http.StatusInternalServerError, err: errors.New("private backend detail"), name: "backend failure"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			runHandlerErrorCase(ctx, t, test)
		})
	}
}

func runHandlerErrorCase(ctx context.Context, t *testing.T, test handlerErrorCase) {
	t.Helper()
	service := &assetServiceFake{err: test.err}
	response := NewAssetHandler(service).ServeAsset(newHandlerTestContext(ctx, handlerTestFilename))
	if response.StatusCode() != test.status {
		t.Fatalf("expected status %d, got %d", test.status, response.StatusCode())
	}
	assertHandlerServiceCall(t, service, handlerTestFilename)
	assertHandlerServiceContext(ctx, t, service)
}

func assertHandlerServiceContext(ctx context.Context, t *testing.T, service *assetServiceFake) {
	t.Helper()
	if service.call == nil {
		t.Fatal("asset service was not called")
	}
	if service.call.ctx != ctx {
		t.Fatal("request context was not propagated to service")
	}
}
