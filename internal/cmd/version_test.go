package cmd

import "testing"

func TestVersionString(t *testing.T) {
	tests := []struct {
		name, version, date, want string
	}{
		{"release stamp", "1.2.3", "2026-07-06", "1.2.3 (2026-07-06)"},
		{"no date", "1.2.3", "", "1.2.3"},
		{"module version normalized", "v0.3.1", "", "0.3.1"},
		{"pseudo-version normalized", "v0.3.2-0.20260712122422-f28e83e31b39", "", "0.3.2-0.20260712122422-f28e83e31b39"},
		{"dev fallback", "dev", "", "dev"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := versionString(tt.version, tt.date); got != tt.want {
				t.Errorf("versionString(%q, %q) = %q, want %q", tt.version, tt.date, got, tt.want)
			}
		})
	}
}
