package service

import "github.com/mcoot/dojo-jj/internal/dependencies"

type DojoService struct {
	filesystemClient dependencies.FileSystemClient
}

func NewDojoService(filesystemClient dependencies.FileSystemClient) *DojoService {
	return &DojoService{
		filesystemClient: filesystemClient,
	}
}

func (s *DojoService) GetWorkspace() error {
	return nil
}
