package handler

import (
	"reflect"
	"testing"

	"github.com/nawaphonOHM/resume-assets/internal/domain"
)

func TestAssetResponseHeaderSelectionAndValidation(t *testing.T) {
	asset := &domain.StaticAsset{ContentLength: -1}
	document := &domain.AssetDocument{Metadata: domain.AssetMetadata{
		ForCaching: map[string]any{
			handlerHeadModeHeader:      "cache",
			"Connection":               "close",
			"Bad Header":               "unsafe name",
			"X-Control":                "bad\r\nvalue",
			handlerContentLengthHeader: "not-a-number",
		},
		ForDownloading: map[string]any{
			handlerHeadModeHeader: "download",
		},
	}}

	for _, test := range []struct {
		want     map[string]string
		filename string
	}{
		{filename: assetResponseTestFilename, want: map[string]string{handlerHeadModeHeader: "cache"}},
		{filename: "app.PDF", want: map[string]string{handlerHeadModeHeader: "download"}},
	} {
		t.Run(test.filename, func(t *testing.T) {
			got := assetResponseHeadersForTest(test.filename, asset, document, nil)
			if !reflect.DeepEqual(got, test.want) {
				t.Fatalf("selected and validated headers: got %#v, want %#v", got, test.want)
			}
		})
	}
}

func TestAssetResponseHeadersWithoutMetadataKeepStorageHeaders(t *testing.T) {
	asset := &domain.StaticAsset{ContentType: "application/octet-stream", ContentLength: 0}
	got := assetResponseHeadersForTest(assetResponseTestFilename, asset, &domain.AssetDocument{}, nil)
	want := map[string]string{
		"Content-Type":   "application/octet-stream",
		"Content-Length": "0",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("storage headers without metadata: got %#v, want %#v", got, want)
	}
}
