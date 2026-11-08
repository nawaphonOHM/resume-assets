package handler

import (
	"net/http"
	"strconv"
	"testing"
	"time"
)

func assertIntegrationAsset(t *testing.T, client *http.Client, baseURL string, service *routerAssetService) {
	t.Helper()
	response, body := performRequest(t, client, http.MethodGet, baseURL+"/assets/app.js")
	if response.StatusCode != http.StatusOK || body != integrationAssetData {
		t.Fatalf("successful asset response: status=%d body=%q", response.StatusCode, body)
	}
	assertIntegrationHeaders(t, response.Header)
	assertIntegrationBodyAndLookup(t, service)
}

func assertIntegrationHeaders(t *testing.T, headers http.Header) {
	t.Helper()
	expected := []struct{ name, value string }{
		{name: "Content-Type", value: integrationContentType},
		{name: "Content-Length", value: strconv.Itoa(len(integrationAssetData))},
		{name: "ETag", value: integrationETag},
		{name: "Last-Modified", value: "Fri, 02 Jan 2026 03:04:05 GMT"},
	}
	for _, header := range expected {
		if headers.Get(header.name) != header.value {
			t.Errorf("asset metadata header %s: got %q, want %q", header.name, headers.Get(header.name), header.value)
		}
	}
}

func assertIntegrationBodyAndLookup(t *testing.T, service *routerAssetService) {
	t.Helper()
	body, filenames := service.snapshot()
	if body == nil {
		t.Fatal("framework stream did not produce an object body")
	}
	assertIntegrationBodyClosed(t, body)
	assertIntegrationFirstLookup(t, filenames)
}

func assertIntegrationFirstLookup(t *testing.T, filenames []string) {
	t.Helper()
	if len(filenames) != 1 || filenames[0] != "app.js" {
		t.Fatalf("unexpected service lookup: %#v", filenames)
	}
}

func assertIntegrationBodyClosed(t *testing.T, body *integrationBody) {
	t.Helper()
	select {
	case <-body.closed:
	case <-time.After(integrationWaitTimeout):
		t.Fatal("framework stream did not close the object body")
	}
}
