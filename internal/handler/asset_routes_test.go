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
	if len(registration.Apis) != 1 {
		t.Fatalf("expected one application route, got %d", len(registration.Apis))
	}
	assertAssetRoute(t, registration.Apis[0])
}

func assertAssetRoute(t *testing.T, route *rest.ExportableAPI) {
	t.Helper()
	if route.Method != rest.GET {
		t.Fatalf("expected GET route, got %v", route.Method)
	}
	if route.Path != RouteAssetPath {
		t.Fatalf("expected route %q, got %q", RouteAssetPath, route.Path)
	}
}
