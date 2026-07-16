package demi

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/rogpeppe/go-internal/testscript"
	"github.com/spf13/cobra"
	"github.com/wroog-com/demiurge/internal/cmd"
	"github.com/wroog-com/demiurge/internal/cmdutil"
	"github.com/wroog-com/demiurge/internal/iostreams"
)

// TestMain hosts two subprocess entry points. The DEMI_TEST_SIGNAL_HANG branch
// is the hung-command helper for TestSignalContext_secondSignalKills and must
// stay first: testscript.Main never returns (it calls os.Exit), so if it ran
// first the re-exec'd helper would run the whole suite instead of hanging.
func TestMain(m *testing.M) {
	if os.Getenv("DEMI_TEST_SIGNAL_HANG") == "1" {
		ctx, stop := signalContext(context.Background())
		defer stop()
		fmt.Println("ready")
		go func() {
			<-ctx.Done()
			fmt.Println("canceled")
		}()
		time.Sleep(10 * time.Second) // deliberately ignores ctx
		os.Exit(0)
	}
	// os.Exit with the real code so scripts observe the 0/1/2 exit contract.
	testscript.Main(m, map[string]func(){
		"demi": func() { os.Exit(int(Main())) },
	})
}

// TestScripts drives the real binary through every script under testdata/scripts.
func TestScripts(t *testing.T) {
	testscript.Run(t, testscript.Params{
		Dir:                 "testdata/scripts",
		RequireExplicitExec: true,
		RequireUniqueNames:  true,
	})
}

func TestRun_success(t *testing.T) {
	t.Setenv("DEMI_DEBUG", "")
	ios, _, outBuf, errBuf := iostreams.Test()
	f := &cmdutil.Factory{IOStreams: ios}

	if code := run(context.Background(), f, []string{"version"}); code != ExitOK {
		t.Errorf("run() = %d, want ExitOK", code)
	}
	if got, want := outBuf.String(), "demi version dev\n"; got != want {
		t.Errorf("Out = %q, want %q", got, want)
	}
	if errBuf.String() != "" {
		t.Errorf("expected nothing on Err for success, got %q", errBuf.String())
	}
}

func TestRun_emptyArgsShowsHelp(t *testing.T) {
	t.Setenv("DEMI_DEBUG", "")
	ios, _, outBuf, errBuf := iostreams.Test()
	f := &cmdutil.Factory{IOStreams: ios}

	if code := run(context.Background(), f, []string{}); code != ExitOK {
		t.Errorf("run() = %d, want ExitOK", code)
	}
	if !strings.Contains(outBuf.String(), "Usage:") {
		t.Errorf("Out = %q, want root help", outBuf.String())
	}
	if errBuf.String() != "" {
		t.Errorf("expected nothing on Err, got %q", errBuf.String())
	}
}

func TestRun_nilArgsDoesNotReadOSArgs(t *testing.T) {
	t.Setenv("DEMI_DEBUG", "")
	ios, _, _, errBuf := iostreams.Test()
	f := &cmdutil.Factory{IOStreams: ios}

	orig := os.Args
	os.Args = []string{"demi", "nonsense"}
	defer func() { os.Args = orig }()

	if code := run(context.Background(), f, nil); code != ExitOK {
		t.Errorf("run() = %d, want ExitOK (help), not the ambient argv result", code)
	}
	if errBuf.String() != "" {
		t.Errorf("expected nothing on Err, got %q", errBuf.String())
	}
}

func TestRun_unknownCommand(t *testing.T) {
	t.Setenv("DEMI_DEBUG", "")
	ios, _, _, errBuf := iostreams.Test()
	f := &cmdutil.Factory{IOStreams: ios}

	if code := run(context.Background(), f, []string{"nonsense"}); code != ExitError {
		t.Errorf("run() = %d, want ExitError", code)
	}
	got := errBuf.String()
	if !strings.Contains(got, `unknown command "nonsense"`) {
		t.Errorf("Err = %q, want unknown-command message", got)
	}
	if !strings.Contains(got, "Usage:") {
		t.Errorf("Err = %q, want usage appended", got)
	}
}

// Cancellation is cooperative: a command that never reads its context
// completes normally.
func TestRun_canceledContextStillRunsVersion(t *testing.T) {
	t.Setenv("DEMI_DEBUG", "")
	ios, _, outBuf, _ := iostreams.Test()
	f := &cmdutil.Factory{IOStreams: ios}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if code := run(ctx, f, []string{"version"}); code != ExitOK {
		t.Errorf("run() = %d, want ExitOK", code)
	}
	if got, want := outBuf.String(), "demi version dev\n"; got != want {
		t.Errorf("Out = %q, want %q", got, want)
	}
}

func TestRun_debugStartupLine(t *testing.T) {
	ios, _, outBuf, errBuf := iostreams.Test()
	ios.SetEnv(map[string]string{"DEMI_DEBUG": "1"})
	f := &cmdutil.Factory{IOStreams: ios}

	if code := run(context.Background(), f, []string{"version"}); code != ExitOK {
		t.Errorf("run() = %d, want ExitOK", code)
	}
	if got, want := errBuf.String(), "DEBUG: demi dev args [\"version\"]\n"; got != want {
		t.Errorf("Err = %q, want %q", got, want)
	}
	if got, want := outBuf.String(), "demi version dev\n"; got != want {
		t.Errorf("Out = %q, want debug output to leave stdout untouched (%q)", got, want)
	}
}

// Mirrors run()'s wiring (run builds its root internally); keep the two in
// sync or this pin stops covering the real pipeline.
func TestExecute_contextCancellationMapsToExitCancel(t *testing.T) {
	ios, _, _, errBuf := iostreams.Test()
	f := &cmdutil.Factory{IOStreams: ios}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	root := cmd.NewRootCmd(f)
	root.AddCommand(&cobra.Command{
		Use: "block",
		RunE: func(c *cobra.Command, args []string) error {
			return c.Context().Err()
		},
	})
	root.SetArgs([]string{"block"})

	executed, err := root.ExecuteContextC(ctx)
	if code := mapError(f, executed, err); code != ExitCancel {
		t.Errorf("mapError() = %d, want ExitCancel", code)
	}
	if errBuf.String() != "" {
		t.Errorf("cancellation should be silent, got %q", errBuf.String())
	}
}

func TestMapError_nil(t *testing.T) {
	ios, _, _, errBuf := iostreams.Test()
	f := &cmdutil.Factory{IOStreams: ios}

	if code := mapError(f, nil, nil); code != ExitOK {
		t.Errorf("mapError(nil) = %d, want ExitOK", code)
	}
	if errBuf.String() != "" {
		t.Errorf("expected no output for nil error, got %q", errBuf.String())
	}
}

func TestMapError_printsError(t *testing.T) {
	ios, _, _, errBuf := iostreams.Test()
	f := &cmdutil.Factory{IOStreams: ios}

	code := mapError(f, nil, errors.New("something failed"))
	if code != ExitError {
		t.Errorf("mapError() = %d, want ExitError", code)
	}
	if got := errBuf.String(); got != "something failed\n" {
		t.Errorf("Err = %q, want %q", got, "something failed\n")
	}
}

func TestMapError_silentErrorNotPrinted(t *testing.T) {
	ios, _, _, errBuf := iostreams.Test()
	f := &cmdutil.Factory{IOStreams: ios}

	code := mapError(f, nil, cmdutil.ErrSilent)
	if code != ExitError {
		t.Errorf("mapError(SilentError) = %d, want ExitError", code)
	}
	if errBuf.String() != "" {
		t.Errorf("SilentError should be suppressed, got %q", errBuf.String())
	}
}

func TestMapError_wrappedSilentErrorNotPrinted(t *testing.T) {
	ios, _, _, errBuf := iostreams.Test()
	f := &cmdutil.Factory{IOStreams: ios}

	wrapped := fmt.Errorf("context: %w", cmdutil.ErrSilent)
	if code := mapError(f, nil, wrapped); code != ExitError {
		t.Errorf("mapError(wrapped SilentError) = %d, want ExitError", code)
	}
	if errBuf.String() != "" {
		t.Errorf("wrapped SilentError should be suppressed, got %q", errBuf.String())
	}
}

func TestMapError_cancellation(t *testing.T) {
	cases := map[string]error{
		"CancelError":      cmdutil.ErrCancel,
		"context.Canceled": context.Canceled,
		"wrapped":          fmt.Errorf("aborting: %w", cmdutil.ErrCancel),
	}
	for name, err := range cases {
		t.Run(name, func(t *testing.T) {
			ios, _, _, errBuf := iostreams.Test()
			f := &cmdutil.Factory{IOStreams: ios}

			if code := mapError(f, nil, err); code != ExitCancel {
				t.Errorf("mapError(%s) = %d, want ExitCancel", name, code)
			}
			if errBuf.String() != "" {
				t.Errorf("cancellation should be silent, got %q", errBuf.String())
			}
		})
	}
}

func TestMapError_deadlineIsPlainError(t *testing.T) {
	ios, _, _, errBuf := iostreams.Test()
	f := &cmdutil.Factory{IOStreams: ios}

	if code := mapError(f, nil, context.DeadlineExceeded); code != ExitError {
		t.Errorf("mapError(DeadlineExceeded) = %d, want ExitError", code)
	}
	if got := errBuf.String(); got != "context deadline exceeded\n" {
		t.Errorf("Err = %q, want deadline message", got)
	}
}

func TestPrintError_appendsUsageOnFlagError(t *testing.T) {
	ios, _, _, errBuf := iostreams.Test()
	c := &cobra.Command{Use: "demi"}

	err := cmdutil.FlagErrorf("unknown flag: --nope")
	printError(ios.Err, err, c)

	out := errBuf.String()
	if !strings.Contains(out, "unknown flag: --nope") {
		t.Errorf("expected error message in output, got %q", out)
	}
	if !strings.Contains(out, c.UsageString()) {
		t.Errorf("expected usage string appended for FlagError, got %q", out)
	}
}

func TestPrintError_appendsUsageOnUnknownCommand(t *testing.T) {
	ios, _, _, errBuf := iostreams.Test()
	c := &cobra.Command{Use: "demi"}

	err := errors.New("unknown command \"nope\" for \"demi\"")
	printError(ios.Err, err, c)

	if !strings.Contains(errBuf.String(), c.UsageString()) {
		t.Errorf("expected usage appended for unknown command, got %q", errBuf.String())
	}
}

func TestPrintError_plainErrorNoUsage(t *testing.T) {
	ios, _, _, errBuf := iostreams.Test()
	c := &cobra.Command{Use: "demi"}

	printError(ios.Err, errors.New("disk on fire"), c)

	out := errBuf.String()
	if out != "disk on fire\n" {
		t.Errorf("expected only the error line, got %q", out)
	}
}

func TestSignalContext_firstSignalCancels(t *testing.T) {
	for _, sig := range []syscall.Signal{syscall.SIGINT, syscall.SIGTERM} {
		t.Run(sig.String(), func(t *testing.T) {
			ctx, stop := signalContext(context.Background())
			defer stop()

			// Registration completes before signalContext returns, so the
			// self-delivered signal cannot race the handler.
			if err := syscall.Kill(syscall.Getpid(), sig); err != nil {
				t.Fatalf("kill: %v", err)
			}

			select {
			case <-ctx.Done():
			case <-time.After(5 * time.Second):
				t.Fatalf("context not canceled after %v", sig)
			}
			if !errors.Is(ctx.Err(), context.Canceled) {
				t.Errorf("ctx.Err() = %v, want context.Canceled", ctx.Err())
			}
		})
	}
}

// Runs against a re-exec'd copy of this test binary, so the shared test
// process is never signaled.
func TestSignalContext_secondSignalKills(t *testing.T) {
	child := exec.Command(os.Args[0])
	child.Env = append(os.Environ(), "DEMI_TEST_SIGNAL_HANG=1")
	stdout, err := child.StdoutPipe()
	if err != nil {
		t.Fatalf("stdout pipe: %v", err)
	}
	if err := child.Start(); err != nil {
		t.Fatalf("start helper: %v", err)
	}
	defer func() { _ = child.Process.Kill() }()

	scanner := bufio.NewScanner(stdout)
	expect := func(want string) {
		lines := make(chan string, 1)
		go func() {
			if scanner.Scan() {
				lines <- scanner.Text()
			} else {
				close(lines)
			}
		}()
		select {
		case got, ok := <-lines:
			if !ok || got != want {
				t.Fatalf("handshake: got %q (ok=%v), want %q", got, ok, want)
			}
		case <-time.After(5 * time.Second):
			t.Fatalf("timeout waiting for %q", want)
		}
	}

	expect("ready") // registration done: no startup race
	if err := child.Process.Signal(syscall.SIGINT); err != nil {
		t.Fatalf("first SIGINT: %v", err)
	}
	expect("canceled")                 // ctx canceled: watcher is running stop()
	time.Sleep(200 * time.Millisecond) // settle: disposition restore is async

	if err := child.Process.Signal(syscall.SIGINT); err != nil {
		t.Fatalf("second SIGINT: %v", err)
	}

	done := make(chan error, 1)
	go func() { done <- child.Wait() }()
	select {
	case err := <-done:
		var exitErr *exec.ExitError
		if !errors.As(err, &exitErr) {
			t.Fatalf("Wait() = %v, want signal-death ExitError", err)
		}
		ws, ok := exitErr.Sys().(syscall.WaitStatus)
		if !ok || !ws.Signaled() || ws.Signal() != syscall.SIGINT {
			t.Fatalf("child did not die by SIGINT: %v (status %#v)", err, exitErr.Sys())
		}
	case <-time.After(5 * time.Second):
		t.Fatal("child survived the second SIGINT")
	}
}
