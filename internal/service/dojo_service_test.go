package service

import (
	"errors"
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

func (s *DojoServiceSuite) Test_GetWorkspace_WhenJJNotOnPath_Fails() {
	s.JJClient.SetJJInstalled(false)

	err := s.service.GetWorkspace()

	s.requireErrorWithCode(err, models.ErrJJNotOnPath)
}

func (s *DojoServiceSuite) Test_GetWorkspace_WhenNotInJJRepo_Fails() {
	s.JJClient.SetGetRepoRootError(
		errors.New("Error: There is no jj repo in \".\""),
	)

	err := s.service.GetWorkspace()

	s.requireErrorWithCode(err, models.ErrNotInJJRepo)
}

func (s *DojoServiceSuite) Test_GetWorkspace_WhenGettingJJRootFails_Fails() {
	s.JJClient.SetGetRepoRootError(errors.New("Error: something bad happened"))

	err := s.service.GetWorkspace()

	s.requireErrorWithCode(err, models.ErrJJGetRootFailed)
}

func (s *DojoServiceSuite) requireErrorWithCode(err error, code models.ErrorCode) {
	s.T().Helper()
	s.Require().Error(err)
	s.Require().IsType(&models.DojoError{}, err)
	s.Require().Equal(code, err.(*models.DojoError).Code)
}
