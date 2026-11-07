package handler

import (
	"net/http"
	"testing"
)

const integrationDeniedFilename = "denied.js"

func TestFrameworkAssetAccessPolicy(t *testing.T) {
	repository := newAccessPolicyRepository()
	storageClient := &accessPolicyStorage{}
	client, baseURL := startAccessPolicyServer(t, repository, storageClient)
	assertPolicyAssetResponse(t, client, baseURL)
	assertPolicyDeniedResponse(t, client, baseURL)
	assertAccessPolicyStorageCalls(t, storageClient)
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
