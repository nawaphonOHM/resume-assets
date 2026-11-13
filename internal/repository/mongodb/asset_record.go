package mongodb

import (
	"github.com/nawaphonOHM/resume-assets/internal/domain"
	libmongodb "github.com/nawaphonOHM/whatever/v3/pkg/mongodb"
)

type assetRecord struct {
	Metadata any    `bson:"metaData"`
	Location any    `bson:"location"`
	Access   any    `bson:"access"`
	Filename string `bson:"filename"`
}

func (r assetRecord) toAssetDocument() *domain.AssetDocument {
	return &domain.AssetDocument{
		Filename: r.Filename,
		Location: stringField(r.Location),
		Access:   stringField(r.Access),
		Metadata: assetMetadata(r.Metadata),
	}
}

func assetMetadata(value any) domain.AssetMetadata {
	fields, ok := documentFields(value)
	if !ok {
		return domain.AssetMetadata{}
	}
	return domain.AssetMetadata{
		ForCaching:     headerFields(fields["forCaching"]),
		ForDownloading: headerFields(fields["forDownloading"]),
	}
}

func headerFields(value any) map[string]any {
	fields, ok := documentFields(value)
	if !ok {
		return nil
	}
	return fields
}

func documentFields(value any) (map[string]any, bool) {
	if fields, ok := value.(libmongodb.D); ok {
		return orderedDocumentFields(fields), true
	}
	return unorderedDocumentFields(value)
}

func unorderedDocumentFields(value any) (map[string]any, bool) {
	if fields, ok := value.(map[string]any); ok {
		return fields, true
	}
	return mongodbMapFields(value)
}

func mongodbMapFields(value any) (map[string]any, bool) {
	fields, ok := value.(libmongodb.M)
	if !ok {
		return nil, false
	}
	return map[string]any(fields), true
}

func orderedDocumentFields(fields libmongodb.D) map[string]any {
	result := make(map[string]any, len(fields))
	for _, field := range fields {
		result[field.Key] = field.Value
	}
	return result
}

func stringField(value any) string {
	field, ok := value.(string)
	if !ok {
		return ""
	}
	return field
}
