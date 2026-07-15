package demi

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"

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
	ctx, stop := signalContext(context.Background())
	defer stop()

	f := &cmdutil.Factory{
		IOStreams: iostreams.System(),
	}
	return run(ctx, f, os.Args[1:])
}

// The first SIGINT/SIGTERM cancels the context; restoring the default
// disposition then lets a second signal kill a command that ignores it.
func signalContext(parent context.Context) (context.Context, context.CancelFunc) {
	ctx, stop := signal.NotifyContext(parent, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-ctx.Done()
		stop()
	}()
	return ctx, stop
}

func run(ctx context.Context, f *cmdutil.Factory, args []string) ExitCode {
	// A nil slice would make the framework fall back to ambient os.Args.
	if args == nil {
		args = []string{}
	}

	root := cmd.NewRootCmd(f)
	f.IOStreams.Debugf("demi %s args %q", root.Version, args)
	root.SetArgs(args)

	// ExecuteContextC returns the command that ran, so printError can show its usage.
	executed, err := root.ExecuteContextC(ctx)
	return mapError(f, executed, err)
}

func mapError(f *cmdutil.Factory, cmd *cobra.Command, err error) ExitCode {
	if err == nil {
		return ExitOK
	}
	if errors.Is(err, cmdutil.ErrSilent) {
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
