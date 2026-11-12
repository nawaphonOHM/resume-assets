package mongodb

import (
	"context"
	"errors"
)

func readAssetRecords(ctx context.Context, cursor interface {
	Next(context.Context) bool
	Decode(any) error
	Err() error
	Close(context.Context) error
}) (records []assetRecord, err error) {
	defer func() {
		err = errors.Join(err, cursor.Close(ctx))
	}()
	records = make([]assetRecord, 0, 2)
	for cursor.Next(ctx) {
		record, decodeErr := decodeAssetCursorRecord(cursor)
		if decodeErr != nil {
			return nil, decodeErr
		}
		records = append(records, record)
	}
	return records, cursor.Err()
}

func decodeAssetCursorRecord(cursor interface{ Decode(any) error }) (assetRecord, error) {
	var record assetRecord
	if err := cursor.Decode(&record); err != nil {
		return assetRecord{}, err
	}
	return record, nil
}
