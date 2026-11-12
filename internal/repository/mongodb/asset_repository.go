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

// FindAsset queries by exact filename and optional exact tags, rejecting ambiguous matches.
func (r *MongoAssetRepository) FindAsset(
	ctx context.Context,
	filename string,
	criteria *domain.AssetCriteria,
) (*domain.AssetDocument, error) {
	if r.collection == nil {
		return nil, ErrNilCollection
	}
	records, err := r.findRecords(ctx, filename, criteria)
	if err != nil {
		return nil, err
	}
	return uniqueAssetRecord(records, filename)
}

func (r *MongoAssetRepository) findRecords(
	ctx context.Context,
	filename string,
	criteria *domain.AssetCriteria,
) ([]assetRecord, error) {
	filter, err := BuildFilter(filename, criteria)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", domain.ErrInvalidAssetCriteria, err)
	}
	return r.aggregateAssetRecords(ctx, filename, filter)
}

func (r *MongoAssetRepository) aggregateAssetRecords(
	ctx context.Context,
	filename string,
	filter mongodb.D,
) ([]assetRecord, error) {
	pipeline := mongodb.A{
		mongodb.D{{Key: "$match", Value: filter}},
		mongodb.D{{Key: "$limit", Value: 2}},
	}
	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, mapFindError(err, filename)
	}
	records, err := readAssetRecords(ctx, cursor)
	if err != nil {
		return nil, fmt.Errorf("read asset metadata for filename %q: %w", filename, err)
	}
	return records, nil
}

func uniqueAssetRecord(records []assetRecord, filename string) (*domain.AssetDocument, error) {
	if len(records) == 0 {
		return nil, domain.ErrAssetNotFound
	}
	if len(records) == 1 {
		return records[0].toAssetDocument(), nil
	}
	return nil, fmt.Errorf(
		"multiple asset metadata documents for filename %q: %w",
		filename,
		domain.ErrAmbiguousAssetMatch,
	)
}

func mapFindError(err error, filename string) error {
	if errors.Is(err, mongodb.ErrNoDocuments) {
		return domain.ErrAssetNotFound
	}
	return fmt.Errorf("find asset metadata for filename %q: %w", filename, err)
}
