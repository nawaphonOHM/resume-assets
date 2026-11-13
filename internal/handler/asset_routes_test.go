package handler

import (
	"testing"

	"github.com/nawaphonOHM/whatever/v3/pkg/rest"
)

func TestNewAPIRegistration(t *testing.T) {
	registration := NewAPIRegistration(NewAssetHandler(&assetServiceFake{}))
	if registration.Group != rest.APIGroupRaw {
		t.Fatalf("expected raw route group, got %q", registration.Group)
	}
	if len(registration.Apis) != 2 {
		t.Fatalf("expected GET and HEAD routes, got %d", len(registration.Apis))
	}
	assertAssetRoute(t, registration.Apis[0], rest.GET)
	assertAssetRoute(t, registration.Apis[1], rest.HEAD)
}

func assertAssetRoute(t *testing.T, route *rest.ExportableAPI, method any) {
	t.Helper()
	if route.Method != method {
		t.Fatalf("expected %v route, got %v", method, route.Method)
	}
	if route.Path != RouteAssetPath {
		t.Fatalf("expected route %q, got %q", RouteAssetPath, route.Path)
	}
}
