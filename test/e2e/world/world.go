//go:build e2e

package world

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

const defaultCommandTimeout = 10 * time.Second

type World struct {
	Root           string
	HomeDir        string
	JJConfigPath   string
	DojoBin        string
	JJBin          string
	TemplateRepo   string
	CommandTimeout time.Duration

	system systemClient
}

type worldOptions struct {
	system systemClient
}

type systemClient interface {
	Getenv(string) string
	LookPath(string) (string, error)
}

type realSystemClient struct{}

func (c *realSystemClient) Getenv(key string) string {
	return os.Getenv(key)
}

func (c *realSystemClient) LookPath(file string) (string, error) {
	return exec.LookPath(file)
}

func New(t testing.TB) *World {
	t.Helper()

	w, err := newWorld(t, worldOptions{})
	if err != nil {
		t.Fatalf("set up e2e world: %v", err)
	}

	return w
}

func newWorld(t testing.TB, opts worldOptions) (*World, error) {
	t.Helper()

	system := opts.system
	if system == nil {
		system = &realSystemClient{}
	}

	dojoBin, err := resolveDojoBin(system)
	if err != nil {
		return nil, err
	}

	jjBin, err := system.LookPath("jj")
	if err != nil {
		return nil, fmt.Errorf("jj is required on PATH for e2e tests: %w", err)
	}

	root := t.TempDir()
	homeDir := filepath.Join(root, "home")
	jjConfigPath := filepath.Join(root, "jj-config", "config.toml")

	if err := os.MkdirAll(homeDir, 0o755); err != nil {
		return nil, fmt.Errorf("create e2e HOME %q: %w", homeDir, err)
	}

	if err := writeJJConfig(jjConfigPath); err != nil {
		return nil, err
	}

	return &World{
		Root:           root,
		HomeDir:        homeDir,
		JJConfigPath:   jjConfigPath,
		DojoBin:        dojoBin,
		JJBin:          jjBin,
		CommandTimeout: defaultCommandTimeout,
		system:         system,
	}, nil
}

func resolveDojoBin(system systemClient) (string, error) {
	raw := system.Getenv("DOJO_BIN")
	if raw == "" {
		return "", errors.New("DOJO_BIN must point to the compiled dojo binary")
	}

	dojoBin := raw
	if !filepath.IsAbs(dojoBin) {
		abs, err := filepath.Abs(dojoBin)
		if err != nil {
			return "", fmt.Errorf("resolve DOJO_BIN %q: %w", raw, err)
		}
		dojoBin = abs
	}

	info, err := os.Stat(dojoBin)
	if err != nil {
		return "", fmt.Errorf("DOJO_BIN %q is invalid at %q: %w", raw, dojoBin, err)
	}

	if info.IsDir() {
		return "", fmt.Errorf("DOJO_BIN %q is invalid at %q: path is a directory", raw, dojoBin)
	}

	if info.Mode().Perm()&0o111 == 0 {
		return "", fmt.Errorf("DOJO_BIN %q is invalid at %q: path is not executable", raw, dojoBin)
	}

	return dojoBin, nil
}
