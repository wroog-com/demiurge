package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestConfigDir_default(t *testing.T) {
	t.Setenv("DEMI_CONFIG_DIR", "")
	t.Setenv("XDG_CONFIG_HOME", "")
	home, _ := os.UserHomeDir()
	want := filepath.Join(home, ".config", "demi")
	got, err := ConfigDir()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != want {
		t.Errorf("ConfigDir() = %q, want %q", got, want)
	}
}

func TestConfigDir_xdg(t *testing.T) {
	t.Setenv("DEMI_CONFIG_DIR", "")
	t.Setenv("XDG_CONFIG_HOME", "/custom/config")
	got, err := ConfigDir()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "/custom/config/demi" {
		t.Errorf("ConfigDir() = %q, want %q", got, "/custom/config/demi")
	}
}

func TestConfigDir_appOverride(t *testing.T) {
	t.Setenv("DEMI_CONFIG_DIR", "/override/config")
	t.Setenv("XDG_CONFIG_HOME", "/custom/config")
	got, err := ConfigDir()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "/override/config" {
		t.Errorf("ConfigDir() = %q, want %q", got, "/override/config")
	}
}

func TestConfigDir_containsDemi(t *testing.T) {
	t.Setenv("DEMI_CONFIG_DIR", "")
	got, err := ConfigDir()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(got, "demi") {
		t.Errorf("ConfigDir() = %q, expected it to contain 'demi'", got)
	}
}

func TestStateDir_default(t *testing.T) {
	t.Setenv("DEMI_STATE_DIR", "")
	t.Setenv("XDG_STATE_HOME", "")
	home, _ := os.UserHomeDir()
	want := filepath.Join(home, ".local", "state", "demi")
	got, err := StateDir()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != want {
		t.Errorf("StateDir() = %q, want %q", got, want)
	}
}

func TestStateDir_xdg(t *testing.T) {
	t.Setenv("DEMI_STATE_DIR", "")
	t.Setenv("XDG_STATE_HOME", "/custom/state")
	got, err := StateDir()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "/custom/state/demi" {
		t.Errorf("StateDir() = %q, want %q", got, "/custom/state/demi")
	}
}

func TestStateDir_appOverride(t *testing.T) {
	t.Setenv("DEMI_STATE_DIR", "/override/state")
	t.Setenv("XDG_STATE_HOME", "/custom/state")
	got, err := StateDir()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "/override/state" {
		t.Errorf("StateDir() = %q, want %q", got, "/override/state")
	}
}

func TestStateDir_containsDemi(t *testing.T) {
	t.Setenv("DEMI_STATE_DIR", "")
	got, err := StateDir()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(got, "demi") {
		t.Errorf("StateDir() = %q, expected it to contain 'demi'", got)
	}
}
