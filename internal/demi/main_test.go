package demi

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/wroog-com/demiurge/internal/cmdutil"
	"github.com/wroog-com/demiurge/internal/iostreams"
)

func TestRun_success(t *testing.T) {
	ios, _, _, errBuf := iostreams.Test()
	f := &cmdutil.Factory{IOStreams: ios}

	if code := run(f); code != ExitOK {
		t.Errorf("run() = %d, want ExitOK", code)
	}
	if errBuf.String() != "" {
		t.Errorf("expected nothing on Err for success, got %q", errBuf.String())
	}
}

func TestMapError_nil(t *testing.T) {
	ios, _, _, errBuf := iostreams.Test()
	f := &cmdutil.Factory{IOStreams: ios}

	if code := mapError(f, nil, nil); code != ExitOK {
		t.Errorf("mapError(nil) = %d, want ExitOK", code)
	}
	if errBuf.String() != "" {
		t.Errorf("expected no output for nil error, got %q", errBuf.String())
	}
}

func TestMapError_printsError(t *testing.T) {
	ios, _, _, errBuf := iostreams.Test()
	f := &cmdutil.Factory{IOStreams: ios}

	code := mapError(f, nil, errors.New("something failed"))
	if code != ExitError {
		t.Errorf("mapError() = %d, want ExitError", code)
	}
	if got := errBuf.String(); got != "something failed\n" {
		t.Errorf("Err = %q, want %q", got, "something failed\n")
	}
}

func TestMapError_silentErrorNotPrinted(t *testing.T) {
	ios, _, _, errBuf := iostreams.Test()
	f := &cmdutil.Factory{IOStreams: ios}

	code := mapError(f, nil, cmdutil.SilentError)
	if code != ExitError {
		t.Errorf("mapError(SilentError) = %d, want ExitError", code)
	}
	if errBuf.String() != "" {
		t.Errorf("SilentError should be suppressed, got %q", errBuf.String())
	}
}

func TestMapError_wrappedSilentErrorNotPrinted(t *testing.T) {
	ios, _, _, errBuf := iostreams.Test()
	f := &cmdutil.Factory{IOStreams: ios}

	wrapped := fmt.Errorf("context: %w", cmdutil.SilentError)
	if code := mapError(f, nil, wrapped); code != ExitError {
		t.Errorf("mapError(wrapped SilentError) = %d, want ExitError", code)
	}
	if errBuf.String() != "" {
		t.Errorf("wrapped SilentError should be suppressed, got %q", errBuf.String())
	}
}

func TestMapError_cancellation(t *testing.T) {
	cases := map[string]error{
		"CancelError":      cmdutil.CancelError,
		"context.Canceled": context.Canceled,
		"wrapped":          fmt.Errorf("aborting: %w", cmdutil.CancelError),
	}
	for name, err := range cases {
		t.Run(name, func(t *testing.T) {
			ios, _, _, errBuf := iostreams.Test()
			f := &cmdutil.Factory{IOStreams: ios}

			if code := mapError(f, nil, err); code != ExitCancel {
				t.Errorf("mapError(%s) = %d, want ExitCancel", name, code)
			}
			if errBuf.String() != "" {
				t.Errorf("cancellation should be silent, got %q", errBuf.String())
			}
		})
	}
}

func TestPrintError_appendsUsageOnFlagError(t *testing.T) {
	ios, _, _, errBuf := iostreams.Test()
	c := &cobra.Command{Use: "demi"}

	err := cmdutil.FlagErrorf("unknown flag: --nope")
	printError(ios.Err, err, c)

	out := errBuf.String()
	if !strings.Contains(out, "unknown flag: --nope") {
		t.Errorf("expected error message in output, got %q", out)
	}
	if !strings.Contains(out, c.UsageString()) {
		t.Errorf("expected usage string appended for FlagError, got %q", out)
	}
}

func TestPrintError_appendsUsageOnUnknownCommand(t *testing.T) {
	ios, _, _, errBuf := iostreams.Test()
	c := &cobra.Command{Use: "demi"}

	err := errors.New("unknown command \"nope\" for \"demi\"")
	printError(ios.Err, err, c)

	if !strings.Contains(errBuf.String(), c.UsageString()) {
		t.Errorf("expected usage appended for unknown command, got %q", errBuf.String())
	}
}

func TestPrintError_plainErrorNoUsage(t *testing.T) {
	ios, _, _, errBuf := iostreams.Test()
	c := &cobra.Command{Use: "demi"}

	printError(ios.Err, errors.New("disk on fire"), c)

	out := errBuf.String()
	if out != "disk on fire\n" {
		t.Errorf("expected only the error line, got %q", out)
	}
}
