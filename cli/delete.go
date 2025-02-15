package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/dominikbraun/timetrace/core"
	"github.com/dominikbraun/timetrace/out"

	"github.com/spf13/cobra"
)

var confirmed bool

type deleteOptions struct {
	Revert bool
}

func deleteCommand(t *core.Timetrace) *cobra.Command {
	delete := &cobra.Command{
		Use:   "delete",
		Short: "Delete a resource",
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	delete.AddCommand(deleteProjectCommand(t))
	delete.AddCommand(deleteRecordCommand(t))
	delete.PersistentFlags().BoolVar(&confirmed, "yes", false, "Do not ask for confirmation")

	return delete
}

func deleteProjectCommand(t *core.Timetrace) *cobra.Command {
	var options deleteOptions
	deleteProject := &cobra.Command{
		Use:   "project <KEY>",
		Short: "Delete a project",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			key := args[0]
			if options.Revert {
				if err := t.RevertProject(key); err != nil {
					out.Err("Failed to revert project: %s", err.Error())
				} else {
					out.Info("Project backup restored successfully")
				}
				return
			}

			project := core.Project{
				Key: key,
			}

			if !askForConfirmation() {
				out.Info("Record NOT deleted.")
				return
			}

			if err := t.BackupProject(key); err != nil {
				out.Err("Failed to backup project before deletion: %s", err.Error())
				return
			}

			if err := t.DeleteProject(project); err != nil {
				out.Err("Failed to delete %s", err.Error())
				return
			}

			out.Success("Deleted project %s", key)
		},
	}

	deleteProject.PersistentFlags().BoolVarP(&options.Revert, "revert", "r", false, "Restores the record to it's state prior to the last 'delete' command.")

	return deleteProject
}

func deleteRecordCommand(t *core.Timetrace) *cobra.Command {
	var options deleteOptions
	// Depending on the use12hours setting, the command syntax either is
	// `record YYYY-MM-DD-HH-MM` or `record YYYY-MM-DD-HH-MMPM`.
	use := fmt.Sprintf("record %s", t.Formatter().RecordKeyLayout())

	deleteRecord := &cobra.Command{
		Use:   use,
		Short: "Delete a record",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			start, err := t.Formatter().ParseRecordKey(args[0])
			if err != nil {
				out.Err("Failed to parse date argument: %s", err.Error())
				return
			}

			if options.Revert {
				if err := t.RevertRecord(start); err != nil {
					out.Err("Failed to revert record: %s", err.Error())
				} else {
					out.Info("Record backup restored successfully")
				}
				return
			}

			record, err := t.LoadRecord(start)
			if err != nil {
				out.Err("Failed to read record: %s", err.Error())
				return
			}

			showRecord(record, t.Formatter())
			if !confirmed {
				if !askForConfirmation() {
					out.Info("Record NOT deleted.")
					return
				}
			}

			if err := t.BackupRecord(start); err != nil {
				out.Err("Failed to backup record before deletion: %s", err.Error())
				return
			}

			if err := t.DeleteRecord(*record); err != nil {
				out.Err("Failed to delete %s", err.Error())
				return
			}

			out.Success("Deleted record %s", args[0])
		},
	}

	deleteRecord.PersistentFlags().BoolVarP(&options.Revert, "revert", "r", false, "Restores the record to it's state prior to the last 'delete' command.")

	return deleteRecord
}

func askForConfirmation() bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Fprint(os.Stderr, "Please confirm [y/N]: ")
	s, _ := reader.ReadString('\n')
	s = strings.TrimSuffix(s, "\n")
	s = strings.ToLower(s)

	return s == "y"
}
