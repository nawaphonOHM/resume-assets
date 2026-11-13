package service

import (
	"testing"

	"github.com/nawaphonOHM/resume-assets/internal/domain"
)

func assertGetAssetSucceeded(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("GetAsset: %v", err)
	}
}

func assertSelectedDocument(t *testing.T, got, want *domain.AssetDocument) {
	t.Helper()
	if got != want {
		t.Fatalf("expected selected document %p, got %p", want, got)
	}
}

func assertCriteriaRepositoryCall(t *testing.T, repository *repositoryFake, want domain.AssetCriteria) {
	t.Helper()
	assertSingleCriteriaLookup(t, repository)
	if want == nil {
		assertAbsentCriteria(t, repository.call.criteria)
		return
	}
	assertCriteriaValue(t, repository.call.criteria, want)
}

func assertSingleCriteriaLookup(t *testing.T, repository *repositoryFake) {
	t.Helper()
	if repository.calls != 1 || repository.call == nil {
		t.Fatalf("expected one metadata lookup, got %d calls", repository.calls)
	}
}

func assertAbsentCriteria(t *testing.T, got *domain.AssetCriteria) {
	t.Helper()
	if got != nil {
		t.Fatalf("expected absent criteria, got %#v", got)
	}
}

func assertReturnedAssetBody(t *testing.T, asset *domain.StaticAsset, body *closeTracker) {
	t.Helper()
	if asset == nil || asset.Body == nil {
		t.Fatal("expected retrieved asset stream")
	}
	closeAssetBody(t, asset.Body)
	assertBodyClosed(t, body)
}
