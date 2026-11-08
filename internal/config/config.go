// Package config provides application configuration loading and validation.
package config

import (
	"errors"
	"os"
	"strings"
)

const (
	// DefaultStorageBucket is empty until configured by the shared storage client.
	DefaultStorageBucket = ""
	// DefaultServerAPIGroup keeps application routes on their raw paths.
	DefaultServerAPIGroup = "raw"
	serverAPIGroupEnv     = "OHM9996_SERVER_API_GROUP"
)

// ErrUnsupportedServerAPIGroup indicates that a non-raw route group was configured.
var ErrUnsupportedServerAPIGroup = errors.New("server API group must be raw")

// Config contains application settings not managed by the shared clients.
type Config struct {
	StorageBucket  string
	ServerAPIGroup string
}

// DefaultConfig returns configuration defaults for the asset server.
func DefaultConfig() *Config {
	return &Config{
		StorageBucket:  DefaultStorageBucket,
		ServerAPIGroup: DefaultServerAPIGroup,
	}
}

// Load reads optional dotenv files and OHM9996_* environment settings.
func Load(filenames ...string) (*Config, error) {
	if err := initializeEnvironment(filenames...); err != nil {
		return nil, err
	}
	return loadConfigFromEnvironment()
}

func initializeEnvironment(filenames ...string) error {
	if err := loadDotEnvFiles(filenames...); err != nil {
		return err
	}
	return ensureDefaultServerAPIGroup()
}

func loadConfigFromEnvironment() (*Config, error) {
	cfg := DefaultConfig()
	applyEnvOverrides(cfg)
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}

// Validate ensures the application keeps raw routing so the asset path is stable.
func (c *Config) Validate() error {
	if c.ServerAPIGroup != DefaultServerAPIGroup {
		return ErrUnsupportedServerAPIGroup
	}
	return nil
}

func applyEnvOverrides(cfg *Config) {
	applyEnvString(&cfg.StorageBucket, "OHM9996_STORAGE_BUCKET")
	applyEnvString(&cfg.ServerAPIGroup, serverAPIGroupEnv)
}

func applyEnvString(target *string, key string) {
	if value := os.Getenv(key); value != "" {
		*target = value
	}
}

func ensureDefaultServerAPIGroup() error {
	value, exists := os.LookupEnv(serverAPIGroupEnv)
	if exists && strings.TrimSpace(value) != "" {
		return nil
	}
	return os.Setenv(serverAPIGroupEnv, DefaultServerAPIGroup)
}
