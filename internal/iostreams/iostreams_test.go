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

func TestTest_isNotTTY(t *testing.T) {
	ios, _, _, _ := Test()
	if ios.IsStdinTTY() {
		t.Error("expected IsStdinTTY to be false in test streams")
	}
	if ios.IsStdoutTTY() {
		t.Error("expected IsStdoutTTY to be false in test streams")
	}
	if ios.IsStderrTTY() {
		t.Error("expected IsStderrTTY to be false in test streams")
	}
}

func TestForceTerminal(t *testing.T) {
	ios, _, _, _ := Test()
	ios.ForceTerminal()
	if !ios.IsStdoutTTY() {
		t.Error("expected IsStdoutTTY to be true after ForceTerminal")
	}
	if !ios.IsStderrTTY() {
		t.Error("expected IsStderrTTY to be true after ForceTerminal")
	}
	if !ios.IsStdinTTY() {
		t.Error("expected IsStdinTTY to be true after ForceTerminal")
	}
}

func TestColorEnabled_noTTY(t *testing.T) {
	ios, _, _, _ := Test()
	if ios.ColorEnabled() {
		t.Error("expected ColorEnabled to be false when not a TTY")
	}
}

func TestColorEnabled_withTTY(t *testing.T) {
	ios, _, _, _ := Test()
	ios.ForceTerminal()
	if !ios.ColorEnabled() {
		t.Error("expected ColorEnabled to be true when TTY and NO_COLOR unset")
	}
}

func TestColorEnabled_noColor(t *testing.T) {
	ios, _, _, _ := Test()
	ios.ForceTerminal()
	t.Setenv("NO_COLOR", "1")
	if ios.ColorEnabled() {
		t.Error("expected ColorEnabled to be false when NO_COLOR is set")
	}
}

func TestForceColorOverrides(t *testing.T) {
	ios, _, _, _ := Test()

	ios.ForceColorEnabled()
	if !ios.ColorEnabled() {
		t.Error("expected ColorEnabled to be true after ForceColorEnabled")
	}

	ios.ForceColorDisabled()
	if ios.ColorEnabled() {
		t.Error("expected ColorEnabled to be false after ForceColorDisabled")
	}
}

func TestIsDebug_unset(t *testing.T) {
	t.Setenv("DEMI_DEBUG", "")
	ios, _, _, _ := Test()
	if ios.IsDebug() {
		t.Error("expected IsDebug to be false when DEMI_DEBUG is unset")
	}
}

func TestIsDebug_falsy(t *testing.T) {
	for _, v := range []string{"false", "0", "no"} {
		t.Setenv("DEMI_DEBUG", v)
		ios, _, _, _ := Test()
		if ios.IsDebug() {
			t.Errorf("expected IsDebug to be false when DEMI_DEBUG=%q", v)
		}
	}
}

func TestIsDebug_set(t *testing.T) {
	t.Setenv("DEMI_DEBUG", "1")
	ios, _, _, _ := Test()
	if !ios.IsDebug() {
		t.Error("expected IsDebug to be true when DEMI_DEBUG=1")
	}
}

func TestDebugf_writesWhenEnabled(t *testing.T) {
	t.Setenv("DEMI_DEBUG", "1")
	ios, _, _, errBuf := Test()
	ios.Debugf("hello %s", "world")
	if errBuf.String() != "DEBUG: hello world\n" {
		t.Errorf("unexpected debug output: %q", errBuf.String())
	}
}

func TestDebugf_silentWhenDisabled(t *testing.T) {
	t.Setenv("DEMI_DEBUG", "")
	ios, _, _, errBuf := Test()
	ios.Debugf("should not appear")
	if errBuf.String() != "" {
		t.Errorf("expected no output, got: %q", errBuf.String())
	}
}
