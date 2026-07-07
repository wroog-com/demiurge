package iostreams

import "testing"

// A zero-value IOStreams has nil file handles; the TTY checks must guard
// against that and report false rather than calling term.IsTerminal on fd 0.
func TestIsTTY_nilFileReturnsFalse(t *testing.T) {
	s := &IOStreams{}

	if s.IsStdinTTY() {
		t.Error("expected IsStdinTTY to be false when inFile is nil")
	}
	if s.IsStdoutTTY() {
		t.Error("expected IsStdoutTTY to be false when outFile is nil")
	}
	if s.IsStderrTTY() {
		t.Error("expected IsStderrTTY to be false when errFile is nil")
	}
}

func TestForceNoTerminal(t *testing.T) {
	ios, _, _, _ := Test()
	ios.ForceTerminal()
	ios.ForceNoTerminal()

	if ios.IsStdinTTY() {
		t.Error("expected IsStdinTTY to be false after ForceNoTerminal")
	}
	if ios.IsStdoutTTY() {
		t.Error("expected IsStdoutTTY to be false after ForceNoTerminal")
	}
	if ios.IsStderrTTY() {
		t.Error("expected IsStderrTTY to be false after ForceNoTerminal")
	}
}

// System() wires real file handles, so the TTY checks take the term.IsTerminal
// path. The result depends on the environment; assert only that it does not panic.
func TestSystem_isTTYUsesRealFds(t *testing.T) {
	s := System()
	_ = s.IsStdinTTY()
	_ = s.IsStdoutTTY()
	_ = s.IsStderrTTY()
}
