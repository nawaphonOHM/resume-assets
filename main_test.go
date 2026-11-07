package main

import (
	"context"
	"errors"
	"testing"

	"github.com/nawaphonOHM/resume-assets/internal/config"
	"github.com/nawaphonOHM/whatever/v3/pkg/mongodb"
	"github.com/nawaphonOHM/whatever/v3/pkg/rest"
	"github.com/nawaphonOHM/whatever/v3/pkg/storage"
)

func TestDisconnectClientNil(t *testing.T) {
	t.Helper()
	disconnectClient(context.Background(), nil)
}

func TestBuildBluePrintSuccess(t *testing.T) {
	blueprint, err := buildBluePrint(config.DefaultConfig(), &mongodb.Client{}, &storage.Client{})
	if err != nil {
		t.Fatalf("buildBluePrint: %v", err)
	}
	if blueprint == nil {
		t.Fatal("expected non-nil blueprint")
	}
}

func TestBuildBluePrintRequiresMongoClient(t *testing.T) {
	blueprint, err := buildBluePrint(config.DefaultConfig(), nil, &storage.Client{})
	assertBuildError(t, blueprint, err, errNilMongoClient)
}

func TestBuildBluePrintRequiresStorageClient(t *testing.T) {
	blueprint, err := buildBluePrint(config.DefaultConfig(), &mongodb.Client{}, nil)
	assertBuildError(t, blueprint, err, errNilStorageClient)
}

func assertBuildError(t *testing.T, blueprint *rest.BluePrint, got, want error) {
	t.Helper()
	if !errors.Is(got, want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
	if blueprint != nil {
		t.Fatal("expected no blueprint on invalid clients")
	}
}

func TestRunRejectsNonRawAPIRouting(t *testing.T) {
	t.Setenv("OHM9996_SERVER_API_GROUP", "standard")
	if err := run(context.Background()); !errors.Is(err, config.ErrUnsupportedServerAPIGroup) {
		t.Fatalf("expected raw routing config error, got %v", err)
	}
}
