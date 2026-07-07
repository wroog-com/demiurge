package cmd

import (
	"strings"
	"testing"

	"github.com/wroog-com/demiurge/internal/cmdutil"
	"github.com/wroog-com/demiurge/internal/iostreams"
)

func TestVersionCmd_metadata(t *testing.T) {
	c := NewVersionCmd(testFactory())

	if c.Use != "version" {
		t.Errorf("Use = %q, want %q", c.Use, "version")
	}
	if !c.Hidden {
		t.Error("version command should be hidden")
	}
}

func TestVersionCmd_writesToOut(t *testing.T) {
	ios, _, out, errBuf := iostreams.Test()
	f := &cmdutil.Factory{IOStreams: ios}

	c := NewVersionCmd(f)
	if err := c.RunE(c, nil); err != nil {
		t.Fatalf("RunE returned error: %v", err)
	}

	got := out.String()
	if !strings.HasPrefix(got, "demi version ") {
		t.Errorf("output = %q, want it to start with %q", got, "demi version ")
	}
	if !strings.HasSuffix(got, "\n") {
		t.Errorf("output = %q, want it to end with a newline", got)
	}
	if errBuf.String() != "" {
		t.Errorf("version should not write to Err, got %q", errBuf.String())
	}
}
