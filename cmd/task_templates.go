package cmd

import (
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

var taskTemplatesCmd = &cobra.Command{
	Use:   "task-templates",
	Short: "Task templates",
}

var generateTaskTemplatesCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate task templates",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		flags := cmd.Flags()
		client, err := getClient(flags)
		if err != nil {
			log.Fatalf("Failed to create client: %s", err)
		}
		err = client.ImportTaskTemplates(args...)
		if err != nil {
			log.Fatalf("Failed to import task templates: %s", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(taskTemplatesCmd)

	rootCmd.PersistentFlags().StringP("workdir", "w", "", "Work directory")
	rootCmd.MarkPersistentFlagRequired("workdir")

	taskTemplatesCmd.AddCommand(generateTaskTemplatesCmd)
}
