package handler

import (
	"net/http"
	"testing"
)

func TestFrameworkAssetRouter(t *testing.T) {
	service := &routerAssetService{}
	baseURL := startFrameworkServer(t, RegisterRoutes(NewAssetHandler(service)))
	client := &http.Client{Transport: &http.Transport{DisableKeepAlives: true}}
	assertIntegrationAsset(t, client, baseURL, service)
	assertIntegrationRoutes(t, client, baseURL, service)
	assertIntegrationProbes(t, client, baseURL)
}

func TestFrameworkAssetQueryBinding(t *testing.T) {
	service := &routerAssetService{}
	baseURL := startFrameworkServer(t, RegisterRoutes(NewAssetHandler(service)))
	client := &http.Client{Transport: &http.Transport{DisableKeepAlives: true}}
	response, body := performRequest(
		t,
		client,
		http.MethodGet,
		baseURL+`/assets/app.js?%7B%22status%22%3A%22available%22%7D`,
	)
	assertQueryResponse(t, response, body)
	assertQueryCriteria(t, service.criteriaSnapshot())
}

func assertQueryResponse(t *testing.T, response *http.Response, body string) {
	t.Helper()
	if response.StatusCode != http.StatusOK || body != integrationAssetData {
		t.Fatalf("query criteria response: status=%d body=%q", response.StatusCode, body)
	}
}

func assertQueryCriteria(t *testing.T, criteria []*string) {
	t.Helper()
	if len(criteria) != 1 {
		t.Fatalf("unexpected criteria passed through the framework: %#v", criteria)
	}
	assertHandlerCriteria(t, criteria[0], handlerTestAvailableCriteria)
}
