package cmd

import (
	"github.com/mcoot/dojo-jj/internal/factory"
	"github.com/spf13/cobra"
)

func BuildCli(app *factory.App) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "dojo",
		Short: "dojo manages Jujutsu workspaces",
	}

	(&GetCmd{}).Mount(rootCmd, app)

	return rootCmd
}
