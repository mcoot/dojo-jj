package factory

import (
	"github.com/mcoot/dojo-jj/internal/dependencies"
	"github.com/mcoot/dojo-jj/internal/service"
)

type App struct {
	DojoService *service.DojoService
}

func BuildApp() (*App, error) {
	filesystemClient := dependencies.NewFileSystemClient()
	jjClient := dependencies.NewJJClient()

	dojoService := service.NewDojoService(filesystemClient, jjClient)

	return &App{
		DojoService: dojoService,
	}, nil
}
