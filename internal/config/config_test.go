package config

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

const testFileMode = 0600

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.StorageBucket != "" || cfg.ServerAPIGroup != "raw" {
		t.Fatalf("unexpected defaults: %+v", cfg)
	}
}

func TestLoadReadsEnvironmentAndDotEnv(t *testing.T) {
	resetConfigEnvironment(t)
	filename := writeTestDotEnv(t, "OHM9996_STORAGE_BUCKET=assets-bucket\n")
	cfg, err := Load(filename)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	assertLoadedConfig(t, cfg, "assets-bucket")
}

func resetConfigEnvironment(t *testing.T) {
	t.Helper()
	t.Setenv("OHM9996_SERVER_API_GROUP", "")
	t.Setenv("OHM9996_STORAGE_BUCKET", "")
	if err := os.Unsetenv("OHM9996_STORAGE_BUCKET"); err != nil {
		t.Fatalf("unset storage bucket environment: %v", err)
	}
}

func writeTestDotEnv(t *testing.T, content string) string {
	t.Helper()
	filename := filepath.Join(t.TempDir(), ".env")
	if err := os.WriteFile(filename, []byte(content), testFileMode); err != nil {
		t.Fatalf("write dotenv file: %v", err)
	}
	return filename
}

func assertLoadedConfig(t *testing.T, cfg *Config, bucket string) {
	t.Helper()
	if cfg.StorageBucket != bucket || cfg.ServerAPIGroup != "raw" {
		t.Fatalf("unexpected loaded config: %+v", cfg)
	}
}

func TestLoadRejectsNonRawAPIGroup(t *testing.T) {
	t.Setenv("OHM9996_SERVER_API_GROUP", "standard")
	filename := writeTestDotEnv(t, "")
	if _, err := Load(filename); !errors.Is(err, ErrUnsupportedServerAPIGroup) {
		t.Fatalf("expected ErrUnsupportedServerAPIGroup, got %v", err)
	}
}

func TestValidateRejectsNonRawAPIGroup(t *testing.T) {
	cfg := DefaultConfig()
	cfg.ServerAPIGroup = "standard"
	if !errors.Is(cfg.Validate(), ErrUnsupportedServerAPIGroup) {
		t.Fatal("expected non-raw API group to be rejected")
	}
}
