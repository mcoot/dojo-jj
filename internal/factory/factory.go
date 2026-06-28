package factory

import (
	"github.com/mcoot/dojo-jj/internal/dependencies"
	"github.com/mcoot/dojo-jj/internal/service"
	"github.com/mcoot/dojo-jj/internal/service/appconfig"
)

type App struct {
	DojoService *service.DojoService
}

func BuildApp() (*App, error) {
	filesystemClient := dependencies.NewFileSystemClient()
	jjClient := dependencies.NewJJClient()

	appconfigLoader := appconfig.NewLoader()

	appConfig, err := appconfigLoader.Load()
	if err != nil {
		return nil, err
	}

	workspacePoolService := service.NewWorkspacePoolService(appConfig, filesystemClient, jjClient)

	dojoService := service.NewDojoService(appConfig, jjClient, workspacePoolService)

	return &App{
		DojoService: dojoService,
	}, nil
}
