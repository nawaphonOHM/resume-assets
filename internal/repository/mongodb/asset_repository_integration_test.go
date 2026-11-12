package mongodb

import (
	"os"
	"testing"
)

func TestMongoAssetRepositoryExactTags(t *testing.T) {
	uri := os.Getenv("ASSET_TEST_MONGODB_URI")
	if uri == "" {
		t.Skip("set ASSET_TEST_MONGODB_URI to run MongoDB expression integration tests")
	}

	ctx, client, databaseName := newMongoIntegrationClient(t, uri)
	collection := client.Collection("assets", databaseName)
	repository := NewMongoAssetRepository(collection)
	insertMongoIntegrationFixtures(ctx, t, collection)
	assertMongoIntegrationBehavior(ctx, t, repository)
}
