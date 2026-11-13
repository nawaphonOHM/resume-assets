package handler

import (
	"context"
	"net/http"

	"github.com/nawaphonOHM/resume-assets/internal/domain"
	"github.com/nawaphonOHM/whatever/v3/pkg/logging"
	"github.com/nawaphonOHM/whatever/v3/pkg/rest"
)

func buildAssetResponse(
	ctx context.Context,
	filename string,
	asset *domain.StaticAsset,
	document *domain.AssetDocument,
) rest.Response {
	if asset == nil || asset.Body == nil {
		return rest.InternalServerError(codeInternalError, "the server failed to serve the requested asset")
	}
	response := rest.Stream(http.StatusOK, asset.Body)
	return withAssetHeaders(
		response,
		assetResponseHeaders(
			assetHeaderResponseContext{ctx: ctx, filename: filename, warn: logging.WarnContext},
			asset,
			document,
		),
	)
}
