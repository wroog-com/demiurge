package cmd

import (
	"errors"
	"testing"

	"github.com/wroog-com/demiurge/internal/cmdutil"
	"github.com/wroog-com/demiurge/internal/iostreams"
)

func testFactory() *cmdutil.Factory {
	ios, _, _, _ := iostreams.Test()
	return &cmdutil.Factory{IOStreams: ios}
}

func TestNewRootCmd_metadata(t *testing.T) {
	root := NewRootCmd(testFactory())

	if root.Use != "demi" {
		t.Errorf("Use = %q, want %q", root.Use, "demi")
	}
	if root.Short == "" {
		t.Error("Short should not be empty")
	}
	if !root.SilenceErrors {
		t.Error("SilenceErrors should be true so demi.Main controls error output")
	}
	if !root.SilenceUsage {
		t.Error("SilenceUsage should be true")
	}
	if !root.CompletionOptions.DisableDefaultCmd {
		t.Error("default completion command should be disabled")
	}
}

func TestNewRootCmd_hasVersionSubcommand(t *testing.T) {
	root := NewRootCmd(testFactory())

	for _, c := range root.Commands() {
		if c.Name() == "version" {
			return
		}
	}
	t.Error("root command should register the version subcommand")
}

func TestNewRootCmd_unknownCommandErrors(t *testing.T) {
	root := NewRootCmd(testFactory())
	root.SetArgs([]string{"does-not-exist"})

	if err := root.Execute(); err == nil {
		t.Error("expected error for unknown command")
	}
}

func TestNewRootCmd_flagErrorsAreWrapped(t *testing.T) {
	root := NewRootCmd(testFactory())
	root.SetArgs([]string{"--nope"})

	err := root.Execute()
	if err == nil {
		t.Fatal("expected error for unknown flag")
	}
	var flagErr *cmdutil.FlagError
	if !errors.As(err, &flagErr) {
		t.Errorf("expected a *cmdutil.FlagError, got %T: %v", err, err)
	}
}
