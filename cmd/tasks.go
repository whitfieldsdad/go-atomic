package cmd

import (
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

var tasksCmd = &cobra.Command{
	Use:   "tasks",
	Short: "Tasks",
}

var listTasksCmd = &cobra.Command{
	Use:   "list",
	Short: "List tasks",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		flags := cmd.Flags()
		client, err := getClient(flags)
		if err != nil {
			log.Fatalf("Failed to create client: %s", err)
		}
		i, _ := flags.GetString("input-path")
		tasks, err := client.ReadTasks(i)
		if err != nil {
			log.Fatalf("Failed to read tasks: %s", err)
		}
		printJson(tasks)
	},
}

var countTasksCmd = &cobra.Command{
	Use:   "count",
	Short: "Count tasks",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

var generateTasksCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate tasks",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		flags := cmd.Flags()
		client, err := getClient(flags)
		if err != nil {
			log.Fatalf("Failed to create client: %s", err)
		}
		i, _ := flags.GetString("input-path")
		o, _ := flags.GetString("output-path")
		err = client.GenerateTasks(i, o)
		if err != nil {
			log.Fatalf("Failed to generate tasks: %s", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(tasksCmd)
	tasksCmd.AddCommand(listTasksCmd, countTasksCmd, generateTasksCmd)

	listTasksCmd.Flags().StringP("input-path", "i", "", "Input file/directory")
	listTasksCmd.MarkFlagRequired("input-path")

	generateTasksCmd.Flags().StringP("input-path", "i", "", "Input file/directory")
	generateTasksCmd.Flags().StringP("output-path", "o", "", "Output directory")
	generateTasksCmd.MarkFlagRequired("input-path")
	generateTasksCmd.MarkFlagRequired("output-path")
}
