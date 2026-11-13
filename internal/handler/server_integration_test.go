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
	assertIntegrationHeadResponses(t, client, baseURL)
	assertIntegrationRoutes(t, client, baseURL, service)
	assertIntegrationProbes(t, client, baseURL)
}

func TestFrameworkAssetQueryBinding(t *testing.T) {
	service := &routerAssetService{}
	baseURL := startFrameworkServer(t, RegisterRoutes(NewAssetHandler(service)))
	client := &http.Client{Transport: &http.Transport{DisableKeepAlives: true}}
	url :=
		baseURL +
			`/assets/app.js?%7B%22status%22%3A%22available%22%7D`
	response, body := performRequest(
		t,
		client,
		http.MethodGet,
		url,
	)
	assertQueryResponse(t, response, body)
	response, body = performRequest(t, client, http.MethodHead, url)
	if response.StatusCode != http.StatusOK || body != "" {
		t.Fatalf("HEAD query criteria response: status=%d body=%q", response.StatusCode, body)
	}
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
	if len(criteria) != 2 {
		t.Fatalf("unexpected criteria passed through the framework: %#v", criteria)
	}
	for _, value := range criteria {
		assertHandlerCriteria(t, value, handlerTestAvailableCriteria)
	}
}
