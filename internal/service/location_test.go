package service

import (
	"errors"
	"testing"
)

const (
	testLocationBucket = "asset-bucket"
	testLocationKey    = "path/app.js"
)

type parseLocationCase struct {
	wantErr       error
	name          string
	location      string
	defaultBucket string
	wantBucket    string
	wantKey       string
}

func TestParseLocation(t *testing.T) {
	tests := []parseLocationCase{
		{
			name:       "s3 URI",
			location:   "s3://asset-bucket/path/app.js",
			wantBucket: testLocationBucket,
			wantKey:    testLocationKey,
		},
		{
			name:       "gs URI",
			location:   "gs://asset-bucket/path/app.js",
			wantBucket: testLocationBucket,
			wantKey:    testLocationKey,
		},
		{
			name:          "s3 URI with default bucket",
			location:      "s3:///path/app.js",
			defaultBucket: testDefaultBucket,
			wantBucket:    testDefaultBucket,
			wantKey:       testLocationKey,
		},
		{
			name:          "gs URI with default bucket",
			location:      "gs:///path/app.js",
			defaultBucket: testDefaultBucket,
			wantBucket:    testDefaultBucket,
			wantKey:       testLocationKey,
		},
		{
			name:          "plain key with leading slash",
			location:      "///path/app.js",
			defaultBucket: testDefaultBucket,
			wantBucket:    testDefaultBucket,
			wantKey:       testLocationKey,
		},
		{
			name:       "trimmed location",
			location:   " s3://asset-bucket/app.js ",
			wantBucket: testLocationBucket,
			wantKey:    testAssetFilename,
		},
		{
			name:       "explicit bucket with no default",
			location:   "s3://asset-bucket/app.js",
			wantBucket: testLocationBucket,
			wantKey:    testAssetFilename,
		},
		{
			name:     "empty default bucket follows reference semantics",
			location: "gs:///app.js",
			wantKey:  testAssetFilename,
		},
		{name: "empty location", wantErr: ErrInvalidLocation},
		{name: "unsupported scheme", location: "https://bucket/app.js", wantErr: ErrInvalidLocation},
		{name: "URI without key", location: "s3://asset-bucket", wantErr: ErrInvalidLocation},
		{name: "empty URI key", location: "s3://asset-bucket///", wantErr: ErrInvalidLocation},
		{name: "empty root key", location: "s3:///", wantErr: ErrInvalidLocation},
		{name: "empty plain key", location: "///", wantErr: ErrInvalidLocation},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assertLocation(t, test)
		})
	}
}

func assertLocation(t *testing.T, test parseLocationCase) {
	t.Helper()
	bucket, key, err := ParseLocation(test.location, test.defaultBucket)
	assertLocationError(t, err, test.wantErr)
	assertLocationTarget(t, bucket, key, test)
}

func assertLocationError(t *testing.T, got, want error) {
	t.Helper()
	if !errors.Is(got, want) {
		t.Fatalf("expected error %v, got %v", want, got)
	}
}

func assertLocationTarget(t *testing.T, bucket, key string, test parseLocationCase) {
	t.Helper()
	if bucket != test.wantBucket || key != test.wantKey {
		t.Fatalf("got bucket=%q key=%q, want bucket=%q key=%q", bucket, key, test.wantBucket, test.wantKey)
	}
}
