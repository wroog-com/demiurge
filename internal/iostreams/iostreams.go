package iostreams

import (
	"bytes"
	"io"
	"os"
)

type IOStreams struct {
	In  io.Reader
	Out io.Writer
	Err io.Writer
}

func System() *IOStreams {
	return &IOStreams{
		In:  os.Stdin,
		Out: os.Stdout,
		Err: os.Stderr,
	}
}

func Test() (*IOStreams, *bytes.Buffer, *bytes.Buffer, *bytes.Buffer) {
	in := &bytes.Buffer{}
	out := &bytes.Buffer{}
	errBuf := &bytes.Buffer{}
	return &IOStreams{In: in, Out: out, Err: errBuf}, in, out, errBuf
}
