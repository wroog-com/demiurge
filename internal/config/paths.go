package config

import (
	"fmt"
	"os"
	"path/filepath"
)

// Priority: DEMI_CONFIG_DIR > XDG_CONFIG_HOME > ~/.config/demi.
func ConfigDir() (string, error) {
	if dir := os.Getenv("DEMI_CONFIG_DIR"); dir != "" {
		return dir, nil
	}
	if dir := os.Getenv("XDG_CONFIG_HOME"); dir != "" {
		return filepath.Join(dir, "demi"), nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("determine home directory for config: %w", err)
	}
	return filepath.Join(home, ".config", "demi"), nil
}

// Priority: DEMI_STATE_DIR > XDG_STATE_HOME > ~/.local/state/demi.
func StateDir() (string, error) {
	if dir := os.Getenv("DEMI_STATE_DIR"); dir != "" {
		return dir, nil
	}
	if dir := os.Getenv("XDG_STATE_HOME"); dir != "" {
		return filepath.Join(dir, "demi"), nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("determine home directory for state: %w", err)
	}
	return filepath.Join(home, ".local", "state", "demi"), nil
}
