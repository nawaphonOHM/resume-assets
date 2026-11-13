package handler

import (
	"strings"
	"time"

	"github.com/nawaphonOHM/resume-assets/internal/domain"
)

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

func (s *routerAssetService) criteriaSnapshot() []*string {
	s.mu.Lock()
	criteria := append([]*string(nil), s.criteria...)
	s.mu.Unlock()
	return criteria
}
