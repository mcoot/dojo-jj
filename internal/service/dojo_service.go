package service

import (
	"errors"
	"strings"

	"github.com/mcoot/dojo-jj/internal/dependencies"
	"github.com/mcoot/dojo-jj/internal/models"
)

type DojoService struct {
	appConfig            *models.AppConfig
	jjClient             dependencies.JJClient
	workspacePoolService *WorkspacePoolService
}

func NewDojoService(
	appConfig *models.AppConfig,
	jjClient dependencies.JJClient,
	workspacePoolService *WorkspacePoolService,
) *DojoService {
	return &DojoService{
		appConfig:            appConfig,
		jjClient:             jjClient,
		workspacePoolService: workspacePoolService,
	}
}

func (s *DojoService) GetWorkspace() error {
	if !s.jjClient.IsJJAvailable() {
		return models.NewDojoError(models.ErrJJNotOnPath, "JJ not found on path")
	}

	// Get the current repo root
	_, err := s.jjClient.GetRepoRoot()
	if err != nil {
		if strings.Contains(err.Error(), "There is no jj repo") {
			return models.NewDojoError(models.ErrNotInJJRepo, "not in a jj repo")
		}
		return models.NewDojoErrorWithCause(models.ErrJJGetRootFailed, "failed to get repo root", err)
	}

	// Acquire a workspace from the pool
	//_, err := s.workspacePoolService.Acquire("some-repo")
	//if err != nil {
	//	return err
	//}

	// Print the command to use this workspace and release it
	return errors.New("not implemented")
}
