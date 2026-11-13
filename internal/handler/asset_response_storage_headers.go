package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/nawaphonOHM/resume-assets/internal/domain"
)

const storageAssetHeaderCapacity = 4

func storageAssetHeaders(asset *domain.StaticAsset) map[string]string {
	if asset == nil {
		return nil
	}
	headers := make(map[string]string, storageAssetHeaderCapacity)
	addStorageHeader(headers, "Content-Type", asset.ContentType)
	addStorageContentLength(headers, asset.ContentLength)
	addStorageHeader(headers, "ETag", asset.ETag)
	addStorageLastModified(headers, asset.LastModified)
	return headers
}

func addStorageHeader(headers map[string]string, name, value string) {
	if value != "" {
		headers[name] = value
	}
}

func addStorageContentLength(headers map[string]string, length int64) {
	if length >= 0 {
		headers[contentLengthHeaderName] = strconv.FormatInt(length, decimalBase)
	}
}

func addStorageLastModified(headers map[string]string, modified time.Time) {
	if !modified.IsZero() {
		headers["Last-Modified"] = modified.UTC().Format(http.TimeFormat)
	}
}
