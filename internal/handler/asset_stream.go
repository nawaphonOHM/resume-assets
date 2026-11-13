package handler

import (
	"context"

	"github.com/nawaphonOHM/resume-assets/internal/domain"
	"github.com/nawaphonOHM/whatever/v3/pkg/logging"
)

func hasAssetBody(asset *domain.StaticAsset) bool {
	return asset != nil && asset.Body != nil
}

func closeAssetStream(ctx context.Context, asset *domain.StaticAsset) {
	if !hasAssetBody(asset) {
		return
	}
	if err := asset.Body.Close(); err != nil {
		logging.ErrorContext(ctx, "failed to close asset stream after missing metadata document", "error", err)
	}
}
