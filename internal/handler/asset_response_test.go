package handler

import (
	"testing"

	"github.com/nawaphonOHM/resume-assets/internal/domain"
)

const (
	assetResponseTestFilename      = "app.js"
	assetResponseTestStorageSize   = int64(4)
	assetResponseTestMetadataSize  = 5
	assetResponseTestConflictCount = 3
)

func TestAssetResponseHeaderPrecedence(t *testing.T) {
	var warnings assetHeaderWarningRecorder
	asset := &domain.StaticAsset{
		ContentType:   "application/octet-stream",
		ContentLength: assetResponseTestStorageSize,
		ETag:          "storage-tag",
	}
	document := &domain.AssetDocument{Metadata: domain.AssetMetadata{ForCaching: map[string]any{
		"content-type":   "application/javascript",
		"content-length": "0005",
		"etag":           "metadata-tag",
		"X-Metadata":     "present",
	}}}

	got := assetResponseHeadersForTest(assetResponseTestFilename, asset, document, warnings.log)
	want := map[string]string{
		"content-type":   "application/javascript",
		"Content-Length": "4",
		"etag":           "metadata-tag",
		"X-Metadata":     "present",
	}
	assertResponseHeaders(t, got, want)
	assertWarningCount(t, warnings.warnings, assetResponseTestConflictCount)
	assertConflictWarning(t, warnings.warnings[0], "content-length", "storage")
	assertConflictWarning(t, warnings.warnings[1], "content-type", "metadata")
	assertConflictWarning(t, warnings.warnings[2], "etag", "metadata")
}

func TestAssetResponseHeaderMatchingContentLengthIsComparedNumerically(t *testing.T) {
	var warnings assetHeaderWarningRecorder
	asset := &domain.StaticAsset{ContentLength: assetResponseTestStorageSize}
	document := &domain.AssetDocument{Metadata: domain.AssetMetadata{ForCaching: map[string]any{
		"Content-Length": "0004",
	}}}

	got := assetResponseHeadersForTest(assetResponseTestFilename, asset, document, warnings.log)
	if got["Content-Length"] != "4" {
		t.Fatalf("storage content length was not retained: %#v", got)
	}
	assertWarningCount(t, warnings.warnings, 0)
}

func TestAssetResponseHeaderUnknownStorageLengthOmitsMetadataLength(t *testing.T) {
	var warnings assetHeaderWarningRecorder
	asset := &domain.StaticAsset{ContentLength: -1}
	document := &domain.AssetDocument{Metadata: domain.AssetMetadata{ForCaching: map[string]any{
		handlerContentLengthHeader: assetResponseTestMetadataSize,
	}}}

	got := assetResponseHeadersForTest(assetResponseTestFilename, asset, document, warnings.log)
	assertResponseHeaderAbsent(t, got, handlerContentLengthHeader)
	assertWarningCount(t, warnings.warnings, 1)
	warning := warnings.warnings[0]
	assertWarningAttributes(t, warning, handlerContentLengthHeader, "storage")
	if warning.message != assetUnknownLengthWarning {
		t.Fatalf("unexpected unknown-length warning: %q", warning.message)
	}
}
