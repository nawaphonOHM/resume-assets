package handler

import (
	"net/http"
	"testing"
)

const integrationExpectedLookups = 9

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
		{method: http.MethodHead, path: "/assets/missing", status: http.StatusNotFound},
		{method: http.MethodHead, path: "/assets/denied", status: http.StatusForbidden},
		{method: http.MethodHead, path: "/assets/broken", status: http.StatusInternalServerError},
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
	assertIntegrationRouteStatus(t, response, body, test)
	assertIntegrationHeadRoute(t, response, body, test)
	assertIntegrationPrivateDetail(t, test, body)
}

func assertIntegrationRouteStatus(t *testing.T, response *http.Response, body string, test integrationRouteCase) {
	t.Helper()
	if response.StatusCode != test.status {
		t.Errorf(
			"%s %s: got status %d, want %d; body=%q",
			test.method, test.path, response.StatusCode, test.status, body,
		)
	}
}

func assertIntegrationHeadRoute(t *testing.T, response *http.Response, body string, test integrationRouteCase) {
	t.Helper()
	if test.method != http.MethodHead {
		return
	}
	if body != "" {
		t.Errorf("HEAD %s returned a body: %q", test.path, body)
	}
	assertHeadMetadataOmittedOnError(t, response, test)
}

func assertHeadMetadataOmittedOnError(t *testing.T, response *http.Response, test integrationRouteCase) {
	t.Helper()
	if test.status == http.StatusOK {
		return
	}
	for _, name := range []string{"ETag", "Last-Modified", "Content-Type", "Content-Length", "X-Head-Mode"} {
		assertHeadErrorHeaderOmitted(t, response, test, name)
	}
}

func assertHeadErrorHeaderOmitted(t *testing.T, response *http.Response, test integrationRouteCase, name string) {
	t.Helper()
	if response.Header.Get(name) != "" {
		t.Errorf(
			"HEAD error %s included asset metadata header %s=%q",
			test.path, name, response.Header.Get(name),
		)
	}
}
