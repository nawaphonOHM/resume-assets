package mongodb

import "github.com/nawaphonOHM/resume-assets/internal/domain"

type assetRecord struct {
	Location any    `bson:"location"`
	Access   any    `bson:"access"`
	Filename string `bson:"filename"`
}

func (r assetRecord) toAssetDocument() *domain.AssetDocument {
	return &domain.AssetDocument{
		Filename: r.Filename,
		Location: stringField(r.Location),
		Access:   stringField(r.Access),
	}
}

func stringField(value any) string {
	field, ok := value.(string)
	if !ok {
		return ""
	}
	return field
}
