package handler

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/nawaphonOHM/resume-assets/internal/domain"
	"github.com/nawaphonOHM/resume-assets/internal/service"
	"github.com/nawaphonOHM/whatever/v3/pkg/storage"
)

type accessPolicyRepository struct {
	records map[string]*domain.AssetDocument
}

func (r *accessPolicyRepository) FindAsset(
	ctx context.Context,
	filename string,
) (*domain.AssetDocument, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return r.records[filename], nil
}

type accessPolicyStorage struct {
	input storage.GetObjectInput
	calls int
}

func (s *accessPolicyStorage) GetObject(
	ctx context.Context,
	input storage.GetObjectInput,
) (*storage.GetObjectOutput, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	s.calls++
	s.input = input
	body := &integrationBody{reader: strings.NewReader(integrationAssetData), closed: make(chan struct{})}
	return &storage.GetObjectOutput{
		Body:          body,
		ContentType:   integrationContentType,
		ContentLength: int64(len(integrationAssetData)),
	}, nil
}

func newAccessPolicyRepository() *accessPolicyRepository {
	return &accessPolicyRepository{records: map[string]*domain.AssetDocument{
		"app.js": {Filename: "app.js", Location: "s3://assets/app.js", Access: "ALLOW"},
		integrationDeniedFilename: {
			Filename: integrationDeniedFilename,
			Location: "s3://assets/denied.js",
			Access:   "DENY",
		},
	}}
}

func startAccessPolicyServer(
	t *testing.T,
	repository domain.AssetRepository,
	storageClient domain.ObjectStorage,
) (*http.Client, string) {
	t.Helper()
	assetService, err := service.NewAssetService(repository, storageClient, "default-bucket")
	if err != nil {
		t.Fatalf("NewAssetService: %v", err)
	}
	baseURL := startFrameworkServer(t, RegisterRoutes(NewAssetHandler(assetService)))
	client := &http.Client{Transport: &http.Transport{DisableKeepAlives: true}}
	return client, baseURL
}
