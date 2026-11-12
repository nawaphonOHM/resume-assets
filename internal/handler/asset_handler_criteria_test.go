package handler

import (
	"context"
	"net/http"
	"testing"

	"github.com/nawaphonOHM/resume-assets/internal/domain"
	"github.com/nawaphonOHM/whatever/v3/pkg/rest"
)

func TestServeAssetPassesScriptStyleQueryCriteria(t *testing.T) {
	criteria := handlerTestAvailableCriteria
	ctx := newHandlerTestContext(context.Background(), handlerTestFilename)
	ctx.request.query = map[string][]string{criteria: {""}}
	service := &assetServiceFake{asset: &domain.StaticAsset{Body: handlerBody{}}}
	response := NewAssetHandler(service).ServeAsset(ctx)
	assertAssetCriteriaRequest(t, response, service, criteria)
}

func TestServeAssetRejectsUnsupportedQueryShapes(t *testing.T) {
	tests := []struct {
		query map[string][]string
		name  string
	}{
		{name: "named parameter", query: map[string][]string{"criteria": {handlerTestAvailableCriteria}}},
		{name: "additional parameter", query: map[string][]string{handlerTestAvailableCriteria: {""}, "other": {""}}},
		{name: "repeated criteria", query: map[string][]string{handlerTestAvailableCriteria: {"", ""}}},
		{name: "parameter value", query: map[string][]string{handlerTestAvailableCriteria: {"unexpected"}}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) { assertUnsupportedAssetQuery(t, test.query) })
	}
}

func assertAssetCriteriaRequest(t *testing.T, response rest.Response, service *assetServiceFake, criteria string) {
	t.Helper()
	if response.StatusCode() != http.StatusOK {
		t.Fatalf("expected status 200, got %d", response.StatusCode())
	}
	assertHandlerServiceCall(t, service, handlerTestFilename)
	assertHandlerCriteria(t, service.call.criteria, criteria)
}

func assertHandlerCriteria(t *testing.T, got *string, want string) {
	t.Helper()
	if got == nil || *got != want {
		t.Fatalf("expected criteria %q, got %#v", want, got)
	}
}

func assertUnsupportedAssetQuery(t *testing.T, query map[string][]string) {
	t.Helper()
	ctx := newHandlerTestContext(context.Background(), handlerTestFilename)
	ctx.request.query = query
	service := &assetServiceFake{}
	response := NewAssetHandler(service).ServeAsset(ctx)
	if response.StatusCode() != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", response.StatusCode())
	}
	if service.calls != 0 {
		t.Fatalf("unsupported query reached service %d times", service.calls)
	}
}
