package handler

import (
	"strconv"
	"strings"

	"github.com/nawaphonOHM/resume-assets/internal/domain"
)

const (
	decimalBase             = 10
	contentLengthHeaderName = "Content-Length"
)

type namedAssetHeader struct {
	name  string
	value string
}

type assetMetadataHeader struct {
	key   string
	name  string
	value string
}

func assetResponseHeaders(
	responseContext assetHeaderResponseContext,
	asset *domain.StaticAsset,
	document *domain.AssetDocument,
) map[string]string {
	var metadataHeaders map[string]string
	if document != nil {
		metadataHeaders = validatedAssetHeaders(selectedAssetHeaders(responseContext.filename, document.Metadata))
	}
	return mergeAssetResponseHeaders(responseContext, storageAssetHeaders(asset), metadataHeaders)
}

func mergeAssetResponseHeaders(
	responseContext assetHeaderResponseContext,
	storageHeaders, metadataHeaders map[string]string,
) map[string]string {
	merged := make(map[string]namedAssetHeader, len(storageHeaders)+len(metadataHeaders))
	for _, name := range sortedHeaderNames(storageHeaders) {
		merged[strings.ToLower(name)] = namedAssetHeader{name: name, value: storageHeaders[name]}
	}
	for _, name := range sortedHeaderNames(metadataHeaders) {
		mergeMetadataAssetHeader(responseContext, merged, name, metadataHeaders[name])
	}
	return responseHeadersByName(merged)
}

func mergeMetadataAssetHeader(
	responseContext assetHeaderResponseContext,
	merged map[string]namedAssetHeader,
	name, value string,
) {
	key := strings.ToLower(name)
	header := assetMetadataHeader{key: key, name: name, value: value}
	if strings.EqualFold(name, contentLengthHeaderName) {
		mergeMetadataContentLength(responseContext, merged, header)
		return
	}
	mergeOtherMetadataHeader(responseContext, merged, header)
}

func mergeMetadataContentLength(
	responseContext assetHeaderResponseContext,
	merged map[string]namedAssetHeader,
	header assetMetadataHeader,
) {
	storageHeader, exists := merged[header.key]
	if !exists {
		responseContext.warnUnknownContentLength(header.name)
		return
	}
	if contentLengthsDiffer(header.value, storageHeader.value) {
		responseContext.warnConflict(header.name, "storage")
	}
}

func contentLengthsDiffer(metadata, storage string) bool {
	metadataLength, metadataErr := strconv.ParseInt(metadata, decimalBase, headerIntegerBitSize)
	storageLength, storageErr := strconv.ParseInt(storage, decimalBase, headerIntegerBitSize)
	return metadataErr != nil || storageErr != nil || metadataLength != storageLength
}

func mergeOtherMetadataHeader(
	responseContext assetHeaderResponseContext,
	merged map[string]namedAssetHeader,
	header assetMetadataHeader,
) {
	storageHeader, exists := merged[header.key]
	if exists && header.value != storageHeader.value {
		responseContext.warnConflict(header.name, "metadata")
	}
	merged[header.key] = namedAssetHeader{name: header.name, value: header.value}
}
