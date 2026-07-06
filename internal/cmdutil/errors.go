package cmdutil

import "errors"

// SilentError signals exit code 1 without printing anything. Use it when the
// subcommand has already printed its own error message.
var SilentError = errors.New("SilentError")
