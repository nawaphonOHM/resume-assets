package handler

import (
	"io"
	"net/http"
	"path"
	"strings"

	"github.com/nawaphonOHM/resume-assets/internal/domain"
	"github.com/nawaphonOHM/whatever/v3/pkg/rest"
)

func buildHeadAssetResponse(filename string, asset *domain.AssetDocument) rest.Response {
	headers := selectedAssetHeaders(filename, asset.Metadata)
	return withAssetHeaders(bodylessResponse(http.StatusOK), validatedAssetHeaders(headers))
}

func selectedAssetHeaders(filename string, metadata domain.AssetMetadata) map[string]any {
	if strings.EqualFold(path.Ext(filename), ".pdf") {
		return metadata.ForDownloading
	}
	return metadata.ForCaching
}

func bodylessResponse(status int) rest.Response {
	return rest.Stream(status, io.NopCloser(strings.NewReader("")))
}

func withAssetHeaders(response rest.Response, headers map[string]string) rest.Response {
	result := response
	for _, name := range sortedHeaderNames(headers) {
		result = result.WithHeader(name, headers[name])
	}
	return result
}
