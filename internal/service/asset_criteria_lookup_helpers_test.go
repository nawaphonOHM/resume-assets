package service

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/nawaphonOHM/resume-assets/internal/domain"
)

type parsedCriteriaTestFixture struct {
	document     *domain.AssetDocument
	repository   *repositoryFake
	storage      *storageFake
	service      *AssetService
	criteriaJSON string
}

func newParsedCriteriaFixture(t *testing.T) parsedCriteriaTestFixture {
	t.Helper()
	document := &domain.AssetDocument{Filename: testAssetFilename, Location: "key", Access: testAccessAllowed}
	repository := &repositoryFake{document: document}
	storage := &storageFake{}
	return parsedCriteriaTestFixture{
		document:     document,
		repository:   repository,
		storage:      storage,
		service:      newServiceForTest(t, repository, storage),
		criteriaJSON: `{"status":"available","nested":{"number":9007199254740993}}`,
	}
}

func assertFoundAsset(t *testing.T, got, want *domain.AssetDocument) {
	t.Helper()
	if got != want {
		t.Fatalf("expected metadata %#v, got %#v", want, got)
	}
}

func assertParsedCriteria(t *testing.T, repository *repositoryFake) {
	t.Helper()
	want := domain.AssetCriteria{
		"status": "available",
		"nested": map[string]any{"number": json.Number("9007199254740993")},
	}
	if repository.call == nil {
		t.Fatal("repository call was not recorded")
	}
	assertCriteriaValue(t, repository.call.criteria, want)
}

func assertCriteriaValue(t *testing.T, got *domain.AssetCriteria, want domain.AssetCriteria) {
	t.Helper()
	if got == nil || !reflect.DeepEqual(*got, want) {
		t.Fatalf("repository criteria: got %#v, want %#v", got, want)
	}
}
