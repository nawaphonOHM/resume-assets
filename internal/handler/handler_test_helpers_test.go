package handler

import (
	"context"
	"io"

	"github.com/nawaphonOHM/resume-assets/internal/domain"
	"github.com/nawaphonOHM/whatever/v3/pkg/rest"
)

type handlerContextEmbedder struct {
	rest.Context
}

type handlerTestRequest struct {
	ctx      context.Context
	filename string
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

type handlerServiceCall struct {
	ctx      context.Context
	filename string
}

type assetServiceFake struct {
	err   error
	call  *handlerServiceCall
	asset *domain.StaticAsset
	calls int
}

func (s *assetServiceFake) GetAsset(ctx context.Context, filename string) (*domain.StaticAsset, error) {
	s.calls++
	s.call = &handlerServiceCall{ctx: ctx, filename: filename}
	return s.asset, s.err
}

type handlerBody struct{}

func (handlerBody) Read([]byte) (int, error) { return 0, io.EOF }
func (handlerBody) Close() error             { return nil }

type handlerTestContextKey string
