package service

import (
	"testing"

	"github.com/nawaphonOHM/resume-assets/internal/domain"
)

func TestFindAllowedAssetParsesAndPassesCriteria(t *testing.T) {
	fixture := newParsedCriteriaFixture(t)
	ctx := t.Context()
	got, err := fixture.service.FindAllowedAsset(ctx, testAssetFilename, &fixture.criteriaJSON)
	if err != nil {
		t.Fatalf("FindAllowedAsset: %v", err)
	}
	assertFoundAsset(t, got, fixture.document)
	assertParsedCriteria(t, fixture.repository)
	assertRepositoryCall(ctx, t, fixture.repository.call, testAssetFilename)
	assertStorageNotCalled(t, fixture.storage)
}

func TestFindAllowedAssetPreservesAbsentCriteria(t *testing.T) {
	criteria := lookupAssetCriteria(t, nil)
	if criteria != nil {
		t.Fatalf("expected absent criteria, got %#v", criteria)
	}
}

func TestFindAllowedAssetPreservesEmptyCriteria(t *testing.T) {
	criteria := lookupAssetCriteria(t, stringPointer(`{}`))
	if criteria == nil {
		t.Fatal("expected criteria to be present")
	}
	assertEmptyCriteria(t, *criteria)
}

func lookupAssetCriteria(t *testing.T, criteriaJSON *string) *domain.AssetCriteria {
	t.Helper()
	repository := &repositoryFake{document: &domain.AssetDocument{
		Filename: testAssetFilename,
		Access:   testAccessAllowed,
	}}
	assetService := newServiceForTest(t, repository, &storageFake{})
	if _, err := assetService.FindAllowedAsset(t.Context(), testAssetFilename, criteriaJSON); err != nil {
		t.Fatalf("FindAllowedAsset: %v", err)
	}
	if repository.call == nil {
		t.Fatal("repository call was not recorded")
	}
	return repository.call.criteria
}
