package handler

import (
	"context"
	"net/http"

	"github.com/nawaphonOHM/resume-assets/internal/domain"
	"github.com/nawaphonOHM/whatever/v3/pkg/logging"
	"github.com/nawaphonOHM/whatever/v3/pkg/rest"
)

func (h *AssetHandler) serveHeadRequest(request assetRequest) rest.Response {
	asset, err := h.assetService.FindAllowedAsset(request.ctx, request.filename, request.criteria)
	if err != nil {
		return bodylessResponse(mapAssetError(request.ctx, err).StatusCode())
	}
	return headAssetResponse(request.ctx, request.filename, asset)
}

func headAssetResponse(ctx context.Context, filename string, asset *domain.AssetDocument) rest.Response {
	if asset == nil {
		logging.ErrorContext(ctx, "asset service returned an empty metadata document")
		return bodylessResponse(http.StatusInternalServerError)
	}
	if asset.Access != "ALLOW" {
		return bodylessResponse(http.StatusForbidden)
	}
	return buildHeadAssetResponse(filename, asset)
}
