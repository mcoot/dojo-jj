package main

import (
	"log"

	"github.com/mcoot/dojo-jj/internal/cmd"
	"github.com/mcoot/dojo-jj/internal/factory"
)

func main() {
	app, err := factory.BuildApp()
	if err != nil {
		log.Fatal(err)
	}

	cli := cmd.BuildCli(app)

	err = cli.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
