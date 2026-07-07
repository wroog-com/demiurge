package cmdutil

import (
	"errors"
	"fmt"
	"testing"
)

func TestSilentError(t *testing.T) {
	err := SilentError
	if err == nil {
		t.Fatal("SilentError should not be nil")
	}
	if !errors.Is(err, SilentError) {
		t.Error("errors.Is should match SilentError")
	}
	if err.Error() == "" {
		t.Error("SilentError should have a non-empty message")
	}
}

func TestSilentError_wrapping(t *testing.T) {
	wrapped := fmt.Errorf("context: %w", SilentError)
	if !errors.Is(wrapped, SilentError) {
		t.Error("wrapped SilentError should be detectable via errors.Is")
	}
}
