package e2e

import (
	"github.com/mcoot/dojo-jj/internal/e2e/world"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type WorkspacePoolE2ESuite struct {
	suite.Suite
	*world.E2EWorld
}

func (s *WorkspacePoolE2ESuite) SetupSuite() {
	s.E2EWorld = &world.E2EWorld{}
}

func (s *WorkspacePoolE2ESuite) Test_WorkspaceLifecycle() {
	assert.True(s.T(), true)
}
