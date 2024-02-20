package cmd

import "github.com/spf13/cobra"

var tasksCmd = &cobra.Command{
	Use:   "tasks",
	Short: "Task commands",
}

var generateTasksCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate tasks from a tarball or directory",
	Run: func(cmd *cobra.Command, args []string) {
		
	},
}

func init() {
	rootCmd.AddCommand(tasksCmd)
}
