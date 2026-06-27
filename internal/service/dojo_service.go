package service

import (
	"github.com/mcoot/dojo-jj/internal/dependencies"
	"github.com/mcoot/dojo-jj/internal/models"
)

type DojoService struct {
	filesystemClient dependencies.FileSystemClient
	jjClient         dependencies.JJClient
}

func NewDojoService(filesystemClient dependencies.FileSystemClient, jjClient dependencies.JJClient) *DojoService {
	return &DojoService{
		filesystemClient: filesystemClient,
		jjClient:         jjClient,
	}
}

func (s *DojoService) GetWorkspace() error {
	if !s.jjClient.IsJJAvailable() {
		return models.NewDojoError(models.ErrJJNotOnPath, "JJ not found on path")
	}

	_, err := s.jjClient.ListWorkspaces()
	if err != nil {
		return models.NewDojoErrorWithCause(
			models.ErrJJFailedToListWorkspaces,
			"Failed to list workspaces",
			err,
		)
	}

	return nil
}
