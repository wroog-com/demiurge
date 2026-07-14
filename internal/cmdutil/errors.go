package cmdutil

import (
	"context"
	"errors"
	"fmt"
)

// ErrSilent signals exit code 1 without printing anything. Use it when the
// subcommand has already printed its own error message.
var ErrSilent = errors.New("ErrSilent")

// ErrCancel signals user-initiated cancellation; it maps to its own exit code.
var ErrCancel = errors.New("ErrCancel")

// IsUserCancellation reports whether err represents the user aborting the
// command. Deadline expiry is deliberately excluded: a timeout is a failure,
// not a cancellation.
func IsUserCancellation(err error) bool {
	return errors.Is(err, ErrCancel) ||
		errors.Is(err, context.Canceled)
}

// A FlagError indicates a problem parsing flags or arguments; such errors cause
// the command's usage to be shown.
type FlagError struct {
	err error
}

func (e *FlagError) Error() string { return e.err.Error() }
func (e *FlagError) Unwrap() error { return e.err }

// FlagErrorf returns a FlagError wrapping fmt.Errorf(format, args...).
func FlagErrorf(format string, args ...any) error {
	return FlagErrorWrap(fmt.Errorf(format, args...))
}

// FlagErrorWrap returns a FlagError wrapping err.
func FlagErrorWrap(err error) error { return &FlagError{err} }
