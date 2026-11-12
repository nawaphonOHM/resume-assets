// Package handler provides HTTP handlers for serving authorized assets.
package handler

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/nawaphonOHM/resume-assets/internal/domain"
	"github.com/nawaphonOHM/whatever/v3/pkg/logging"
	"github.com/nawaphonOHM/whatever/v3/pkg/rest"
)

const (
	codeAssetNotFound   = "ASSET_NOT_FOUND"
	codeAccessDenied    = "ASSET_ACCESS_DENIED"
	codeInvalidName     = "INVALID_ASSET_FILENAME"
	codeInvalidCriteria = "INVALID_ASSET_CRITERIA"
	codeAmbiguousAsset  = "AMBIGUOUS_ASSET_MATCH"
	codeInternalError   = "INTERNAL_SERVER_ERROR"
	codeServiceOffline  = "SERVICE_UNAVAILABLE"
)

// AssetHandler serves asset requests through the application service.
type AssetHandler struct {
	assetService domain.AssetService
}

// NewAssetHandler creates a handler backed by the supplied asset service.
func NewAssetHandler(assetService domain.AssetService) *AssetHandler {
	return &AssetHandler{assetService: assetService}
}

// ServeAsset handles GET /assets/:filename.
func (h *AssetHandler) ServeAsset(c rest.Context) rest.Response {
	if !h.isConfigured() {
		return rest.InternalServerError(codeServiceOffline, "asset service is not configured")
	}
	return h.serveRequest(c)
}

func (h *AssetHandler) isConfigured() bool {
	return h != nil && h.assetService != nil
}

func (h *AssetHandler) serveRequest(c rest.Context) rest.Response {
	ctx := requestContext(c)
	filename := requestFilename(c)
	if !validFilename(filename) {
		return rest.BadRequest(codeInvalidName, "filename must be a single path segment")
	}
	criteria, err := requestAssetCriteria(c)
	if err != nil {
		return rest.BadRequest(codeInvalidCriteria, "asset criteria must be one URL-encoded JSON object")
	}
	return h.retrieveAsset(ctx, filename, criteria)
}

func requestFilename(c rest.Context) string {
	if c == nil {
		return ""
	}
	return c.Param("filename")
}

func (h *AssetHandler) retrieveAsset(ctx context.Context, filename string, criteria *string) rest.Response {
	asset, err := h.assetService.GetAsset(ctx, filename, criteria)
	if err != nil {
		return mapAssetError(ctx, err)
	}
	if !hasAssetBody(asset) {
		logging.ErrorContext(ctx, "asset service returned an empty stream")
		return rest.InternalServerError(codeInternalError, "the server failed to serve the requested asset")
	}
	return buildAssetResponse(asset)
}

func hasAssetBody(asset *domain.StaticAsset) bool {
	return asset != nil && asset.Body != nil
}

func validFilename(filename string) bool {
	return filename != "" && !isDotSegment(filename) && !strings.ContainsAny(filename, "/\x00")
}

func isDotSegment(filename string) bool {
	return filename == "." || filename == ".."
}

func requestContext(c rest.Context) context.Context {
	if c != nil && c.Context() != nil {
		return c.Context()
	}
	return context.Background()
}

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
