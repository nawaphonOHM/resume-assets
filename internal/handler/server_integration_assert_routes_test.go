package handler

import (
	"net/http"
	"strings"
	"testing"
)

const integrationExpectedLookups = 4

type integrationRouteCase struct {
	method string
	path   string
	status int
}

func assertIntegrationRoutes(t *testing.T, client *http.Client, baseURL string, service *routerAssetService) {
	t.Helper()
	cases := []integrationRouteCase{
		{method: http.MethodGet, path: "/assets/missing", status: http.StatusNotFound},
		{method: http.MethodGet, path: "/assets/denied", status: http.StatusForbidden},
		{method: http.MethodGet, path: "/assets/broken", status: http.StatusInternalServerError},
		{method: http.MethodPost, path: "/assets/app.js", status: http.StatusMethodNotAllowed},
		{method: http.MethodGet, path: "/assets/nested/app.js", status: http.StatusNotFound},
	}
	for _, test := range cases {
		assertIntegrationRoute(t, client, baseURL, test)
	}
	_, filenames := service.snapshot()
	if len(filenames) != integrationExpectedLookups {
		t.Fatalf("unsupported/nested routes reached the asset service: %#v", filenames)
	}
}

func assertIntegrationRoute(t *testing.T, client *http.Client, baseURL string, test integrationRouteCase) {
	t.Helper()
	response, body := performRequest(t, client, test.method, baseURL+test.path)
	if response.StatusCode != test.status {
		t.Errorf(
			"%s %s: got status %d, want %d; body=%q",
			test.method, test.path, response.StatusCode, test.status, body,
		)
	}
	assertIntegrationPrivateDetail(t, test, body)
}

func assertIntegrationPrivateDetail(t *testing.T, test integrationRouteCase, body string) {
	t.Helper()
	if test.path == "/assets/broken" && strings.Contains(body, integrationPrivateDetail) {
		t.Error("internal backend error was exposed in the response")
	}
}

func assertIntegrationProbes(t *testing.T, client *http.Client, baseURL string) {
	t.Helper()
	for _, probe := range []string{"/health", "/ready"} {
		response, _ := performRequest(t, client, http.MethodGet, baseURL+probe)
		if response.StatusCode != http.StatusOK {
			t.Errorf("probe %s returned %d", probe, response.StatusCode)
		}
	}
}
