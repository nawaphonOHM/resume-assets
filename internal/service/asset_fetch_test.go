package service

import (
	"context"
	"testing"
	"time"
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
	fixedTime := time.Date(testAssetYear, time.January, 2, testAssetHour, testAssetMinute, testAssetSecond, 0, time.UTC)
	body := &closeTracker{}
	storageClient := newStorageFixture(body, fixedTime)
	ctx := context.WithValue(context.Background(), testContextKey("request"), "request")
	assetService := newServiceForTest(t, repositoryForLocation(test.location), storageClient)
	asset := getAssetForTest(ctx, t, assetService)
	assertStorageInput(t, storageClient, test.bucket, test.key)
	assertStorageContext(ctx, t, storageClient.call)
	assertStreamOutput(t, asset, body, fixedTime)
}
