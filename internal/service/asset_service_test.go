package service

import (
	"context"
	"errors"
	"io"
	"testing"

	"github.com/nawaphonOHM/resume-assets/internal/domain"
	"github.com/nawaphonOHM/whatever/v3/pkg/storage"
)

const (
	testAssetFilename = "app.js"
	testAccessAllowed = "ALLOW"
	testDefaultBucket = "default-bucket"
)

type testContextKey string

type repositoryCall struct {
	criteria *domain.AssetCriteria
	ctx      context.Context
	filename string
}

type repositoryFake struct {
	document *domain.AssetDocument
	err      error
	call     *repositoryCall
}

func (r *repositoryFake) FindAsset(
	ctx context.Context,
	filename string,
	criteria *domain.AssetCriteria,
) (*domain.AssetDocument, error) {
	r.call = &repositoryCall{ctx: ctx, filename: filename, criteria: criteria}
	return r.document, r.err
}

type storageCall struct {
	ctx   context.Context
	input storage.GetObjectInput
}

type storageFake struct {
	err    error
	call   *storageCall
	output *storage.GetObjectOutput
	calls  int
}

func (s *storageFake) GetObject(ctx context.Context, input storage.GetObjectInput) (*storage.GetObjectOutput, error) {
	s.calls++
	s.call = &storageCall{ctx: ctx, input: input}
	return s.output, s.err
}

type closeTracker struct {
	closed bool
}

func (*closeTracker) Read([]byte) (int, error) {
	return 0, io.EOF
}

func (b *closeTracker) Close() error {
	b.closed = true
	return nil
}

func newServiceForTest(
	t *testing.T,
	repository domain.AssetRepository,
	storageClient domain.ObjectStorage,
) *AssetService {
	t.Helper()
	service, err := NewAssetService(repository, storageClient, testDefaultBucket)
	if err != nil {
		t.Fatalf("NewAssetService: %v", err)
	}
	return service
}

func TestNewAssetServiceRejectsNilRepository(t *testing.T) {
	if _, err := NewAssetService(nil, &storageFake{}, ""); !errors.Is(err, ErrNilRepository) {
		t.Fatalf("expected ErrNilRepository, got %v", err)
	}
}

func TestNewAssetServiceRejectsNilStorage(t *testing.T) {
	if _, err := NewAssetService(&repositoryFake{}, nil, ""); !errors.Is(err, ErrNilStorage) {
		t.Fatalf("expected ErrNilStorage, got %v", err)
	}
}
