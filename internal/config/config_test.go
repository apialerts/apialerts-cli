package config

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func setupTestDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	configDirOverride = dir
	t.Cleanup(func() { configDirOverride = "" })
	return dir
}

func TestLoadNoFile(t *testing.T) {
	setupTestDir(t)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if cfg.APIKey != "" {
		t.Fatalf("expected empty API key, got %q", cfg.APIKey)
	}
}

func TestSaveAndLoad(t *testing.T) {
	setupTestDir(t)

	err := Save(&CLIConfig{APIKey: "test-key-123"})
	if err != nil {
		t.Fatalf("expected no error saving, got %v", err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("expected no error loading, got %v", err)
	}
	if cfg.APIKey != "test-key-123" {
		t.Fatalf("expected API key %q, got %q", "test-key-123", cfg.APIKey)
	}
}

func TestSaveFilePermissions(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("file permissions not supported on Windows")
	}
	dir := setupTestDir(t)

	err := Save(&CLIConfig{APIKey: "test-key"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	info, err := os.Stat(filepath.Join(dir, configFile))
	if err != nil {
		t.Fatalf("expected config file to exist, got %v", err)
	}

	perm := info.Mode().Perm()
	if perm != 0600 {
		t.Fatalf("expected file permissions 0600, got %04o", perm)
	}
}

func TestLoadCorruptedFile(t *testing.T) {
	dir := setupTestDir(t)

	err := os.WriteFile(filepath.Join(dir, configFile), []byte("not json"), 0600)
	if err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	_, err = Load()
	if err == nil {
		t.Fatal("expected error for corrupted config, got nil")
	}
}

func TestGetAPIKeyNotConfigured(t *testing.T) {
	setupTestDir(t)

	_, err := GetAPIKey()
	if err == nil {
		t.Fatal("expected error when no key configured, got nil")
	}
}

func TestGetAPIKeyConfigured(t *testing.T) {
	setupTestDir(t)

	err := Save(&CLIConfig{APIKey: "my-key"})
	if err != nil {
		t.Fatalf("expected no error saving, got %v", err)
	}

	key, err := GetAPIKey()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if key != "my-key" {
		t.Fatalf("expected key %q, got %q", "my-key", key)
	}
}

func TestSaveOverwrite(t *testing.T) {
	setupTestDir(t)

	Save(&CLIConfig{APIKey: "old-key"})
	Save(&CLIConfig{APIKey: "new-key"})

	cfg, err := Load()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if cfg.APIKey != "new-key" {
		t.Fatalf("expected key %q, got %q", "new-key", cfg.APIKey)
	}
}
