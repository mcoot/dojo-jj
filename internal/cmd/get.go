package cmd

import (
	"github.com/mcoot/dojo-jj/internal/factory"
	"github.com/spf13/cobra"
)

type GetCmd struct {
}

func (c *GetCmd) Run(app *factory.App) error {
	return app.DojoService.GetWorkspace()
}

func (c *GetCmd) Mount(parent *cobra.Command, app *factory.App) {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "launch a subshell in a workspace",
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.Run(app)
		},
	}

	parent.AddCommand(cmd)
}
