package handler

import (
	"context"
	"errors"
	"io"

	"github.com/nawaphonOHM/resume-assets/internal/domain"
	"github.com/nawaphonOHM/whatever/v3/pkg/rest"
)

type handlerContextEmbedder struct {
	rest.Context
}

type handlerTestRequest struct {
	queryBindErr error
	ctx          context.Context
	query        map[string][]string
	filename     string
}

type handlerTestContext struct {
	request *handlerTestRequest
	handlerContextEmbedder
}

func (c *handlerTestContext) Param(name string) string {
	if c.request != nil && name == "filename" {
		return c.request.filename
	}
	return ""
}

func (c *handlerTestContext) Context() context.Context {
	if c.request == nil {
		return nil
	}
	return c.request.ctx
}

func (c *handlerTestContext) ShouldBindQuery(target any) error {
	query, ok := target.(*map[string][]string)
	if !ok {
		return errors.New("query binding target must be a map of string slices")
	}
	return copyTestQuery(query, c.request)
}

func copyTestQuery(query *map[string][]string, request *handlerTestRequest) error {
	if request == nil {
		return nil
	}
	if request.queryBindErr != nil {
		return request.queryBindErr
	}
	copyTestQueryValues(query, request.query)
	return nil
}

func copyTestQueryValues(query *map[string][]string, values map[string][]string) {
	for key, queryValues := range values {
		(*query)[key] = append([]string(nil), queryValues...)
	}
}

type handlerServiceCall struct {
	criteria *string
	ctx      context.Context
	filename string
}

type assetServiceFake struct {
	err      error
	call     *handlerServiceCall
	asset    *domain.StaticAsset
	document *domain.AssetDocument
	calls    int
}

func (s *assetServiceFake) FindAllowedAsset(
	ctx context.Context,
	filename string,
	criteria *string,
) (*domain.AssetDocument, error) {
	s.calls++
	s.call = &handlerServiceCall{ctx: ctx, filename: filename, criteria: criteria}
	return s.document, s.err
}

func (s *assetServiceFake) GetAsset(
	ctx context.Context,
	filename string,
	criteria *string,
) (*domain.StaticAsset, *domain.AssetDocument, error) {
	s.calls++
	s.call = &handlerServiceCall{ctx: ctx, filename: filename, criteria: criteria}
	return s.asset, s.document, s.err
}

type handlerBody struct{}

func (handlerBody) Read([]byte) (int, error) { return 0, io.EOF }
func (handlerBody) Close() error             { return nil }

type handlerTestContextKey string
