//go:build e2e

package e2e

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mcoot/dojo-jj/test/e2e/world"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type WorkspacePoolE2ESuite struct {
	suite.Suite
	*world.World
}

func TestWorkspacePoolE2ESuite(t *testing.T) {
	suite.Run(t, new(WorkspacePoolE2ESuite))
}

func (s *WorkspacePoolE2ESuite) SetupSuite() {
	s.World = world.New(s.T())
}

func (s *WorkspacePoolE2ESuite) Test_DojoHelp_FromCopiedRepo() {
	repo := s.NewRepoFromTemplate(s.T(), "help")

	result := s.RunDojo(s.T(), repo, "--help")

	result.RequireSuccess(s.T())
	assert.Contains(s.T(), result.Stdout, "dojo manages Jujutsu workspaces")
	assert.Contains(s.T(), result.Stdout, "get")
}

func (s *WorkspacePoolE2ESuite) Test_DojoGet_WhenJJMissingFromPath_Fails() {
	repo := s.NewRepoFromTemplate(s.T(), "get-without-jj")
	emptyPath := s.T().TempDir()

	result := s.RunDojoWithEnv(s.T(), repo, map[string]string{
		"PATH": emptyPath,
	}, "get")

	result.RequireFailure(s.T())
	assert.Contains(s.T(), result.Stderr, "JJ not found on path")
}

func (s *WorkspacePoolE2ESuite) Test_DojoGet_WhenNotInJJRepo_Fails() {
	repo := s.NewRepoFromTemplate(s.T(), "get-without-jj")
	// Delete the .jj directory
	_ = os.RemoveAll(filepath.Join(repo, ".jj"))

	result := s.RunDojo(s.T(), repo, "get")

	result.RequireFailure(s.T())
	assert.Contains(s.T(), result.Stderr, "not in a jj repo")
}
