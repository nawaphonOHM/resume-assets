package mongodb

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/nawaphonOHM/resume-assets/internal/domain"
	libmongodb "github.com/nawaphonOHM/whatever/v3/pkg/testing/mongodb"
)

const (
	mongoIntegrationTimeout = 30 * time.Second
	mongoCleanupTimeout     = 10 * time.Second
)

func newMongoIntegrationClient(t *testing.T, uri string) (context.Context, *libmongodb.TestClient, string) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), mongoIntegrationTimeout)
	t.Cleanup(cancel)
	databaseName := fmt.Sprintf("resume_assets_criteria_test_%d", time.Now().UnixNano())
	client := libmongodb.NewTestClientURI(t, uri, libmongodb.WithDatabase(databaseName))
	t.Cleanup(func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), mongoCleanupTimeout)
		defer cleanupCancel()
		if err := client.DropDatabase(cleanupCtx, databaseName); err != nil {
			t.Errorf("drop integration database: %v", err)
		}
	})
	if err := client.Ping(ctx); err != nil {
		t.Fatalf("ping test MongoDB: %v", err)
	}
	return ctx, client, databaseName
}

func integrationAsset(filename string, tags libmongodb.D) libmongodb.D {
	document := libmongodb.D{
		{Key: "filename", Value: filename},
		{Key: "location", Value: "s3://assets/" + filename},
		{Key: "access", Value: "ALLOW"},
	}
	if tags != nil {
		document = append(document, libmongodb.E{Key: "tags", Value: tags})
	}
	return document
}

type mongoRepositoryQuery struct {
	criteria *domain.AssetCriteria
	want     error
	filename string
}

func assertMongoRepositoryError(
	ctx context.Context,
	t *testing.T,
	repository *MongoAssetRepository,
	query mongoRepositoryQuery,
) {
	t.Helper()
	asset, err := repository.FindAsset(ctx, query.filename, query.criteria)
	if asset != nil {
		t.Fatalf("expected no asset for %q, got %#v", query.filename, asset)
	}
	if !errors.Is(err, query.want) {
		t.Fatalf("expected error %v for %q, got %v", query.want, query.filename, err)
	}
}
