package handler

import (
	"net/http"
	"testing"
)

const integrationDeniedFilename = "denied.js"

type policyHeadCase struct {
	url    string
	access string
	status int
}

func TestFrameworkAssetAccessPolicy(t *testing.T) {
	repository := newAccessPolicyRepository()
	storageClient := &accessPolicyStorage{}
	client, baseURL := startAccessPolicyServer(t, repository, storageClient)
	assertPolicyAssetResponse(t, client, baseURL)
	assertPolicyHeadResponses(t, client, baseURL)
	assertPolicyDeniedResponse(t, client, baseURL)
	assertPolicyMissingTagsResponse(t, client, baseURL)
	assertAccessPolicyStorageCalls(t, storageClient)
}

func assertPolicyHeadResponses(t *testing.T, client *http.Client, baseURL string) {
	t.Helper()
	assertPolicyHeadStatus(t, client, policyHeadCase{
		status: http.StatusOK, url: baseURL + "/assets/app.js", access: "ALLOW",
	})
	response := assertPolicyHeadStatus(t, client, policyHeadCase{
		status: http.StatusForbidden,
		url:    baseURL + "/assets/" + integrationDeniedFilename,
		access: "DENY",
	})
	if response.Header.Get("ETag") != "" {
		t.Fatalf("DENY HEAD response exposed asset metadata: ETag=%q", response.Header.Get("ETag"))
	}
}

func assertPolicyHeadStatus(t *testing.T, client *http.Client, test policyHeadCase) *http.Response {
	t.Helper()
	response, body := performRequest(t, client, http.MethodHead, test.url)
	if response.StatusCode != test.status || body != "" {
		t.Fatalf("%s HEAD response: status=%d body=%q", test.access, response.StatusCode, body)
	}
	return response
}

func assertAccessPolicyStorageCalls(t *testing.T, storageClient *accessPolicyStorage) {
	t.Helper()
	if storageClient.calls != 1 {
		t.Fatalf("expected storage to be fetched once for ALLOW only, got %d calls", storageClient.calls)
	}
	assertAccessPolicyStorageRequest(t, storageClient.input.Bucket, storageClient.input.Key)
}

func assertAccessPolicyStorageRequest(t *testing.T, bucket, key string) {
	t.Helper()
	if bucket != "assets" || key != "app.js" {
		t.Fatalf("unexpected storage request: bucket=%q key=%q", bucket, key)
	}
}

func assertPolicyAssetResponse(t *testing.T, client *http.Client, baseURL string) {
	t.Helper()
	response, body := performRequest(t, client, http.MethodGet, baseURL+"/assets/app.js")
	if response.StatusCode != http.StatusOK || body != integrationAssetData {
		t.Fatalf("ALLOW response: status=%d body=%q", response.StatusCode, body)
	}
}

func assertPolicyDeniedResponse(t *testing.T, client *http.Client, baseURL string) {
	t.Helper()
	response, _ := performRequest(t, client, http.MethodGet, baseURL+"/assets/"+integrationDeniedFilename)
	if response.StatusCode != http.StatusForbidden {
		t.Fatalf("DENY response: got status %d, want %d", response.StatusCode, http.StatusForbidden)
	}
}

func assertPolicyMissingTagsResponse(t *testing.T, client *http.Client, baseURL string) {
	t.Helper()
	response, _ := performRequest(
		t, client, http.MethodGet, baseURL+`/assets/app.js?%7B%22status%22%3A%22available%22%7D`,
	)
	if response.StatusCode != http.StatusNotFound {
		t.Fatalf("missing tags response: got status %d, want %d", response.StatusCode, http.StatusNotFound)
	}
}
