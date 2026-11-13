package service

import (
	"context"
	"io"
	"testing"

	"github.com/nawaphonOHM/resume-assets/internal/domain"
)

func assertStorageInput(t *testing.T, client *storageFake, bucket, key string) {
	t.Helper()
	assertStorageCallCount(t, client, 1)
	if client.call == nil {
		t.Fatal("storage call was not recorded")
	}
	assertStorageTarget(t, client.call, bucket, key)
}

func assertStorageTarget(t *testing.T, call *storageCall, bucket, key string) {
	t.Helper()
	if call.input.Bucket != bucket {
		t.Fatalf("unexpected storage bucket: %q", call.input.Bucket)
	}
	if call.input.Key != key {
		t.Fatalf("unexpected storage key: %q", call.input.Key)
	}
}

func assertStorageContext(ctx context.Context, t *testing.T, call *storageCall) {
	t.Helper()
	if call == nil || call.ctx != ctx {
		t.Fatal("request context was not propagated to storage")
	}
}

func assertStaticAsset(t *testing.T, asset *domain.StaticAsset, expected domain.StaticAsset) {
	t.Helper()
	if asset == nil {
		t.Fatal("expected an asset")
	}
	if *asset != expected {
		t.Fatalf("unexpected asset output: got %#v, want %#v", asset, expected)
	}
}

func assertNoAssetDocument(t *testing.T, document *domain.AssetDocument) {
	t.Helper()
	if document != nil {
		t.Fatalf("expected nil asset document, got %#v", document)
	}
}

func closeAssetBody(t *testing.T, body io.ReadCloser) {
	t.Helper()
	if err := body.Close(); err != nil {
		t.Fatalf("close returned body: %v", err)
	}
}

func assertBodyOpen(t *testing.T, body *closeTracker) {
	t.Helper()
	if body.closed {
		t.Fatal("service closed the object body before returning it")
	}
}

func assertBodyClosed(t *testing.T, body *closeTracker) {
	t.Helper()
	if !body.closed {
		t.Fatal("returned object body was not closed by its caller")
	}
}
