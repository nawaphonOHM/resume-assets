package handler

import (
	"errors"

	"github.com/nawaphonOHM/resume-assets/internal/domain"
)

func integrationAssetError(filename string) error {
	if filename == "missing" {
		return domain.ErrAssetNotFound
	}
	if filename == "denied" {
		return domain.ErrAssetAccessDenied
	}
	return integrationBackendError(filename)
}

func integrationBackendError(filename string) error {
	if filename == "broken" {
		return errors.New(integrationPrivateDetail)
	}
	return nil
}
