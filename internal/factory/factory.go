package factory

type App struct{}

func BuildApp() (*App, error) {
	return &App{}, nil
}
