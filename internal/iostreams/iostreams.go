package iostreams

import (
	"bytes"
	"io"
	"os"

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
}

func System() *IOStreams {
	return &IOStreams{
		In:      os.Stdin,
		Out:     os.Stdout,
		Err:     os.Stderr,
		inFile:  os.Stdin,
		outFile: os.Stdout,
		errFile: os.Stderr,
	}
}

func Test() (*IOStreams, *bytes.Buffer, *bytes.Buffer, *bytes.Buffer) {
	in := &bytes.Buffer{}
	out := &bytes.Buffer{}
	errBuf := &bytes.Buffer{}
	ios := &IOStreams{In: in, Out: out, Err: errBuf}
	ios.SetStdinTTY(false)
	ios.SetStdoutTTY(false)
	ios.SetStderrTTY(false)
	return ios, in, out, errBuf
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

func (s *IOStreams) ColorEnabled() bool {
	if s.colorOverride != nil {
		return *s.colorOverride
	}
	if os.Getenv("NO_COLOR") != "" {
		return false
	}
	return s.IsStdoutTTY()
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
