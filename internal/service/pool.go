package service

import (
	"errors"

	"github.com/mcoot/dojo-jj/internal/dependencies"
	"github.com/mcoot/dojo-jj/internal/models"
)

type WorkspacePoolService struct {
	appConfig        *models.AppConfig
	filesystemClient dependencies.FileSystemClient
	jjClient         dependencies.JJClient
}

func NewWorkspacePoolService(
	appConfig *models.AppConfig,
	filesystemClient dependencies.FileSystemClient,
	jjClient dependencies.JJClient,
) *WorkspacePoolService {
	return &WorkspacePoolService{
		appConfig:        appConfig,
		filesystemClient: filesystemClient,
		jjClient:         jjClient,
	}
}

func (s *WorkspacePoolService) Acquire(repo models.RepoRoot) (*models.WorkspaceLease, error) {
	// Find the pool for this repo and read the state

	// If there are no free workspaces, create a new one

	// Take the lease on the lowest-ordered workspace

	// Return the workspace
	return nil, errors.New("not implemented")
}

func (s *WorkspacePoolService) Release(lease models.WorkspaceLease) {
	// Release the lease on the workspace
}
