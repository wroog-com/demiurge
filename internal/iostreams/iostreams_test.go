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

	_, _ = ios.Out.Write([]byte("via ios out"))
	if out.String() != "via ios out" {
		t.Errorf("out: expected %q, got %q", "via ios out", out.String())
	}

	_, _ = ios.Err.Write([]byte("via ios err"))
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

func TestColorEnabled_envMatrix(t *testing.T) {
	tests := []struct {
		name string
		tty  bool
		env  map[string]string
		want bool
	}{
		{"tty, no signals", true, nil, true},
		{"no tty, no signals", false, nil, false},
		{"NO_COLOR beats tty", true, map[string]string{"NO_COLOR": "1"}, false},
		{"CLICOLOR=0 disables on tty", true, map[string]string{"CLICOLOR": "0"}, false},
		{"CLICOLOR=1 does not force through pipe", false, map[string]string{"CLICOLOR": "1"}, false},
		{"TERM=dumb disables on tty", true, map[string]string{"TERM": "dumb"}, false},
		{"CLICOLOR_FORCE forces through pipe", false, map[string]string{"CLICOLOR_FORCE": "1"}, true},
		{"CLICOLOR_FORCE=0 does not force", false, map[string]string{"CLICOLOR_FORCE": "0"}, false},
		{"CLICOLOR_FORCE beats NO_COLOR", false, map[string]string{"CLICOLOR_FORCE": "1", "NO_COLOR": "1"}, true},
		{"CLICOLOR_FORCE beats TERM=dumb", false, map[string]string{"CLICOLOR_FORCE": "1", "TERM": "dumb"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ios, _, _, _ := Test()
			if tt.tty {
				ios.ForceTerminal()
			}
			ios.SetEnv(tt.env)
			if got := ios.ColorEnabled(); got != tt.want {
				t.Errorf("ColorEnabled() = %v, want %v", got, tt.want)
			}
		})
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

func TestSystem_readsProcessEnv(t *testing.T) {
	t.Setenv("DEMI_DEBUG", "1")
	if !System().IsDebug() {
		t.Error("expected System() to read DEMI_DEBUG from the process environment")
	}
}

func TestIsDebug_unset(t *testing.T) {
	ios, _, _, _ := Test()
	if ios.IsDebug() {
		t.Error("expected IsDebug to be false when DEMI_DEBUG is unset")
	}
}

func TestIsDebug_falsy(t *testing.T) {
	for _, v := range []string{"false", "0", "no", "False", "NO", "FALSE"} {
		ios, _, _, _ := Test()
		ios.SetEnv(map[string]string{"DEMI_DEBUG": v})
		if ios.IsDebug() {
			t.Errorf("expected IsDebug to be false when DEMI_DEBUG=%q", v)
		}
	}
}

func TestIsDebug_set(t *testing.T) {
	ios, _, _, _ := Test()
	ios.SetEnv(map[string]string{"DEMI_DEBUG": "1"})
	if !ios.IsDebug() {
		t.Error("expected IsDebug to be true when DEMI_DEBUG=1")
	}
}

func TestDebugf_writesWhenEnabled(t *testing.T) {
	ios, _, _, errBuf := Test()
	ios.SetEnv(map[string]string{"DEMI_DEBUG": "1"})
	ios.Debugf("hello %s", "world")
	if errBuf.String() != "DEBUG: hello world\n" {
		t.Errorf("unexpected debug output: %q", errBuf.String())
	}
}

func TestDebugf_silentWhenDisabled(t *testing.T) {
	ios, _, _, errBuf := Test()
	ios.Debugf("should not appear")
	if errBuf.String() != "" {
		t.Errorf("expected no output, got: %q", errBuf.String())
	}
}
