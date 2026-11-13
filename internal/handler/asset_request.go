package handler

import (
	"context"
	"strings"

	"github.com/nawaphonOHM/whatever/v3/pkg/rest"
)

type assetRequest struct {
	response rest.Response
	ctx      context.Context
	criteria *string
	filename string
	valid    bool
}

func parseAssetRequest(c rest.Context) assetRequest {
	ctx := requestContext(c)
	filename := requestFilename(c)
	if !validFilename(filename) {
		return assetRequest{
			ctx: ctx, response: rest.BadRequest(codeInvalidName, "filename must be a single path segment"),
		}
	}
	criteria, err := requestAssetCriteria(c)
	if err != nil {
		return assetRequest{
			ctx: ctx,
			response: rest.BadRequest(
				codeInvalidCriteria,
				"asset criteria must be one URL-encoded JSON object",
			),
		}
	}
	return assetRequest{ctx: ctx, filename: filename, criteria: criteria, valid: true}
}

func requestFilename(c rest.Context) string {
	if c == nil {
		return ""
	}
	return c.Param("filename")
}

func validFilename(filename string) bool {
	return filename != "" && !isDotSegment(filename) && !strings.ContainsAny(filename, "/\x00")
}

func isDotSegment(filename string) bool {
	return filename == "." || filename == ".."
}

func requestContext(c rest.Context) context.Context {
	if c != nil && c.Context() != nil {
		return c.Context()
	}
	return context.Background()
}
