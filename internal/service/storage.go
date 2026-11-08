package service

import (
	"errors"
	"fmt"

	"github.com/nawaphonOHM/resume-assets/internal/domain"
	"github.com/nawaphonOHM/whatever/v3/pkg/storage"
)

func mapStorageError(err error) error {
	if errors.Is(err, storage.ErrObjectNotFound) {
		return fmt.Errorf("%w: %w", domain.ErrAssetNotFound, err)
	}
	return fmt.Errorf("fetch asset object: %w", err)
}

func toStaticAsset(output *storage.GetObjectOutput) (*domain.StaticAsset, error) {
	if output == nil || output.Body == nil {
		return nil, ErrInvalidStorageResponse
	}
	return &domain.StaticAsset{
		LastModified:  output.LastModified,
		Body:          output.Body,
		ContentType:   output.ContentType,
		ETag:          output.ETag,
		ContentLength: output.ContentLength,
	}, nil
}
