package cmdutil

import (
	"context"
	"errors"
	"fmt"
	"testing"
)

func TestErrSilent(t *testing.T) {
	err := ErrSilent
	if err == nil {
		t.Fatal("ErrSilent should not be nil")
	}
	if !errors.Is(err, ErrSilent) {
		t.Error("errors.Is should match ErrSilent")
	}
	if err.Error() == "" {
		t.Error("ErrSilent should have a non-empty message")
	}
}

func TestErrSilent_wrapping(t *testing.T) {
	wrapped := fmt.Errorf("context: %w", ErrSilent)
	if !errors.Is(wrapped, ErrSilent) {
		t.Error("wrapped ErrSilent should be detectable via errors.Is")
	}
}

func TestIsUserCancellation(t *testing.T) {
	cases := []struct {
		name string
		err  error
		want bool
	}{
		{"ErrCancel", ErrCancel, true},
		{"wrapped ErrCancel", fmt.Errorf("aborting: %w", ErrCancel), true},
		{"context.Canceled", context.Canceled, true},
		{"context.DeadlineExceeded", context.DeadlineExceeded, true},
		{"unrelated error", errors.New("boom"), false},
		{"ErrSilent", ErrSilent, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := IsUserCancellation(tc.err); got != tc.want {
				t.Errorf("IsUserCancellation(%v) = %v, want %v", tc.err, got, tc.want)
			}
		})
	}
}

func TestFlagError(t *testing.T) {
	inner := errors.New("bad flag")
	fe := FlagErrorWrap(inner)

	if fe.Error() != "bad flag" {
		t.Errorf("Error() = %q, want %q", fe.Error(), "bad flag")
	}
	if !errors.Is(fe, inner) {
		t.Error("FlagError should unwrap to the wrapped error")
	}

	var target *FlagError
	if !errors.As(fe, &target) {
		t.Error("errors.As should recognise a *FlagError")
	}
}

func TestFlagErrorf(t *testing.T) {
	fe := FlagErrorf("unknown flag: %s", "--nope")
	if fe.Error() != "unknown flag: --nope" {
		t.Errorf("Error() = %q, want %q", fe.Error(), "unknown flag: --nope")
	}

	var target *FlagError
	if !errors.As(fe, &target) {
		t.Error("FlagErrorf should produce a *FlagError")
	}
}
