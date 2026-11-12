package service

import "github.com/nawaphonOHM/resume-assets/internal/domain"

func validateAssetDocument(document *domain.AssetDocument, filename string) error {
	if err := validateAssetAccess(document); err != nil {
		return err
	}
	if document.Filename != filename {
		return domain.ErrInvalidAssetMetadata
	}
	return nil
}

func validateAssetAccess(document *domain.AssetDocument) error {
	if document == nil {
		return domain.ErrInvalidAssetMetadata
	}
	if document.Access != "ALLOW" {
		return domain.ErrAssetAccessDenied
	}
	return nil
}
