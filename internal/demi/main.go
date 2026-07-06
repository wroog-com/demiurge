package demi

import (
	"errors"
	"fmt"

	"github.com/wroog-com/demiurge/internal/cmd"
	"github.com/wroog-com/demiurge/internal/cmdutil"
	"github.com/wroog-com/demiurge/internal/iostreams"
)

type ExitCode int

const (
	ExitOK    ExitCode = 0
	ExitError ExitCode = 1
)

func Main() ExitCode {
	ioStreams := iostreams.System()

	if err := cmd.NewRootCmd(ioStreams).Execute(); err != nil {
		if !errors.Is(err, cmdutil.SilentError) {
			fmt.Fprintln(ioStreams.Err, err)
		}
		return ExitError
	}
	return ExitOK
}
