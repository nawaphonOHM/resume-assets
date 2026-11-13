package service

import (
	"testing"
	"time"

	"github.com/nawaphonOHM/resume-assets/internal/domain"
)

func assertInvalidCriteriaRejected(
	t *testing.T,
	assetService *AssetService,
	repository *repositoryFake,
	storageClient *storageFake,
) {
	t.Helper()
	criteriaJSON := `[]`
	asset, document, err := assetService.GetAsset(t.Context(), testAssetFilename, &criteriaJSON)
	assertNoStaticAsset(t, asset)
	assertNoAssetDocument(t, document)
	assertAssetError(t, err, domain.ErrInvalidAssetCriteria)
	if repository.call != nil {
		t.Fatalf("invalid criteria reached repository: %#v", repository.call)
	}
	assertStorageNotCalled(t, storageClient)
}

func runCriteriaLookupCase(t *testing.T, test criteriaLookupCase) {
	t.Helper()
	fixture := newCriteriaLookupFixture(t)
	asset, gotDocument, err := fixture.service.GetAsset(t.Context(), testAssetFilename, test.criteria)
	assertGetAssetSucceeded(t, err)
	assertSelectedDocument(t, gotDocument, fixture.document)
	assertCriteriaRepositoryCall(t, fixture.repository, test.want)
	assertReturnedAssetBody(t, asset, fixture.body)
}

type criteriaLookupFixture struct {
	document   *domain.AssetDocument
	repository *repositoryFake
	body       *closeTracker
	service    *AssetService
}

func newCriteriaLookupFixture(t *testing.T) criteriaLookupFixture {
	document := &domain.AssetDocument{
		Filename: testAssetFilename,
		Location: "key",
		Access:   testAccessAllowed,
	}
	repository := &repositoryFake{document: document}
	body := &closeTracker{}
	storageClient := newStorageFixture(body, time.Time{})
	return criteriaLookupFixture{
		document: document, repository: repository, body: body,
		service: newServiceForTest(t, repository, storageClient),
	}
}
