package mongodb

import (
	"encoding/json"
	"reflect"
	"testing"

	libmongodb "github.com/nawaphonOHM/whatever/v3/pkg/testing/mongodb"
)

func TestCriteriaNumberExpressionEvaluatesDecimalConversion(t *testing.T) {
	const number = "9007199254740993"
	got, err := criteriaScalarExpression("$$assetValue", json.Number(number))
	if err != nil {
		t.Fatalf("build numeric expression: %v", err)
	}
	want := libmongodb.D{{Key: "$eq", Value: libmongodb.A{
		"$$assetValue",
		libmongodb.D{{Key: "$toDecimal", Value: libmongodb.D{{Key: "$literal", Value: number}}}},
	}}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("numeric expression = %#v, want %#v", got, want)
	}
}
