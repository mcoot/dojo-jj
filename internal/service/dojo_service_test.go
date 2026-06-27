package service

import (
	"testing"

	"github.com/mcoot/dojo-jj/internal/models"
	"github.com/stretchr/testify/suite"
)

type DojoServiceSuite struct {
	suite.Suite
	*DojoServiceTestFixture
}

func TestServiceSuite(t *testing.T) {
	suite.Run(t, &DojoServiceSuite{})
}

func (s *DojoServiceSuite) SetupTest() {
	s.DojoServiceTestFixture = NewServiceTestFixture()
}

func (s *DojoServiceSuite) Test_WhenJJNotOnPath_ThenGetWorkspaceShouldReturnError() {
	s.JJClient.SetJJInstalled(false)

	err := s.service.GetWorkspace()

	s.requireErrorWithCode(err, models.ErrJJNotOnPath)
}

func (s *DojoServiceSuite) requireErrorWithCode(err error, code models.ErrorCode) {
	s.T().Helper()
	s.Require().Error(err)
	s.Require().IsType(&models.DojoError{}, err)
	s.Require().Equal(code, err.(*models.DojoError).Code)
}
