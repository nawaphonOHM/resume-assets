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
	ErrInvalidAssetMetadata = errors.New("invalid asset metadata")
)

// AssetDocument is the metadata stored in the MongoDB assets collection.
type AssetDocument struct {
	Filename string `bson:"filename"`
	Location string `bson:"location"`
	Access   string `bson:"access"`
}

// StaticAsset contains a streamed object and the metadata needed for its response.
type StaticAsset struct {
	LastModified  time.Time
	Body          io.ReadCloser
	ContentType   string
	ETag          string
	ContentLength int64
}

// AssetRepository looks up asset metadata by its exact filename.
type AssetRepository interface {
	// FindAsset resolves metadata by the exact stored filename.
	FindAsset(context.Context, string) (*AssetDocument, error)
}

// ObjectStorage retrieves object streams and metadata from cloud storage.
type ObjectStorage interface {
	// GetObject retrieves an object stream and its available metadata.
	GetObject(context.Context, storage.GetObjectInput) (*storage.GetObjectOutput, error)
}

// AssetService retrieves an allowed asset by filename.
type AssetService interface {
	// GetAsset retrieves an asset after applying the access policy.
	GetAsset(context.Context, string) (*StaticAsset, error)
}
