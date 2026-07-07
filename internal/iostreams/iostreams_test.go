package iostreams

import (
	"os"
	"testing"
)

func TestSystem(t *testing.T) {
	s := System()
	if s.In != os.Stdin {
		t.Error("expected In to be os.Stdin")
	}
	if s.Out != os.Stdout {
		t.Error("expected Out to be os.Stdout")
	}
	if s.Err != os.Stderr {
		t.Error("expected Err to be os.Stderr")
	}
}

func TestTest_buffersAreLinked(t *testing.T) {
	ios, _, out, errBuf := Test()

	ios.Out.Write([]byte("via ios out"))
	if out.String() != "via ios out" {
		t.Errorf("out: expected %q, got %q", "via ios out", out.String())
	}

	ios.Err.Write([]byte("via ios err"))
	if errBuf.String() != "via ios err" {
		t.Errorf("errBuf: expected %q, got %q", "via ios err", errBuf.String())
	}
}
