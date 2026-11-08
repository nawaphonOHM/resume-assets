package handler

import (
	"io"
	"net/http"
	"testing"
)

func performRequest(t *testing.T, client *http.Client, method, url string) (*http.Response, string) {
	t.Helper()
	request := newIntegrationRequest(t, method, url)
	response := doIntegrationRequest(t, client, request)
	defer closeIntegrationResponse(t, response)
	return response, readIntegrationResponse(t, response)
}

func newIntegrationRequest(t *testing.T, method, url string) *http.Request {
	t.Helper()
	request, err := http.NewRequest(method, url, http.NoBody)
	if err != nil {
		t.Fatalf("create request: %v", err)
	}
	return request
}

func doIntegrationRequest(t *testing.T, client *http.Client, request *http.Request) *http.Response {
	t.Helper()
	response, err := client.Do(request)
	if err != nil {
		t.Fatalf("request %s %s: %v", request.Method, request.URL, err)
	}
	return response
}

func readIntegrationResponse(t *testing.T, response *http.Response) string {
	t.Helper()
	data, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("read response body: %v", err)
	}
	return string(data)
}

func closeIntegrationResponse(t *testing.T, response *http.Response) {
	t.Helper()
	if err := response.Body.Close(); err != nil {
		t.Errorf("close response body: %v", err)
	}
}
