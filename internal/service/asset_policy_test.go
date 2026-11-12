package service

import (
	"context"
	"testing"

	"github.com/nawaphonOHM/resume-assets/internal/domain"
)

type policyCase struct {
	repoErr  error
	wantErr  error
	document *domain.AssetDocument
	name     string
}

func TestFindAllowedAssetAccessCases(t *testing.T) {
	filename := "bundle/" + testAssetFilename
	tests := []policyCase{
		{
			name: "allowed",
			document: &domain.AssetDocument{
				Filename: filename,
				Location: testAssetFilename,
				Access:   testAccessAllowed,
			},
		},
		{
			name:     "denied",
			document: &domain.AssetDocument{Filename: filename, Access: "DENY"},
			wantErr:  domain.ErrAssetAccessDenied,
		},
		{
			name:    "missing",
			repoErr: domain.ErrAssetNotFound,
			wantErr: domain.ErrAssetNotFound,
		},
		{
			name:     "lowercase access fails closed",
			document: &domain.AssetDocument{Filename: filename, Access: "allow"},
			wantErr:  domain.ErrAssetAccessDenied,
		},
		{
			name:    "nil metadata",
			wantErr: domain.ErrInvalidAssetMetadata,
		},
		{
			name:     "mismatched metadata filename",
			document: &domain.AssetDocument{Filename: "other.js", Access: testAccessAllowed},
			wantErr:  domain.ErrInvalidAssetMetadata,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			runPolicyCase(t, test, filename)
		})
	}
}

func runPolicyCase(t *testing.T, test policyCase, filename string) {
	t.Helper()
	ctx := context.WithValue(context.Background(), testContextKey("request"), "request")
	repository := &repositoryFake{document: test.document, err: test.repoErr}
	storageClient := &storageFake{}
	service := newServiceForTest(t, repository, storageClient)

	asset, err := service.FindAllowedAsset(ctx, filename, nil)
	assertPolicyError(t, err, test.wantErr)
	assertPolicyAsset(t, asset, test.document, test.wantErr)
	assertRepositoryCall(ctx, t, repository.call, filename)
	assertStorageNotCalled(t, storageClient)
}
