package appconfig

import "github.com/mcoot/dojo-jj/internal/models"

type Loader struct{}

func NewLoader() *Loader {
	return &Loader{}
}

func (loader *Loader) Load() (*models.AppConfig, error) {
	return &models.AppConfig{
		RootDir: "~/.dojo",
	}, nil
}
