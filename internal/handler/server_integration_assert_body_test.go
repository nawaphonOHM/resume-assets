package handler

import (
	"reflect"
	"testing"
	"time"
)

func assertIntegrationBodyAndLookup(t *testing.T, service *routerAssetService, expectedFilenames ...string) {
	t.Helper()
	body, filenames := service.snapshot()
	if body == nil {
		t.Fatal("framework stream did not produce an object body")
	}
	assertIntegrationBodyClosed(t, body)
	assertIntegrationLookups(t, filenames, expectedFilenames)
}

func assertIntegrationLookups(t *testing.T, filenames, expected []string) {
	t.Helper()
	if !reflect.DeepEqual(filenames, expected) {
		t.Fatalf("unexpected service lookups: %#v, want %#v", filenames, expected)
	}
}

func assertIntegrationBodyClosed(t *testing.T, body *integrationBody) {
	t.Helper()
	select {
	case <-body.closed:
	case <-time.After(integrationWaitTimeout):
		t.Fatal("framework stream did not close the object body")
	}
}
