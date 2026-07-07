package cmd

import (
	"testing"
)

func TestFormat(t *testing.T) {
	got := Format("1.2.3", "2026-07-06")
	want := "demi version 1.2.3 (2026-07-06)\n"
	if got != want {
		t.Errorf("Format() = %q, want %q", got, want)
	}
}

func TestFormat_noDate(t *testing.T) {
	got := Format("1.2.3", "")
	want := "demi version 1.2.3\n"
	if got != want {
		t.Errorf("Format() = %q, want %q", got, want)
	}
}
