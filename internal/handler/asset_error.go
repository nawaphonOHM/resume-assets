package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/nawaphonOHM/resume-assets/internal/domain"
	"github.com/nawaphonOHM/whatever/v3/pkg/logging"
	"github.com/nawaphonOHM/whatever/v3/pkg/rest"
)

func mapAssetError(ctx context.Context, err error) rest.Response {
	if errors.Is(err, domain.ErrInvalidAssetCriteria) {
		return rest.BadRequest(codeInvalidCriteria, "asset criteria must be one valid JSON object")
	}
	if errors.Is(err, domain.ErrAmbiguousAssetMatch) {
		return rest.Error(http.StatusConflict, codeAmbiguousAsset, "the filename and criteria match multiple assets")
	}
	return mapAssetLookupError(ctx, err)
}

func mapAssetLookupError(ctx context.Context, err error) rest.Response {
	if errors.Is(err, domain.ErrAssetNotFound) {
		return rest.NotFound(codeAssetNotFound, "the requested asset was not found")
	}
	if errors.Is(err, domain.ErrAssetAccessDenied) {
		return rest.Forbidden(codeAccessDenied, "access to the requested asset is denied")
	}
	logging.ErrorContext(ctx, "failed to serve asset", "error", err)
	return rest.InternalServerError(codeInternalError, "the server failed to serve the requested asset")
}
