package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/nawaphonOHM/resume-assets/internal/domain"
	"github.com/nawaphonOHM/whatever/v3/pkg/rest"
)

const (
	decimalBase             = 10
	contentLengthHeaderName = "Content-Length"
)

func buildAssetResponse(asset *domain.StaticAsset) rest.Response {
	if asset == nil || asset.Body == nil {
		return rest.InternalServerError(codeInternalError, "the server failed to serve the requested asset")
	}
	response := rest.Stream(http.StatusOK, asset.Body)
	response = withOptionalHeader(response, "Content-Type", asset.ContentType)
	response = withOptionalContentLength(response, asset.ContentLength)
	response = withOptionalHeader(response, "ETag", asset.ETag)
	response = withOptionalLastModified(response, asset.LastModified)
	return response
}

func withOptionalHeader(response rest.Response, name, value string) rest.Response {
	if value != "" {
		return response.WithHeader(name, value)
	}
	return response
}

func withOptionalContentLength(response rest.Response, length int64) rest.Response {
	if length >= 0 {
		return response.WithHeader(contentLengthHeaderName, strconv.FormatInt(length, decimalBase))
	}
	return response
}

func withOptionalLastModified(response rest.Response, lastModified time.Time) rest.Response {
	if !lastModified.IsZero() {
		return response.WithHeader("Last-Modified", lastModified.UTC().Format(http.TimeFormat))
	}
	return response
}
