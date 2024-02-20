package cmd

import (
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/whitfieldsdad/go-atomic-red-team/pkg/atomic"
)

var tasksCmd = &cobra.Command{
	Use:   "tasks",
	Short: "Task commands",
}

var generateTasksCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate tasks from a tarball or directory",
	Run: func(cmd *cobra.Command, args []string) {
		inputPath, _ := cmd.Flags().GetString("input-path")
		outputPath, _ := cmd.Flags().GetString("output-path")
		client := atomic.NewClient()
		err := client.ImportTaskTemplates(inputPath, outputPath)
		if err != nil {
			log.Fatalf("Failed to import task templates: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(tasksCmd)
	tasksCmd.AddCommand(generateTasksCmd)

	generateTasksCmd.Flags().StringP("input-path", "i", "", "Input file/directory")
	generateTasksCmd.Flags().StringP("output-path", "o", "", "Output directory")

	generateTasksCmd.MarkFlagRequired("input-path")
	generateTasksCmd.MarkFlagRequired("output-path")
}
