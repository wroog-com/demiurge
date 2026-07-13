package cmd

import (
	"errors"
	"strings"
	"testing"

	"github.com/wroog-com/demiurge/internal/cmdutil"
	"github.com/wroog-com/demiurge/internal/iostreams"
)

func runRoot(t *testing.T, args ...string) (stdout, stderr string, err error) {
	t.Helper()
	ios, _, out, errBuf := iostreams.Test()
	root := NewRootCmd(&cmdutil.Factory{IOStreams: ios})
	root.SetArgs(args)
	err = root.Execute()
	return out.String(), errBuf.String(), err
}

func TestVersionCmd_metadata(t *testing.T) {
	c := NewVersionCmd(testFactory())
	if c.Use != "version" {
		t.Errorf("Use = %q, want %q", c.Use, "version")
	}
	if c.Hidden {
		t.Error("version command must be visible")
	}
}

func TestVersionCmd_writesToOut(t *testing.T) {
	stdout, stderr, err := runRoot(t, "version")
	if err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}
	if !strings.HasPrefix(stdout, "demi version ") {
		t.Errorf("output = %q, want prefix %q", stdout, "demi version ")
	}
	if !strings.HasSuffix(stdout, "\n") {
		t.Errorf("output = %q, want trailing newline", stdout)
	}
	if stderr != "" {
		t.Errorf("version must not write to Err, got %q", stderr)
	}
}

func TestVersionFlag_matchesVersionCmd(t *testing.T) {
	flagOut, flagErr, err := runRoot(t, "--version")
	if err != nil {
		t.Fatalf("--version returned error: %v", err)
	}
	cmdOut, _, err := runRoot(t, "version")
	if err != nil {
		t.Fatalf("version returned error: %v", err)
	}
	if flagOut == "" || flagOut != cmdOut {
		t.Errorf("--version = %q, version = %q; the two surfaces must be byte-identical", flagOut, cmdOut)
	}
	if flagErr != "" {
		t.Errorf("--version must write to Out only, got Err %q", flagErr)
	}
}

func TestVersionFlag_hasNoShorthand(t *testing.T) {
	_, _, err := runRoot(t, "-v")
	if err == nil {
		t.Fatal("-v must stay unclaimed so it remains available later; --version has no shorthand")
	}
	var flagErr *cmdutil.FlagError
	if !errors.As(err, &flagErr) {
		t.Fatalf("-v error = %T (%v), want *cmdutil.FlagError so usage prints alongside it", err, err)
	}
	if !strings.Contains(err.Error(), "unknown shorthand flag: 'v'") {
		t.Errorf("-v error = %q, want unknown-shorthand for 'v' specifically", err.Error())
	}
}
