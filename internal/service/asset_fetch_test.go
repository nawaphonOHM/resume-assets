package service

import (
	"context"
	"testing"
	"time"

	"github.com/nawaphonOHM/resume-assets/internal/domain"
)

const (
	testAssetYear          = 2026
	testAssetContentLength = int64(42)
	testAssetHour          = 3
	testAssetMinute        = 4
	testAssetSecond        = 5
	testAssetObjectKey     = "assets/app.js"
)

type locationCase struct {
	name     string
	location string
	bucket   string
	key      string
}

func TestGetAssetResolvesLocationAndPreservesStream(t *testing.T) {
	tests := []locationCase{
		{name: "s3 URI", location: "s3://asset-bucket/assets/app.js", bucket: "asset-bucket", key: testAssetObjectKey},
		{name: "gs URI", location: "gs://asset-bucket/assets/app.js", bucket: "asset-bucket", key: testAssetObjectKey},
		{
			name: "s3 URI uses default bucket", location: "s3:///assets/app.js",
			bucket: testDefaultBucket, key: testAssetObjectKey,
		},
		{
			name: "gs URI uses default bucket", location: "gs:///assets/app.js",
			bucket: testDefaultBucket, key: testAssetObjectKey,
		},
		{
			name: "plain key uses default bucket", location: "/assets/app.js",
			bucket: testDefaultBucket, key: testAssetObjectKey,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			runLocationCase(t, test)
		})
	}
}

func runLocationCase(t *testing.T, test locationCase) {
	t.Helper()
	fixture := newLocationFixture(t, test)
	asset, document := getAssetForTest(fixture.ctx, t, fixture.service)
	assertLocationDocument(t, document, fixture.repository.document)
	assertLocationRepositoryCall(fixture.ctx, t, fixture.repository)
	assertStorageInput(t, fixture.storageClient, test.bucket, test.key)
	assertStorageContext(fixture.ctx, t, fixture.storageClient.call)
	assertStreamOutput(t, asset, fixture.body, fixture.fixedTime)
}

type locationFixture struct {
	ctx           context.Context
	fixedTime     time.Time
	body          *closeTracker
	storageClient *storageFake
	repository    *repositoryFake
	service       *AssetService
}

func newLocationFixture(t *testing.T, test locationCase) locationFixture {
	fixedTime := time.Date(testAssetYear, time.January, 2, testAssetHour, testAssetMinute, testAssetSecond, 0, time.UTC)
	body := &closeTracker{}
	storageClient := newStorageFixture(body, fixedTime)
	repository := repositoryForLocation(test.location)
	return locationFixture{
		ctx:       context.WithValue(context.Background(), testContextKey("request"), "request"),
		fixedTime: fixedTime, body: body, storageClient: storageClient,
		repository: repository, service: newServiceForTest(t, repository, storageClient),
	}
}

func assertLocationDocument(t *testing.T, got, want *domain.AssetDocument) {
	t.Helper()
	if got != want {
		t.Fatal("GetAsset did not return the document used for the storage lookup")
	}
}

func assertLocationRepositoryCall(ctx context.Context, t *testing.T, repository *repositoryFake) {
	t.Helper()
	if repository.calls != 1 {
		t.Fatalf("expected one metadata lookup, got %d", repository.calls)
	}
	assertRepositoryCall(ctx, t, repository.call, testAssetFilename)
}
