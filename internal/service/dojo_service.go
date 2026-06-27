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
		return &models.DojoError{
			Code:    models.ErrJJNotOnPath,
			Message: "JJ not found on path",
		}
	}

	return nil
}
