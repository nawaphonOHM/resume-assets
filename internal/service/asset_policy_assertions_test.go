package service

import (
	"context"
	"errors"
	"testing"

	"github.com/nawaphonOHM/resume-assets/internal/domain"
)

func assertPolicyError(t *testing.T, got, want error) {
	t.Helper()
	if !errors.Is(got, want) {
		t.Fatalf("expected error %v, got %v", want, got)
	}
}

func assertPolicyAsset(t *testing.T, got, want *domain.AssetDocument, wantErr error) {
	t.Helper()
	if wantErr == nil && got != want {
		t.Fatalf("expected original metadata %#v, got %#v", want, got)
	}
}

func assertRepositoryCall(ctx context.Context, t *testing.T, got *repositoryCall, filename string) {
	t.Helper()
	if got == nil {
		t.Fatal("repository call was not recorded")
	}
	assertRepositoryFilename(t, got, filename)
	assertRepositoryContext(ctx, t, got)
}

func assertRepositoryFilename(t *testing.T, got *repositoryCall, filename string) {
	t.Helper()
	if got.filename != filename {
		t.Fatalf("expected exact filename %q, got %q", filename, got.filename)
	}
}

func assertRepositoryContext(ctx context.Context, t *testing.T, got *repositoryCall) {
	t.Helper()
	if got.ctx != ctx {
		t.Fatal("request context was not propagated to repository")
	}
}

func assertStorageNotCalled(t *testing.T, client *storageFake) {
	t.Helper()
	if client.calls != 0 {
		t.Fatalf("expected no storage calls during authorization, got %d", client.calls)
	}
}
