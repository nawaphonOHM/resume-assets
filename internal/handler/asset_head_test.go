package handler

import (
	"maps"
	"testing"

	"github.com/nawaphonOHM/resume-assets/internal/domain"
)

const (
	handlerAllowedAccess       = "ALLOW"
	handlerETagHeader          = "ETag"
	handlerContentLengthHeader = "Content-Length"
	handlerHeadModeHeader      = "X-Mode"
	handlerAssetLength         = 1296
)

func TestServeAssetHeadSelectsMetadataHeaders(t *testing.T) {
	document := &domain.AssetDocument{
		Filename: handlerTestFilename,
		Access:   handlerAllowedAccess,
		Metadata: domain.AssetMetadata{
			ForCaching: map[string]any{
				"ETag":                "cache-tag",
				handlerHeadModeHeader: "cache",
			},
			ForDownloading: map[string]any{
				handlerContentLengthHeader: int32(handlerAssetLength),
				"Content-Type":             "application/pdf",
				handlerHeadModeHeader:      "download",
			},
		},
	}
	tests := []struct {
		wantHeaders map[string]string
		name        string
		filename    string
	}{
		{
			name: "non PDF uses caching metadata", filename: "app.js",
			wantHeaders: map[string]string{handlerETagHeader: "cache-tag", handlerHeadModeHeader: "cache"},
		},
		{
			name: "PDF extension is case insensitive", filename: "app.PDF",
			wantHeaders: map[string]string{
				handlerContentLengthHeader: "1296", "Content-Type": "application/pdf",
				handlerHeadModeHeader: "download",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := validatedAssetHeaders(selectedAssetHeaders(test.filename, document.Metadata))
			if !maps.Equal(got, test.wantHeaders) {
				t.Fatalf("selected headers: got %#v, want %#v", got, test.wantHeaders)
			}
		})
	}
}
