package handler

import (
	"net/http"
	"strconv"
	"testing"
)

func assertIntegrationAsset(t *testing.T, client *http.Client, baseURL string, service *routerAssetService) {
	t.Helper()
	response, body := performRequest(t, client, http.MethodGet, baseURL+"/assets/app.js")
	if response.StatusCode != http.StatusOK || body != integrationAssetData {
		t.Fatalf("successful asset response: status=%d body=%q", response.StatusCode, body)
	}
	assertIntegrationGetHeaders(t, response.Header, []integrationHeader{
		{name: integrationHeaderContentType, value: integrationContentType},
		{name: handlerContentLengthHeader, value: strconv.Itoa(len(integrationAssetData))},
		{name: handlerETagHeader, value: integrationMetadataETag},
		{name: integrationHeaderLastModified, value: "Fri, 02 Jan 2026 03:04:05 GMT"},
		{name: integrationHeaderMode, value: "cache"},
	})
	assertIntegrationBodyAndLookup(t, service, "app.js")
}

func assertIntegrationPDFAsset(t *testing.T, client *http.Client, baseURL string, service *routerAssetService) {
	t.Helper()
	response, body := performRequest(t, client, http.MethodGet, baseURL+"/assets/app.PDF")
	if response.StatusCode != http.StatusOK || body != integrationAssetData {
		t.Fatalf("successful PDF asset response: status=%d body=%q", response.StatusCode, body)
	}
	assertIntegrationGetHeaders(t, response.Header, []integrationHeader{
		{name: integrationHeaderContentType, value: integrationPDFContentType},
		{name: handlerContentLengthHeader, value: strconv.Itoa(len(integrationAssetData))},
		{name: handlerETagHeader, value: integrationETag},
		{name: integrationHeaderLastModified, value: "Fri, 02 Jan 2026 03:04:05 GMT"},
		{name: integrationHeaderMode, value: "download"},
	})
	assertIntegrationBodyAndLookup(t, service, "app.js", "app.PDF")
}

func assertIntegrationHeadResponses(t *testing.T, client *http.Client, baseURL string) {
	t.Helper()
	assertCachingHeadResponse(t, client, baseURL)
	assertDownloadingHeadResponse(t, client, baseURL)
}

func assertCachingHeadResponse(t *testing.T, client *http.Client, baseURL string) {
	t.Helper()
	response, body := performRequest(t, client, http.MethodHead, baseURL+"/assets/app.js")
	assertSuccessfulHeadResponse(t, response, body, "non-PDF")
	assertHeadResponseHeader(t, response.Header, handlerETagHeader, integrationMetadataETag)
	assertHeadResponseHeader(t, response.Header, integrationHeaderMode, "cache")
	assertHeadResponseHeader(t, response.Header, integrationHeaderContentType, "")
	assertHeadResponseHeader(t, response.Header, handlerContentLengthHeader, "")
}

func assertDownloadingHeadResponse(t *testing.T, client *http.Client, baseURL string) {
	t.Helper()
	response, body := performRequest(t, client, http.MethodHead, baseURL+"/assets/app.PDF")
	assertSuccessfulHeadResponse(t, response, body, "PDF")
	assertHeadResponseHeader(t, response.Header, handlerETagHeader, "")
	assertHeadResponseHeader(t, response.Header, integrationHeaderMode, "download")
	assertHeadResponseHeader(t, response.Header, integrationHeaderContentType, integrationPDFContentType)
	assertHeadResponseHeader(t, response.Header, handlerContentLengthHeader, strconv.Itoa(len(integrationAssetData)))
}

func assertSuccessfulHeadResponse(t *testing.T, response *http.Response, body, assetType string) {
	t.Helper()
	if response.StatusCode != http.StatusOK || body != "" {
		t.Fatalf("%s HEAD response: status=%d body=%q", assetType, response.StatusCode, body)
	}
}

func assertHeadResponseHeader(t *testing.T, headers http.Header, name, want string) {
	t.Helper()
	if got := headers.Get(name); got != want {
		t.Errorf("HEAD header %s: got %q, want %q", name, got, want)
	}
}

type integrationHeader struct {
	name  string
	value string
}

func assertIntegrationGetHeaders(t *testing.T, headers http.Header, expected []integrationHeader) {
	t.Helper()
	for _, header := range expected {
		if headers.Get(header.name) != header.value {
			t.Errorf("asset metadata header %s: got %q, want %q", header.name, headers.Get(header.name), header.value)
		}
	}
}
