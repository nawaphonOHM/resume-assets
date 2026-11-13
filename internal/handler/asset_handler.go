// Package handler provides HTTP handlers for serving authorized assets.
package handler

import (
	"context"
	"net/http"

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

// ServeAssetHead handles HEAD /assets/:filename without retrieving the object.
func (h *AssetHandler) ServeAssetHead(c rest.Context) rest.Response {
	if !h.isConfigured() {
		return bodylessResponse(http.StatusInternalServerError)
	}
	request := parseAssetRequest(c)
	if !request.valid {
		return bodylessResponse(request.response.StatusCode())
	}
	return h.serveHeadRequest(request)
}

func (h *AssetHandler) isConfigured() bool {
	return h != nil && h.assetService != nil
}

func (h *AssetHandler) serveRequest(c rest.Context) rest.Response {
	request := parseAssetRequest(c)
	if !request.valid {
		return request.response
	}
	return h.retrieveAsset(request.ctx, request.filename, request.criteria)
}

func (h *AssetHandler) retrieveAsset(ctx context.Context, filename string, criteria *string) rest.Response {
	asset, document, err := h.assetService.GetAsset(ctx, filename, criteria)
	if err != nil {
		return mapAssetError(ctx, err)
	}
	return retrievedAssetResponse(ctx, filename, asset, document)
}

func retrievedAssetResponse(
	ctx context.Context,
	filename string,
	asset *domain.StaticAsset,
	document *domain.AssetDocument,
) rest.Response {
	if document == nil {
		closeAssetStream(ctx, asset)
		logging.ErrorContext(ctx, "asset service returned an empty metadata document")
		return rest.InternalServerError(codeInternalError, "the server failed to serve the requested asset")
	}
	if !hasAssetBody(asset) {
		logging.ErrorContext(ctx, "asset service returned an empty stream")
		return rest.InternalServerError(codeInternalError, "the server failed to serve the requested asset")
	}
	return buildAssetResponse(ctx, filename, asset, document)
}
