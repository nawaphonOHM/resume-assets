package handler

import (
	"net/http"
	"strings"
	"testing"
)

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
