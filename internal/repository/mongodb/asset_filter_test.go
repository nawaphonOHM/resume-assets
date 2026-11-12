package mongodb

import (
	"encoding/json"
	"testing"

	"github.com/nawaphonOHM/resume-assets/internal/domain"
	libmongodb "github.com/nawaphonOHM/whatever/v3/pkg/testing/mongodb"
)

func TestBuildFilterMatchesExactFilename(t *testing.T) {
	const filename = "bundle/" + testFilename
	filter := buildTestFilter(t, filename, nil)
	assertSingleFilenameFilter(t, filter, filename)
}

func buildTestFilter(t *testing.T, filename string, criteria *domain.AssetCriteria) libmongodb.D {
	t.Helper()
	filter, err := BuildFilter(filename, criteria)
	if err != nil {
		t.Fatalf("build filter: %v", err)
	}
	return filter
}

func assertSingleFilenameFilter(t *testing.T, filter libmongodb.D, filename string) {
	t.Helper()
	if len(filter) != 1 {
		t.Fatalf("expected one filter field, got %d", len(filter))
	}
	assertFilenameElement(t, filter[0], filename)
}

func assertFilenameElement(t *testing.T, element libmongodb.E, filename string) {
	t.Helper()
	if element.Key != "filename" || element.Value != filename {
		t.Fatalf("expected exact filename filter, got %#v", element)
	}
}

func TestBuildFilterUsesLiteralCriteriaAndRecursiveExpressions(t *testing.T) {
	criteria := domain.AssetCriteria{
		"$ne": "$where",
		"a.b": []any{json.Number("9007199254740993"), nil},
	}
	filter := buildTestFilter(t, testFilename, &criteria)
	assertRecursiveFilterExpressions(t, filter)
}

func assertRecursiveFilterExpressions(t *testing.T, filter libmongodb.D) {
	t.Helper()
	assertFilterHasExpression(t, filter)
	assertFilterEntries(t, filter)
	assertNoCriteriaOperators(t, filter)
}

func assertFilterHasExpression(t *testing.T, filter libmongodb.D) {
	t.Helper()
	if len(filter) != 2 {
		t.Fatalf("expected filename plus criteria expression, got %#v", filter)
	}
	assertExpressionElement(t, filter[1])
}

func assertExpressionElement(t *testing.T, element libmongodb.E) {
	t.Helper()
	if element.Key != "$expr" {
		t.Fatalf("expected criteria expression, got %#v", element)
	}
	if _, ok := element.Value.(libmongodb.D); !ok {
		t.Fatalf("expected BSON expression document, got %T", element.Value)
	}
}

func assertNoCriteriaOperators(t *testing.T, filter libmongodb.D) {
	t.Helper()
	if containsFilterEntry(filter, "$ne", nil) || containsFilterEntry(filter, "$where", nil) {
		t.Fatal("criteria keys or values were interpolated as operators")
	}
}

func assertFilterEntries(t *testing.T, filter any) {
	t.Helper()
	entries := []struct {
		value any
		key   string
	}{
		{key: "$objectToArray"},
		{key: "$filter"},
		{key: "$arrayElemAt"},
		{key: "$literal", value: "$ne"},
		{key: "$literal", value: "$where"},
		{key: "$literal", value: "a.b"},
		{key: "$toDecimal"},
	}
	for _, entry := range entries {
		if !containsFilterEntry(filter, entry.key, entry.value) {
			t.Errorf("filter does not contain %q with value %#v", entry.key, entry.value)
		}
	}
}

func TestBuildFilterRejectsUnsupportedCriteriaValues(t *testing.T) {
	criteria := domain.AssetCriteria{"unsupported": make(chan int)}
	if _, err := BuildFilter(testFilename, &criteria); err == nil {
		t.Fatal("expected unsupported criteria value to fail filter construction")
	}
}
