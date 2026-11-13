package handler

import (
	"slices"
	"strings"
)

func validatedAssetHeaders(headers map[string]any) map[string]string {
	result := make(map[string]string, len(headers))
	seen := make(map[string]struct{}, len(headers))
	for _, name := range sortedAssetHeaderNames(headers) {
		addValidatedAssetHeader(result, seen, name, headers[name])
	}
	return result
}

func addValidatedAssetHeader(result map[string]string, seen map[string]struct{}, name string, candidate any) {
	value, ok := validatedAssetHeader(name, candidate)
	if !ok {
		return
	}
	key := strings.ToLower(name)
	if _, exists := seen[key]; exists {
		return
	}
	result[name] = value
	seen[key] = struct{}{}
}

func sortedAssetHeaderNames(headers map[string]any) []string {
	names := make([]string, 0, len(headers))
	for name := range headers {
		names = append(names, name)
	}
	slices.Sort(names)
	return names
}

func sortedHeaderNames(headers map[string]string) []string {
	names := make([]string, 0, len(headers))
	for name := range headers {
		names = append(names, name)
	}
	slices.Sort(names)
	return names
}

func validatedAssetHeader(name string, candidate any) (string, bool) {
	if !safeAssetHeaderName(name) {
		return "", false
	}
	value, ok := assetHeaderValue(candidate)
	return value, ok && safeAssetHeaderValue(name, value)
}

func safeAssetHeaderName(name string) bool {
	return validHeaderName(name) && !isHopByHopHeader(name)
}

func safeAssetHeaderValue(name, value string) bool {
	return validHeaderValue(value) &&
		(!strings.EqualFold(name, contentLengthHeaderName) || validContentLength(value))
}
