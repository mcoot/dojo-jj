//go:build e2e

package world

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
	"testing"
	"time"
)

type CommandResult struct {
	Args     []string
	Dir      string
	Stdout   string
	Stderr   string
	ExitCode int
	Duration time.Duration
	Err      error
}

func (w *World) RunDojo(t testing.TB, dir string, args ...string) CommandResult {
	t.Helper()

	return w.RunDojoWithEnv(t, dir, nil, args...)
}

func (w *World) RunDojoWithEnv(
	t testing.TB,
	dir string,
	envOverrides map[string]string,
	args ...string,
) CommandResult {
	t.Helper()

	return w.runCommand(t, w.DojoBin, dir, envOverrides, args...)
}

func (w *World) RunJJ(t testing.TB, dir string, args ...string) CommandResult {
	t.Helper()

	return w.runCommand(t, w.JJBin, dir, nil, args...)
}

func (r CommandResult) RequireSuccess(t testing.TB) {
	t.Helper()

	if r.ExitCode != 0 || r.Err != nil {
		t.Fatalf("expected command to succeed\n%s", r.DebugString())
	}
}

func (r CommandResult) RequireFailure(t testing.TB) {
	t.Helper()

	if r.ExitCode == 0 && r.Err == nil {
		t.Fatalf("expected command to fail\n%s", r.DebugString())
	}
}

func (r CommandResult) DebugString() string {
	return fmt.Sprintf(
		"args: %s\ndir: %s\nexit_code: %d\nduration: %s\nerror: %v\nstdout:\n%s\nstderr:\n%s",
		strings.Join(r.Args, " "),
		r.Dir,
		r.ExitCode,
		r.Duration,
		r.Err,
		r.Stdout,
		r.Stderr,
	)
}

func (w *World) runCommand(
	t testing.TB,
	executable string,
	dir string,
	envOverrides map[string]string,
	args ...string,
) CommandResult {
	t.Helper()

	timeout := w.CommandTimeout
	if timeout == 0 {
		timeout = defaultCommandTimeout
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, executable, args...)
	cmd.Dir = dir
	cmd.Env = w.env(envOverrides)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	start := time.Now()
	err := cmd.Run()
	duration := time.Since(start)

	exitCode := 0
	if ctxErr := ctx.Err(); ctxErr != nil {
		err = ctxErr
		exitCode = -1
	} else if err != nil {
		exitCode = -1

		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			exitCode = exitErr.ExitCode()
		}
	}

	return CommandResult{
		Args:     append([]string{executable}, args...),
		Dir:      dir,
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
		ExitCode: exitCode,
		Duration: duration,
		Err:      err,
	}
}

func (w *World) env(overrides map[string]string) []string {
	values := environToMap(os.Environ())
	values["DOJO_BIN"] = w.DojoBin
	values["HOME"] = w.HomeDir
	values["JJ_CONFIG"] = w.JJConfigPath
	values["NO_COLOR"] = "1"
	values["PATH"] = w.system.Getenv("PATH")

	for key, value := range overrides {
		values[key] = value
	}

	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	env := make([]string, 0, len(keys))
	for _, key := range keys {
		env = append(env, key+"="+values[key])
	}

	return env
}

func environToMap(environ []string) map[string]string {
	values := make(map[string]string, len(environ))
	for _, entry := range environ {
		key, value, found := strings.Cut(entry, "=")
		if !found {
			continue
		}
		values[key] = value
	}

	return values
}
