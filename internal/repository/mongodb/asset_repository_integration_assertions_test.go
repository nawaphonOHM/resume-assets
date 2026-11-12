package mongodb

import (
	"context"
	"testing"

	"github.com/nawaphonOHM/resume-assets/internal/domain"
)

func assertMongoIntegrationBehavior(ctx context.Context, t *testing.T, repository *MongoAssetRepository) {
	t.Helper()
	assertMongoSingleVariant(ctx, t, repository)
	assertMongoNonmatchingVariants(ctx, t, repository)
	assertMongoAmbiguousVariants(ctx, t, repository)
	assertMongoEmptyAndNullTags(ctx, t, repository)
}

func assertMongoSingleVariant(ctx context.Context, t *testing.T, repository *MongoAssetRepository) {
	t.Helper()
	criteria := domain.AssetCriteria{
		integrationVariantKey: "single",
		integrationNestedKey:  integrationNestedCriteria(),
		integrationStatusKey:  integrationStatusAvailable,
	}
	asset, err := repository.FindAsset(ctx, integrationSameFilename, &criteria)
	if err != nil {
		t.Fatalf("find one nested, reordered-object match: %v", err)
	}
	assertMongoSelectedAsset(t, asset)
}

func assertMongoSelectedAsset(t *testing.T, asset *domain.AssetDocument) {
	t.Helper()
	if asset.Filename != integrationSameFilename || asset.Access != "ALLOW" {
		t.Fatalf("unexpected selected asset: %#v", asset)
	}
}

func assertMongoNonmatchingVariants(ctx context.Context, t *testing.T, repository *MongoAssetRepository) {
	t.Helper()
	arrayOrder := domain.AssetCriteria{
		integrationVariantKey: "array-order",
		integrationNestedKey:  integrationNestedCriteria(),
		integrationStatusKey:  integrationStatusAvailable,
	}
	extra := domain.AssetCriteria{
		integrationVariantKey: "single",
		integrationNestedKey:  integrationNestedCriteria(),
		integrationStatusKey:  integrationStatusAvailable,
		"extra":               true,
	}
	missing := domain.AssetCriteria{integrationVariantKey: "missing", integrationStatusKey: integrationStatusAvailable}
	for _, criteria := range []*domain.AssetCriteria{&arrayOrder, &extra, &missing} {
		assertMongoRepositoryError(ctx, t, repository, mongoRepositoryQuery{
			filename: integrationSameFilename, criteria: criteria, want: domain.ErrAssetNotFound,
		})
	}
	assertMongoRepositoryError(ctx, t, repository, mongoRepositoryQuery{
		filename: "not-found.js", want: domain.ErrAssetNotFound,
	})
}

func assertMongoAmbiguousVariants(ctx context.Context, t *testing.T, repository *MongoAssetRepository) {
	t.Helper()
	criteria := domain.AssetCriteria{
		integrationVariantKey: "duplicate",
		integrationNestedKey:  integrationNestedCriteria(),
		integrationStatusKey:  integrationStatusAvailable,
	}
	assertMongoRepositoryError(ctx, t, repository, mongoRepositoryQuery{
		filename: integrationSameFilename, criteria: &criteria, want: domain.ErrAmbiguousAssetMatch,
	})
	assertMongoRepositoryError(ctx, t, repository, mongoRepositoryQuery{
		filename: integrationSameFilename, want: domain.ErrAmbiguousAssetMatch,
	})
}

func assertMongoEmptyAndNullTags(ctx context.Context, t *testing.T, repository *MongoAssetRepository) {
	t.Helper()
	empty := domain.AssetCriteria{}
	assertMongoRepositoryMatch(ctx, t, repository, mongoRepositoryQuery{filename: "empty-tags.js", criteria: &empty})
	assertMongoRepositoryError(ctx, t, repository, mongoRepositoryQuery{
		filename: "missing-tags.js", criteria: &empty, want: domain.ErrAssetNotFound,
	})
	null := domain.AssetCriteria{"nullable": nil}
	assertMongoRepositoryMatch(ctx, t, repository, mongoRepositoryQuery{
		filename: "null-vs-missing.js", criteria: &null,
	})
}

func assertMongoRepositoryMatch(
	ctx context.Context,
	t *testing.T,
	repository *MongoAssetRepository,
	query mongoRepositoryQuery,
) {
	t.Helper()
	if _, err := repository.FindAsset(ctx, query.filename, query.criteria); err != nil {
		t.Fatalf("expected a match for %q: %v", query.filename, err)
	}
}
