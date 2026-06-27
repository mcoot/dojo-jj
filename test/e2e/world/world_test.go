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

	"github.com/stretchr/testify/suite"
)

type WorldSuite struct {
	suite.Suite
}

func TestWorldSuite(t *testing.T) {
	suite.Run(t, new(WorldSuite))
}

func (s *WorldSuite) TestNew_WhenDojoBinMissing_ThenFailsWithHelpfulMessage() {
	system := s.newFakeSystemClient()

	_, err := newWorld(s.T(), worldOptions{
		system: system,
	})

	s.Require().Error(err)
	s.Assert().Contains(err.Error(), "DOJO_BIN")
}

func (s *WorldSuite) TestNew_WhenDojoBinNotExecutable_ThenFailsWithHelpfulMessage() {
	dojoBin := filepath.Join(s.T().TempDir(), "dojo")
	s.Require().NoError(os.WriteFile(dojoBin, []byte("not executable"), 0o644))
	system := s.newFakeSystemClient()
	system.SetEnv("DOJO_BIN", dojoBin)

	_, err := newWorld(s.T(), worldOptions{
		system: system,
	})

	s.Require().Error(err)
	s.Assert().Contains(err.Error(), "DOJO_BIN")
	s.Assert().Contains(err.Error(), dojoBin)
}

func (s *WorldSuite) TestNew_WhenDojoBinRelative_ThenStoresAbsolutePath() {
	root := s.T().TempDir()
	dojoBin := filepath.Join(root, "bin", "dojo")
	s.Require().NoError(os.MkdirAll(filepath.Dir(dojoBin), 0o755))
	s.Require().NoError(os.WriteFile(dojoBin, []byte(""), 0o755))
	s.Require().NoError(os.Chmod(dojoBin, 0o755))
	s.T().Chdir(root)
	system := s.newFakeSystemClient()
	system.SetEnv("DOJO_BIN", filepath.Join("bin", "dojo"))

	w, err := newWorld(s.T(), worldOptions{
		system: system,
	})

	s.Require().NoError(err)
	s.Assert().True(filepath.IsAbs(w.DojoBin))
	s.Assert().Equal(dojoBin, w.DojoBin)
}

func (s *WorldSuite) TestNew_WhenJJMissing_ThenFailsWithHelpfulMessage() {
	dojoBin := filepath.Join(s.T().TempDir(), "dojo")
	s.Require().NoError(os.WriteFile(dojoBin, []byte(""), 0o755))
	s.Require().NoError(os.Chmod(dojoBin, 0o755))
	system := s.newFakeSystemClient()
	system.SetEnv("DOJO_BIN", dojoBin)
	system.SetLookPathError("jj", errors.New("not found"))

	_, err := newWorld(s.T(), worldOptions{
		system: system,
	})

	s.Require().Error(err)
	s.Assert().Contains(err.Error(), "jj is required on PATH")
}

func (s *WorldSuite) TestRunCommand_WhenCommandSucceeds_ThenCapturesStdoutAndExitCode() {
	w := s.newWorldWithHelperDojo()
	dir := s.T().TempDir()

	result := w.RunDojoWithEnv(s.T(), dir, helperEnv(), "-test.run=TestHelperProcess", "--", "success")

	result.RequireSuccess(s.T())
	s.Assert().Equal(0, result.ExitCode)
	s.Assert().Contains(result.Stdout, "helper stdout")
	s.Assert().Empty(result.Stderr)
	s.Assert().Equal(dir, result.Dir)
	s.Assert().Equal(w.DojoBin, result.Args[0])
}

func (s *WorldSuite) TestRunCommand_WhenCommandFails_ThenCapturesStderrAndExitCode() {
	w := s.newWorldWithHelperDojo()

	result := w.RunDojoWithEnv(s.T(), s.T().TempDir(), helperEnv(), "-test.run=TestHelperProcess", "--", "failure")

	result.RequireFailure(s.T())
	s.Assert().Equal(17, result.ExitCode)
	s.Assert().Contains(result.Stderr, "helper stderr")
	s.Assert().Error(result.Err)
}

func (s *WorldSuite) TestEnvironment_ThenUsesTempHomeJJConfigAndNoColor() {
	w := s.newWorldWithHelperDojo()
	emptyPath := s.T().TempDir()

	env := environToMap(w.env(map[string]string{
		"PATH": emptyPath,
	}))

	s.Assert().Equal(w.DojoBin, env["DOJO_BIN"])
	s.Assert().Equal(w.HomeDir, env["HOME"])
	s.Assert().Equal(w.JJConfigPath, env["JJ_CONFIG"])
	s.Assert().Equal("1", env["NO_COLOR"])
	s.Assert().Equal(emptyPath, env["PATH"])
}

func (s *WorldSuite) TestWriteJJConfig_ThenContainsDeterministicSettings() {
	w := s.newWorldWithHelperDojo()

	config, err := os.ReadFile(w.JJConfigPath)
	s.Require().NoError(err)

	text := string(config)
	s.Assert().Contains(text, `user.name = "Dojo E2E"`)
	s.Assert().Contains(text, `user.email = "dojo-e2e@example.invalid"`)
	s.Assert().Contains(text, `ui.color = "never"`)
	s.Assert().Contains(text, `git.colocate = true`)
	s.Assert().Contains(text, `behavior = "drop"`)
	s.Assert().Contains(text, `backend = "none"`)
}

func (s *WorldSuite) TestTemplateRepo_WhenCreated_ThenJJStatusSucceeds() {
	w := s.newWorldWithRealJJ()

	templateRepo := w.ensureTemplateRepo(s.T())
	result := w.RunJJ(s.T(), templateRepo, "status")

	result.RequireSuccess(s.T())
	s.Assert().DirExists(filepath.Join(templateRepo, ".jj"))
	s.Assert().DirExists(filepath.Join(templateRepo, ".git"))
}

func (s *WorldSuite) TestNewRepoFromTemplate_ThenCopyHasJJAndGitMetadata() {
	w := s.newWorldWithRealJJ()

	repo := w.NewRepoFromTemplate(s.T(), "metadata")
	result := w.RunJJ(s.T(), repo, "status")

	result.RequireSuccess(s.T())
	s.Assert().DirExists(filepath.Join(repo, ".jj"))
	s.Assert().DirExists(filepath.Join(repo, ".git"))
}

func (s *WorldSuite) TestNewRepoFromTemplate_WhenCopyMutates_ThenTemplateUnchanged() {
	w := s.newWorldWithRealJJ()

	repo := w.NewRepoFromTemplate(s.T(), "copy")
	copyOnly := filepath.Join(repo, "copy-only.txt")
	s.Require().NoError(os.WriteFile(copyOnly, []byte("copy"), 0o644))

	s.Assert().NoFileExists(filepath.Join(w.TemplateRepo, "copy-only.txt"))
}

func (s *WorldSuite) TestNewRepoFromTemplate_WhenTwoCopiesMutate_ThenCopiesAreIndependent() {
	w := s.newWorldWithRealJJ()

	first := w.NewRepoFromTemplate(s.T(), "first")
	second := w.NewRepoFromTemplate(s.T(), "second")
	s.Require().NoError(os.WriteFile(filepath.Join(first, "value.txt"), []byte("first"), 0o644))
	s.Require().NoError(os.WriteFile(filepath.Join(second, "value.txt"), []byte("second"), 0o644))

	firstValue, err := os.ReadFile(filepath.Join(first, "value.txt"))
	s.Require().NoError(err)
	secondValue, err := os.ReadFile(filepath.Join(second, "value.txt"))
	s.Require().NoError(err)

	s.Assert().Equal("first", string(firstValue))
	s.Assert().Equal("second", string(secondValue))
}

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}

	mode := helperMode()
	switch mode {
	case "success":
		fmt.Fprintln(os.Stdout, "helper stdout")
		os.Exit(0)
	case "failure":
		fmt.Fprintln(os.Stderr, "helper stderr")
		os.Exit(17)
	default:
		fmt.Fprintf(os.Stderr, "unknown helper mode %q\n", mode)
		os.Exit(2)
	}
}

func (s *WorldSuite) newWorldWithHelperDojo() *World {
	dojoBin, err := os.Executable()
	s.Require().NoError(err)
	system := s.newFakeSystemClient()
	system.SetEnv("DOJO_BIN", dojoBin)

	w, err := newWorld(s.T(), worldOptions{
		system: system,
	})
	s.Require().NoError(err)
	w.CommandTimeout = 2 * time.Second

	return w
}

func (s *WorldSuite) newWorldWithRealJJ() *World {
	dojoBin, err := os.Executable()
	s.Require().NoError(err)
	jjBin, err := exec.LookPath("jj")
	s.Require().NoError(err)
	system := s.newFakeSystemClient()
	system.SetEnv("DOJO_BIN", dojoBin)
	system.SetLookPath("jj", jjBin)

	w, err := newWorld(s.T(), worldOptions{
		system: system,
	})
	s.Require().NoError(err)

	return w
}

type fakeSystemClient struct {
	env        map[string]string
	pathValues map[string]string
	pathErrors map[string]error
}

func (s *WorldSuite) newFakeSystemClient() *fakeSystemClient {
	exe, err := os.Executable()
	s.Require().NoError(err)

	client := &fakeSystemClient{
		env: map[string]string{
			"PATH": os.Getenv("PATH"),
		},
		pathValues: map[string]string{},
		pathErrors: map[string]error{},
	}
	client.SetLookPath("jj", exe)

	return client
}

func (c *fakeSystemClient) Getenv(key string) string {
	return c.env[key]
}

func (c *fakeSystemClient) LookPath(file string) (string, error) {
	if err, ok := c.pathErrors[file]; ok {
		return "", err
	}

	if path, ok := c.pathValues[file]; ok {
		return path, nil
	}

	return "", errors.New("not found")
}

func (c *fakeSystemClient) SetEnv(key string, value string) {
	c.env[key] = value
}

func (c *fakeSystemClient) SetLookPath(file string, path string) {
	c.pathValues[file] = path
	delete(c.pathErrors, file)
}

func (c *fakeSystemClient) SetLookPathError(file string, err error) {
	c.pathErrors[file] = err
	delete(c.pathValues, file)
}

func helperEnv() map[string]string {
	return map[string]string{
		"GO_WANT_HELPER_PROCESS": "1",
	}
}

func helperMode() string {
	for index, arg := range os.Args {
		if arg == "--" && index+1 < len(os.Args) {
			return os.Args[index+1]
		}
	}

	return ""
}
