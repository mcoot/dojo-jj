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

	dojoService := service.NewDojoService(filesystemClient)

	return &App{
		DojoService: dojoService,
	}, nil
}
