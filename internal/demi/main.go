package demi

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wroog-com/demiurge/internal/cmd"
	"github.com/wroog-com/demiurge/internal/cmdutil"
	"github.com/wroog-com/demiurge/internal/iostreams"
)

type ExitCode int

const (
	ExitOK     ExitCode = 0
	ExitError  ExitCode = 1
	ExitCancel ExitCode = 2
)

func Main() ExitCode {
	f := &cmdutil.Factory{
		IOStreams: iostreams.System(),
	}
	return run(f)
}

func run(f *cmdutil.Factory) ExitCode {
	// ExecuteC returns the command that ran, so printError can show its usage.
	executed, err := cmd.NewRootCmd(f).ExecuteC()
	return mapError(f, executed, err)
}

func mapError(f *cmdutil.Factory, cmd *cobra.Command, err error) ExitCode {
	if err == nil {
		return ExitOK
	}
	if errors.Is(err, cmdutil.SilentError) {
		return ExitError
	}
	if cmdutil.IsUserCancellation(err) {
		return ExitCancel
	}
	printError(f.IOStreams.Err, err, cmd)
	return ExitError
}

// printError writes err, appending usage for flag/argument and unknown-command
// errors.
func printError(out io.Writer, err error, cmd *cobra.Command) {
	fmt.Fprintln(out, err)

	var flagError *cmdutil.FlagError
	if errors.As(err, &flagError) || strings.HasPrefix(err.Error(), "unknown command ") {
		if cmd != nil {
			if !strings.HasSuffix(err.Error(), "\n") {
				fmt.Fprintln(out)
			}
			fmt.Fprint(out, cmd.UsageString())
		}
	}
}
