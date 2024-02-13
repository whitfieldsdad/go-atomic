package cmd

import (
	"fmt"
	"strings"

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

var listTaskTemplatesCmd = &cobra.Command{
	Use:   "list",
	Short: "List task templates",
	Run: func(cmd *cobra.Command, args []string) {
		flags := cmd.Flags()
		client, err := getClient(flags)
		if err != nil {
			log.Fatalf("Failed to create client: %s", err)
		}
		query, err := getTaskTemplateQuery(flags)
		if err != nil {
			log.Fatalf("Failed to get task query: %s", err)
		}
		templates, err := client.ListTaskTemplates(query)
		if err != nil {
			log.Fatalf("Failed to list task templates: %s", err)
		}
		for _, template := range templates {
			fmt.Printf("- %s: %s (%s)\n", template.Id, template.Name, strings.Join(template.Tags, ", "))
		}
	},
}

var countTaskTemplatesCmd = &cobra.Command{
	Use:   "count",
	Short: "Count task templates",
	Run: func(cmd *cobra.Command, args []string) {
		flags := cmd.Flags()
		client, err := getClient(flags)
		if err != nil {
			log.Fatalf("Failed to create client: %s", err)
		}
		query, err := getTaskTemplateQuery(flags)
		if err != nil {
			log.Fatalf("Failed to get task query: %s", err)
		}
		total, err := client.CountTaskTemplates(query)
		if err != nil {
			log.Fatalf("Failed to count task templates: %s", err)
		}
		fmt.Println(total)
	},
}

// TODO
var runTaskTemplatesCmd = &cobra.Command{
	Use:   "run",
	Short: "Run task templates",
	Run: func(cmd *cobra.Command, args []string) {
		flags := cmd.Flags()
		client, err := getClient(flags)
		if err != nil {
			log.Fatalf("Failed to create client: %s", err)
		}
		query, err := getTaskTemplateQuery(flags)
		if err != nil {
			log.Fatalf("Failed to get task query: %s", err)
		}
		templates, err := client.ListTaskTemplates(query)
		if err != nil {
			log.Fatalf("Failed to list task templates: %s", err)
		}
		for _, template := range templates {
			task, err := template.GetTask(nil)
			if err != nil {
				log.Fatalf("Failed to get task: %s", err)
				continue
			}
			result, err := task.Run(cmd.Context())
			if err != nil {
				log.Fatalf("Failed to run task: %s", err)
			}
			fmt.Println(result)
		}
	},
}

func init() {
	rootCmd.AddCommand(taskTemplatesCmd)

	rootCmd.PersistentFlags().StringP("workdir", "w", "", "Work directory")
	rootCmd.MarkPersistentFlagRequired("workdir")

	taskTemplatesCmd.AddCommand(generateTaskTemplatesCmd, listTaskTemplatesCmd, countTaskTemplatesCmd, runTaskTemplatesCmd)

	taskTemplatesCmd.PersistentFlags().StringArrayP("task-template-id", "", []string{}, "Task template ID")
	taskTemplatesCmd.PersistentFlags().StringArrayP("platform", "", []string{}, "Platform (windows, linux, darwin)")
	taskTemplatesCmd.PersistentFlags().BoolP("auto", "", false, "Auto")
	taskTemplatesCmd.PersistentFlags().BoolP("elevation-required", "", false, "Elevation required")
	taskTemplatesCmd.PersistentFlags().StringArrayP("tag", "", nil, "Tags")
	taskTemplatesCmd.PersistentFlags().StringArrayP("attack-technique-id", "", nil, "MITRE ATT&CK technique ID (e.g. T1057)")
}
