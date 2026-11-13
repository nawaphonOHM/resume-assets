// Package domain defines the asset metadata models and application interfaces.
package domain

import (
	"context"
	"errors"
	"io"
	"time"

	"github.com/nawaphonOHM/whatever/v3/pkg/storage"
)

var (
	ErrAssetNotFound        = errors.New("asset not found")
	ErrAssetAccessDenied    = errors.New("asset access denied")
	ErrInvalidAssetCriteria = errors.New("invalid asset criteria")
	ErrAmbiguousAssetMatch  = errors.New("asset match is ambiguous")
	ErrInvalidAssetMetadata = errors.New("invalid asset metadata")
)

// AssetCriteria is an exact JSON object matched against an asset's tags.
// Nested objects use map[string]any, arrays use []any, and numbers use json.Number.
type AssetCriteria map[string]any

// AssetDocument is the metadata stored in the MongoDB assets collection.
type AssetDocument struct {
	Metadata AssetMetadata `bson:"metaData"`
	Filename string        `bson:"filename"`
	Location string        `bson:"location"`
	Access   string        `bson:"access"`
}

// AssetMetadata contains optional response headers stored with an asset.
type AssetMetadata struct {
	ForCaching     map[string]any `bson:"forCaching"`
	ForDownloading map[string]any `bson:"forDownloading"`
}

// StaticAsset contains a streamed object and the metadata needed for its response.
type StaticAsset struct {
	LastModified  time.Time
	Body          io.ReadCloser
	ContentType   string
	ETag          string
	ContentLength int64
}

// AssetRepository looks up asset metadata by its exact filename and optional tags.
type AssetRepository interface {
	// FindAsset resolves metadata by the exact stored filename and optional exact tags.
	FindAsset(context.Context, string, *AssetCriteria) (*AssetDocument, error)
}

// ObjectStorage retrieves object streams and metadata from cloud storage.
type ObjectStorage interface {
	// GetObject retrieves an object stream and its available metadata.
	GetObject(context.Context, storage.GetObjectInput) (*storage.GetObjectOutput, error)
}

// AssetService retrieves an allowed asset by filename and optional JSON criteria.
type AssetService interface {
	// FindAllowedAsset looks up an allowed asset without retrieving its object.
	FindAllowedAsset(context.Context, string, *string) (*AssetDocument, error)

	// GetAsset retrieves an asset after applying the access policy.
	GetAsset(context.Context, string, *string) (*StaticAsset, *AssetDocument, error)
}
