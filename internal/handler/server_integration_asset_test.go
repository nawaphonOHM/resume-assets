package handler

import (
	"context"
	"errors"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/nawaphonOHM/resume-assets/internal/domain"
)

const (
	integrationAssetData     = "asset-data"
	integrationContentType   = "application/octet-stream"
	integrationETag          = `"asset-etag"`
	integrationPrivateDetail = "private backend detail"
	integrationAssetYear     = 2026
	integrationAssetHour     = 3
	integrationAssetMinute   = 4
	integrationAssetSecond   = 5
)

type routerAssetService struct {
	body      *integrationBody
	filenames []string
	mu        sync.Mutex
}

func (s *routerAssetService) GetAsset(ctx context.Context, filename string) (*domain.StaticAsset, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	s.recordFilename(filename)
	if err := integrationAssetError(filename); err != nil {
		return nil, err
	}
	return s.newAsset(), nil
}

func integrationAssetError(filename string) error {
	if filename == "missing" {
		return domain.ErrAssetNotFound
	}
	if filename == "denied" {
		return domain.ErrAssetAccessDenied
	}
	return integrationBackendError(filename)
}

func integrationBackendError(filename string) error {
	if filename == "broken" {
		return errors.New(integrationPrivateDetail)
	}
	return nil
}

func (s *routerAssetService) recordFilename(filename string) {
	s.mu.Lock()
	s.filenames = append(s.filenames, filename)
	s.mu.Unlock()
}

func (s *routerAssetService) newAsset() *domain.StaticAsset {
	body := &integrationBody{reader: strings.NewReader(integrationAssetData), closed: make(chan struct{})}
	s.mu.Lock()
	s.body = body
	s.mu.Unlock()
	return &domain.StaticAsset{
		Body:          body,
		ContentType:   integrationContentType,
		ContentLength: int64(len(integrationAssetData)),
		ETag:          integrationETag,
		LastModified: time.Date(
			integrationAssetYear, time.January, 2, integrationAssetHour, integrationAssetMinute,
			integrationAssetSecond, 0, time.UTC,
		),
	}
}

func (s *routerAssetService) snapshot() (*integrationBody, []string) {
	s.mu.Lock()
	body := s.body
	filenames := append([]string(nil), s.filenames...)
	s.mu.Unlock()
	return body, filenames
}

type integrationBody struct {
	reader io.Reader
	closed chan struct{}
	once   sync.Once
}

func (b *integrationBody) Read(p []byte) (int, error) { return b.reader.Read(p) }
func (b *integrationBody) Close() error {
	b.once.Do(func() { close(b.closed) })
	return nil
}
