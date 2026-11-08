// Package mongodb implements MongoDB access for asset metadata.
package mongodb

import (
	"context"
	"errors"
	"fmt"

	"github.com/nawaphonOHM/resume-assets/internal/domain"
	"github.com/nawaphonOHM/whatever/v3/pkg/mongodb"
)

var ErrNilCollection = errors.New("mongodb collection is nil")

// MongoAssetRepository queries the assets collection for asset metadata.
type MongoAssetRepository struct {
	collection *mongodb.Collection
}

// NewMongoAssetRepository creates a repository backed by the supplied collection.
func NewMongoAssetRepository(collection *mongodb.Collection) *MongoAssetRepository {
	return &MongoAssetRepository{collection: collection}
}

// BuildFilter builds an exact filename match for the assets collection.
func BuildFilter(filename string) mongodb.D {
	return mongodb.D{{Key: "filename", Value: filename}}
}

// FindAsset queries the collection using the filename exactly as provided.
func (r *MongoAssetRepository) FindAsset(ctx context.Context, filename string) (*domain.AssetDocument, error) {
	if r.collection == nil {
		return nil, ErrNilCollection
	}

	var record assetRecord
	if err := r.collection.FindOne(ctx, BuildFilter(filename)).Decode(&record); err != nil {
		return nil, mapFindError(err, filename)
	}

	return record.toAssetDocument(), nil
}

func mapFindError(err error, filename string) error {
	if errors.Is(err, mongodb.ErrNoDocuments) {
		return domain.ErrAssetNotFound
	}
	return fmt.Errorf("find asset metadata for filename %q: %w", filename, err)
}
