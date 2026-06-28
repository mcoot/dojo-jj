package service

import (
	"github.com/mcoot/dojo-jj/internal/dependencies/mocks"
	"github.com/mcoot/dojo-jj/internal/models"
)

type DojoServiceTestFixture struct {
	filesystemClient *mocks.FileSystemClient
	JJClient         *mocks.JJClient

	service *DojoService
}

func NewServiceTestFixture() *DojoServiceTestFixture {
	testConfig := &models.AppConfig{
		RootDir: "/test-root-dir",
	}

	filesystemClient := mocks.NewFileSystemClient()
	JJClient := mocks.NewJJClient()

	workspacePoolService := NewWorkspacePoolService(testConfig, filesystemClient, JJClient)
	service := NewDojoService(testConfig, JJClient, workspacePoolService)

	return &DojoServiceTestFixture{
		filesystemClient: filesystemClient,
		JJClient:         JJClient,
		service:          service,
	}
}
