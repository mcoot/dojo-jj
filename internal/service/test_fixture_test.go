package service

import "github.com/mcoot/dojo-jj/internal/dependencies/mocks"

type DojoServiceTestFixture struct {
	filesystemClient *mocks.FileSystemClient
	JJClient         *mocks.JJClient

	service *DojoService
}

func NewServiceTestFixture() *DojoServiceTestFixture {
	filesystemClient := mocks.NewFileSystemClient()
	JJClient := mocks.NewJJClient()

	service := NewDojoService(filesystemClient, JJClient)

	return &DojoServiceTestFixture{
		filesystemClient: filesystemClient,
		JJClient:         JJClient,
		service:          service,
	}
}
