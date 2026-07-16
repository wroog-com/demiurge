package iostreams

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/term"
)

type IOStreams struct {
	In  io.Reader
	Out io.Writer
	Err io.Writer

	inFile  *os.File
	outFile *os.File
	errFile *os.File

	stdinTTYOverride  *bool
	stdoutTTYOverride *bool
	stderrTTYOverride *bool
	colorOverride     *bool

	getenv func(string) string
}

func System() *IOStreams {
	return &IOStreams{
		In:      os.Stdin,
		Out:     os.Stdout,
		Err:     os.Stderr,
		inFile:  os.Stdin,
		outFile: os.Stdout,
		errFile: os.Stderr,
		getenv:  os.Getenv,
	}
}

func Test() (*IOStreams, *bytes.Buffer, *bytes.Buffer, *bytes.Buffer) {
	in := &bytes.Buffer{}
	out := &bytes.Buffer{}
	errBuf := &bytes.Buffer{}
	ios := &IOStreams{In: in, Out: out, Err: errBuf, getenv: func(string) string { return "" }}
	ios.SetStdinTTY(false)
	ios.SetStdoutTTY(false)
	ios.SetStderrTTY(false)
	return ios, in, out, errBuf
}

func (s *IOStreams) SetEnv(env map[string]string) {
	s.getenv = func(key string) string { return env[key] }
}

func (s *IOStreams) env(key string) string {
	if s.getenv == nil {
		return os.Getenv(key)
	}
	return s.getenv(key)
}

func (s *IOStreams) IsStdinTTY() bool {
	if s.stdinTTYOverride != nil {
		return *s.stdinTTYOverride
	}
	if s.inFile == nil {
		return false
	}
	return term.IsTerminal(int(s.inFile.Fd()))
}

func (s *IOStreams) IsStdoutTTY() bool {
	if s.stdoutTTYOverride != nil {
		return *s.stdoutTTYOverride
	}
	if s.outFile == nil {
		return false
	}
	return term.IsTerminal(int(s.outFile.Fd()))
}

func (s *IOStreams) IsStderrTTY() bool {
	if s.stderrTTYOverride != nil {
		return *s.stderrTTYOverride
	}
	if s.errFile == nil {
		return false
	}
	return term.IsTerminal(int(s.errFile.Fd()))
}

// CLICOLOR_FORCE beats NO_COLOR, CLICOLOR=0, and TERM=dumb.
func (s *IOStreams) ColorEnabled() bool {
	if s.colorOverride != nil {
		return *s.colorOverride
	}
	if s.colorForced() {
		return true
	}
	if s.colorDisabled() {
		return false
	}
	return s.IsStdoutTTY()
}

func (s *IOStreams) colorForced() bool {
	v := s.env("CLICOLOR_FORCE")
	return v != "" && v != "0"
}

func (s *IOStreams) colorDisabled() bool {
	return s.env("NO_COLOR") != "" || s.env("CLICOLOR") == "0" || s.env("TERM") == "dumb"
}

func (s *IOStreams) SetStdinTTY(v bool)  { s.stdinTTYOverride = &v }
func (s *IOStreams) SetStdoutTTY(v bool) { s.stdoutTTYOverride = &v }
func (s *IOStreams) SetStderrTTY(v bool) { s.stderrTTYOverride = &v }

func (s *IOStreams) ForceTerminal() {
	s.SetStdinTTY(true)
	s.SetStdoutTTY(true)
	s.SetStderrTTY(true)
}

func (s *IOStreams) ForceNoTerminal() {
	s.SetStdinTTY(false)
	s.SetStdoutTTY(false)
	s.SetStderrTTY(false)
}

func (s *IOStreams) ForceColorEnabled()  { v := true; s.colorOverride = &v }
func (s *IOStreams) ForceColorDisabled() { v := false; s.colorOverride = &v }

func (s *IOStreams) IsDebug() bool {
	switch strings.ToLower(s.env("DEMI_DEBUG")) {
	case "", "false", "0", "no":
		return false
	default:
		return true
	}
}

func (s *IOStreams) Debugf(format string, args ...any) {
	if s.IsDebug() {
		fmt.Fprintf(s.Err, "DEBUG: "+format+"\n", args...)
	}
}
