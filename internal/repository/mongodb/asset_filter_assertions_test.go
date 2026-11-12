package mongodb

import (
	"reflect"

	libmongodb "github.com/nawaphonOHM/whatever/v3/pkg/testing/mongodb"
)

func containsFilterEntry(value any, key string, want any) bool {
	if document, ok := value.(libmongodb.D); ok {
		return containsDocumentEntry(document, key, want)
	}
	array, ok := value.(libmongodb.A)
	return ok && containsArrayEntry(array, key, want)
}

func containsDocumentEntry(document libmongodb.D, key string, want any) bool {
	for _, element := range document {
		if containsDocumentElement(element, key, want) {
			return true
		}
	}
	return false
}

func containsDocumentElement(element libmongodb.E, key string, want any) bool {
	if filterEntryMatches(element.Key, element.Value, key, want) {
		return true
	}
	return containsFilterEntry(element.Value, key, want)
}

func filterEntryMatches(gotKey string, gotValue any, wantKey string, wantValue any) bool {
	return gotKey == wantKey && (wantValue == nil || reflect.DeepEqual(gotValue, wantValue))
}

func containsArrayEntry(array libmongodb.A, key string, want any) bool {
	for _, item := range array {
		if containsFilterEntry(item, key, want) {
			return true
		}
	}
	return false
}
