package handler

import (
	"context"
	"sync"

	"github.com/nawaphonOHM/resume-assets/internal/domain"
)

const (
	integrationAssetData          = "asset-data"
	integrationContentType        = "application/octet-stream"
	integrationETag               = `"asset-etag"`
	integrationMetadataETag       = `"metadata-etag"`
	integrationPDFContentType     = "application/pdf"
	integrationHeaderContentType  = "Content-Type"
	integrationHeaderLastModified = "Last-Modified"
	integrationHeaderMode         = "X-Head-Mode"
	integrationPrivateDetail      = "private backend detail"
	integrationAssetYear          = 2026
	integrationAssetHour          = 3
	integrationAssetMinute        = 4
	integrationAssetSecond        = 5
)

type routerAssetService struct {
	body      *integrationBody
	filenames []string
	criteria  []*string
	mu        sync.Mutex
}

func (s *routerAssetService) FindAllowedAsset(
	ctx context.Context,
	filename string,
	criteria *string,
) (*domain.AssetDocument, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	s.recordRequest(filename, criteria)
	if err := integrationAssetError(filename); err != nil {
		return nil, err
	}
	return newIntegrationAssetDocument(filename), nil
}

func newIntegrationAssetDocument(filename string) *domain.AssetDocument {
	return &domain.AssetDocument{
		Filename: filename,
		Access:   "ALLOW",
		Metadata: domain.AssetMetadata{
			ForCaching: map[string]any{
				"etag":        integrationMetadataETag,
				"X-Head-Mode": "cache",
			},
			ForDownloading: map[string]any{
				"content-length": int32(len(integrationAssetData)),
				"content-type":   integrationPDFContentType,
				"X-Head-Mode":    "download",
			},
		},
	}
}

func (s *routerAssetService) GetAsset(
	ctx context.Context,
	filename string,
	criteria *string,
) (*domain.StaticAsset, *domain.AssetDocument, error) {
	if err := ctx.Err(); err != nil {
		return nil, nil, err
	}
	s.recordRequest(filename, criteria)
	if err := integrationAssetError(filename); err != nil {
		return nil, nil, err
	}
	return s.newAsset(), newIntegrationAssetDocument(filename), nil
}

func (s *routerAssetService) recordRequest(filename string, criteria *string) {
	s.mu.Lock()
	s.filenames = append(s.filenames, filename)
	s.criteria = append(s.criteria, cloneCriteria(criteria))
	s.mu.Unlock()
}

func cloneCriteria(criteria *string) *string {
	if criteria != nil {
		cloned := *criteria
		return &cloned
	}
	return nil
}
