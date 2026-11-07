// Package main wires the asset HTTP server and managed backend clients.
package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/nawaphonOHM/resume-assets/internal/config"
	"github.com/nawaphonOHM/resume-assets/internal/handler"
	mongorepo "github.com/nawaphonOHM/resume-assets/internal/repository/mongodb"
	"github.com/nawaphonOHM/resume-assets/internal/service"
	"github.com/nawaphonOHM/whatever/v3/pkg/logging"
	"github.com/nawaphonOHM/whatever/v3/pkg/mongodb"
	"github.com/nawaphonOHM/whatever/v3/pkg/rest"
	"github.com/nawaphonOHM/whatever/v3/pkg/storage"
)

const assetsCollection = "assets"

var (
	errNilMongoClient   = errors.New("mongodb client is nil")
	errNilStorageClient = errors.New("storage client is nil")
)

func validateClients(mongoClient *mongodb.Client, storageClient *storage.Client) error {
	if mongoClient == nil {
		return errNilMongoClient
	}
	if storageClient == nil {
		return errNilStorageClient
	}
	return nil
}

func buildBluePrint(
	cfg *config.Config,
	mongoClient *mongodb.Client,
	storageClient *storage.Client,
) (*rest.BluePrint, error) {
	if err := validateBlueprintInputs(cfg, mongoClient, storageClient); err != nil {
		return nil, err
	}

	repository := mongorepo.NewMongoAssetRepository(mongoClient.Collection(assetsCollection))
	assetService, err := service.NewAssetService(repository, storageClient, cfg.StorageBucket)
	if err != nil {
		return nil, fmt.Errorf("initialize asset service: %w", err)
	}
	return handler.RegisterRoutes(handler.NewAssetHandler(assetService)), nil
}

func validateBlueprintInputs(
	cfg *config.Config,
	mongoClient *mongodb.Client,
	storageClient *storage.Client,
) error {
	if err := validateClients(mongoClient, storageClient); err != nil {
		return err
	}
	if cfg == nil {
		return errors.New("config is nil")
	}
	return cfg.Validate()
}

func startServer(ctx context.Context, cfg *config.Config, mongoClient *mongodb.Client) error {
	storageClient, err := storage.Connect(ctx)
	if err != nil {
		return fmt.Errorf("connect storage: %w", err)
	}
	blueprint, err := buildBluePrint(cfg, mongoClient, storageClient)
	if err != nil {
		return err
	}
	return rest.StartREST(blueprint)
}

func disconnectClient(ctx context.Context, client *mongodb.Client) {
	if client == nil {
		return
	}
	if err := client.Disconnect(ctx); err != nil {
		logging.WarnContext(ctx, "failed to disconnect mongodb client", "error", err)
	}
}

func run(ctx context.Context) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}
	mongoClient, err := mongodb.Connect(ctx)
	if err != nil {
		return fmt.Errorf("connect mongodb: %w", err)
	}
	defer disconnectClient(ctx, mongoClient)
	return startServer(ctx, cfg, mongoClient)
}

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		logging.ErrorContext(ctx, "server stopped with error", "error", err)
		os.Exit(1)
	}
}
