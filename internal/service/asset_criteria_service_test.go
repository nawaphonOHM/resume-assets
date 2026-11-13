package service

import (
	"testing"

	"github.com/nawaphonOHM/resume-assets/internal/domain"
)

type criteriaLookupCase struct {
	criteria *string
	want     domain.AssetCriteria
	name     string
}

func TestGetAssetRejectsInvalidCriteriaBeforeLookupOrStorage(t *testing.T) {
	repository := &repositoryFake{}
	storageClient := &storageFake{}
	assetService := newServiceForTest(t, repository, storageClient)
	assertInvalidCriteriaRejected(t, assetService, repository, storageClient)
}

func TestGetAssetReturnsDocumentFromCriteriaLookup(t *testing.T) {
	criteriaCases := []criteriaLookupCase{
		{name: "absent"},
		{name: "empty object", criteria: stringPointer(`{}`), want: domain.AssetCriteria{}},
		{
			name: "criteria values", criteria: stringPointer(`{"status":"available"}`),
			want: domain.AssetCriteria{"status": "available"},
		},
	}
	for _, test := range criteriaCases {
		t.Run(test.name, func(t *testing.T) {
			runCriteriaLookupCase(t, test)
		})
	}
}
