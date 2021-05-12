package cli

import (
	"github.com/dominikbraun/timetrace/core"
	"github.com/dominikbraun/timetrace/out"

	"github.com/spf13/cobra"
)

func createCommand(t *core.Timetrace) *cobra.Command {
	create := &cobra.Command{
		Use:   "create",
		Short: "Create a new resource",
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	create.AddCommand(createProjectCommand(t))

	return create
}

func createProjectCommand(t *core.Timetrace) *cobra.Command {
	createProject := &cobra.Command{
		Use:   "project <KEY>",
		Short: "Create a new project",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			key := args[0]

			project := core.Project{
				Key: key,
			}

			if err := t.SaveProject(project, false); err != nil {
				out.Err("Failed to create project: %s", err.Error())
				return
			}

			out.Success("Created project %s", key)
		},
	}

	return createProject
}
