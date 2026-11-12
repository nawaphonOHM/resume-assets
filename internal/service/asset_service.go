// Package service provides access policy for asset retrieval.
package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/nawaphonOHM/resume-assets/internal/domain"
	"github.com/nawaphonOHM/whatever/v3/pkg/storage"
)

var (
	ErrNilRepository          = errors.New("asset repository cannot be nil")
	ErrNilStorage             = errors.New("object storage client cannot be nil")
	ErrInvalidStorageResponse = errors.New("invalid storage response")
)

// AssetService applies the asset access policy before retrieval.
type AssetService struct {
	repository    domain.AssetRepository
	storage       domain.ObjectStorage
	defaultBucket string
}

// NewAssetService creates an asset service with its metadata and storage backends.
func NewAssetService(
	repository domain.AssetRepository,
	storageClient domain.ObjectStorage,
	defaultBucket string,
) (*AssetService, error) {
	if repository == nil {
		return nil, ErrNilRepository
	}
	if storageClient == nil {
		return nil, ErrNilStorage
	}
	return &AssetService{
		repository:    repository,
		storage:       storageClient,
		defaultBucket: defaultBucket,
	}, nil
}

// FindAllowedAsset looks up metadata and permits only an exact "ALLOW" access value.
func (s *AssetService) FindAllowedAsset(
	ctx context.Context,
	filename string,
	criteriaJSON *string,
) (*domain.AssetDocument, error) {
	criteria, err := parseOptionalAssetCriteria(criteriaJSON)
	if err != nil {
		return nil, err
	}
	return s.findAllowedWithCriteria(ctx, filename, criteria)
}

func (s *AssetService) findAllowedWithCriteria(
	ctx context.Context,
	filename string,
	criteria *domain.AssetCriteria,
) (*domain.AssetDocument, error) {
	document, err := s.repository.FindAsset(ctx, filename, criteria)
	if err != nil {
		return nil, fmt.Errorf("lookup asset metadata for filename %q: %w", filename, err)
	}
	if err := validateAssetDocument(document, filename); err != nil {
		return nil, err
	}
	return document, nil
}

// GetAsset retrieves an allowed object and its response metadata.
func (s *AssetService) GetAsset(
	ctx context.Context,
	filename string,
	criteriaJSON *string,
) (*domain.StaticAsset, error) {
	document, err := s.FindAllowedAsset(ctx, filename, criteriaJSON)
	if err != nil {
		return nil, err
	}
	return s.fetchAsset(ctx, document.Location)
}

func (s *AssetService) fetchAsset(ctx context.Context, location string) (*domain.StaticAsset, error) {
	bucket, key, err := ParseLocation(location, s.defaultBucket)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", domain.ErrInvalidAssetMetadata, err)
	}

	output, err := s.storage.GetObject(ctx, storage.GetObjectInput{Bucket: bucket, Key: key})
	if err != nil {
		return nil, mapStorageError(err)
	}
	return toStaticAsset(output)
}
