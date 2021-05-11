package cli

import (
	"time"

	"github.com/dominikbraun/timetrace/core"
	"github.com/dominikbraun/timetrace/out"

	"github.com/spf13/cobra"
)

func getCommand() *cobra.Command {
	get := &cobra.Command{
		Use: "get",
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	get.AddCommand(getProjectCommand())
	get.AddCommand(getRecordCommand())

	return get
}

func getProjectCommand() *cobra.Command {
	getProject := &cobra.Command{
		Use:  "project <KEY>",
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			key := args[0]

			project, err := core.LoadProject(key)
			if err != nil {
				out.Err("Failed to get project: %s", key)
				return
			}

			out.Table([]string{"Key"}, [][]string{{project.Key}})
		},
	}

	return getProject
}

func getRecordCommand() *cobra.Command {
	getRecord := &cobra.Command{
		Use:  "record YYYY-MM-DD-HH-MM",
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			start, err := time.Parse("2006-01-02-15-04", args[0])
			if err != nil {
				out.Err("Failed to parse date argument: %s", err.Error())
				return
			}

			record, err := core.LoadRecord(start)
			if err != nil {
				out.Err("Failed to read record: %s", err.Error())
				return
			}

			isBillable := defaultBool

			if record.IsBillable {
				isBillable = "yes"
			}

			end := defaultString

			if record.End != nil {
				end = record.End.Format("15:04")
			}

			project := defaultString

			if record.Project != nil {
				project = record.Project.Key
			}

			rows := [][]string{
				{
					record.Start.Format("15:04"),
					end,
					project,
					isBillable,
				},
			}

			out.Table([]string{"Start", "End", "Project", "Billable"}, rows)
		},
	}

	return getRecord
}
