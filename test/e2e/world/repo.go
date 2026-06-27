//go:build e2e

package world

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"
)

const jjConfig = `user.name = "Dojo E2E"
user.email = "dojo-e2e@example.invalid"
ui.color = "never"
git.colocate = true

[signing]
behavior = "drop"
backend = "none"
`

func writeJJConfig(path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create jj config directory %q: %w", filepath.Dir(path), err)
	}

	if err := os.WriteFile(path, []byte(jjConfig), 0o644); err != nil {
		return fmt.Errorf("write jj config %q: %w", path, err)
	}

	return nil
}

func (w *World) NewRepoFromTemplate(t testing.TB, name string) string {
	t.Helper()

	if filepath.Base(name) != name || name == "." || name == "" {
		t.Fatalf("invalid e2e repo name %q", name)
	}

	templateRepo := w.ensureTemplateRepo(t)
	repoRoot := filepath.Join(t.TempDir(), "repos", name)

	if err := os.MkdirAll(filepath.Dir(repoRoot), 0o755); err != nil {
		t.Fatalf("create e2e repo parent %q: %v", filepath.Dir(repoRoot), err)
	}

	if err := copyTree(templateRepo, repoRoot); err != nil {
		t.Fatalf("copy e2e repo template from %q to %q: %v", templateRepo, repoRoot, err)
	}

	return repoRoot
}

func (w *World) ensureTemplateRepo(t testing.TB) string {
	t.Helper()

	if w.TemplateRepo != "" {
		return w.TemplateRepo
	}

	templateRepo := filepath.Join(w.Root, "template-repo")
	result := w.RunJJ(t, w.Root, "git", "init", templateRepo)
	result.RequireSuccess(t)

	readmePath := filepath.Join(templateRepo, "README.md")
	if err := os.WriteFile(readmePath, []byte("# Dojo E2E Template\n"), 0o644); err != nil {
		t.Fatalf("write template fixture %q: %v", readmePath, err)
	}

	result = w.RunJJ(t, templateRepo, "describe", "-m", "template root")
	result.RequireSuccess(t)

	w.TemplateRepo = templateRepo

	return w.TemplateRepo
}

func copyTree(src string, dst string) error {
	info, err := os.Lstat(src)
	if err != nil {
		return err
	}

	switch {
	case info.Mode()&os.ModeSymlink != 0:
		target, err := os.Readlink(src)
		if err != nil {
			return err
		}
		return os.Symlink(target, dst)
	case info.IsDir():
		if err := os.MkdirAll(dst, info.Mode().Perm()); err != nil {
			return err
		}

		entries, err := os.ReadDir(src)
		if err != nil {
			return err
		}

		for _, entry := range entries {
			if err := copyTree(filepath.Join(src, entry.Name()), filepath.Join(dst, entry.Name())); err != nil {
				return err
			}
		}

		return os.Chmod(dst, info.Mode().Perm())
	case info.Mode().IsRegular():
		return copyFile(src, dst, info.Mode().Perm())
	default:
		return fmt.Errorf("unsupported file type at %q", src)
	}
}

func copyFile(src string, dst string, mode os.FileMode) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.OpenFile(dst, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, mode)
	if err != nil {
		return err
	}

	if _, err := io.Copy(out, in); err != nil {
		_ = out.Close()
		return err
	}

	return out.Close()
}
