package cmd

import (
	"encoding/json"
	"fmt"
	"runtime"
	"slices"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/whitfieldsdad/go-atomic-red-team/pkg/atomic"
)

var tasksCmd = &cobra.Command{
	Use:   "tasks",
	Short: "Tasks",
}

var runTasksCmd = &cobra.Command{
	Use:   "run",
	Short: "Run tasks",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		client, err := getClient(cmd.Flags())
		if err != nil {
			log.Fatalf("Failed to create client: %s", err)
		}
		if len(args) > 0 {
			tasks, err := client.ReadTasks(args)
			if err != nil {
				log.Fatalf("Failed to read tasks: %s", err)
			}
			for _, task := range tasks {
				result := task.Run(ctx)
				blob, err := json.Marshal(result)
				if err != nil {
					log.Fatalf("Failed to marshal task result: %s", err)
				}
				fmt.Println(string(blob))
			}
		}
	},
}

var listTasksCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List tasks",
	Run: func(cmd *cobra.Command, args []string) {
		flags := cmd.Flags()
		client, err := getClient(flags)
		if err != nil {
			log.Fatalf("Failed to create client: %s", err)
		}
		query, err := getTaskQuery(flags)
		if err != nil {
			log.Fatalf("Failed to get task query: %s", err)
		}
		tasks, err := client.ListTasks(query)
		if err != nil {
			log.Fatalf("Failed to list tasks: %s", err)
		}
		for _, task := range tasks {
			blob, err := json.Marshal(task)
			if err != nil {
				log.Fatalf("Failed to marshal task: %s", err)
			}
			fmt.Println(string(blob))
		}
	},
}

var countTasksCmd = &cobra.Command{
	Use:   "count",
	Short: "Count tasks",
	Run: func(cmd *cobra.Command, args []string) {
		flags := cmd.Flags()
		client, err := getClient(flags)
		if err != nil {
			log.Fatalf("Failed to create client: %s", err)
		}
		query, err := getTaskQuery(flags)
		if err != nil {
			log.Fatalf("Failed to get task query: %s", err)
		}
		tasks, err := client.ListTasks(query)
		if err != nil {
			log.Fatalf("Failed to list tasks: %s", err)
		}
		total := len(tasks)
		fmt.Println(total)
	},
}

var generateTasksCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate tasks",
	Run: func(cmd *cobra.Command, args []string) {
		flags := cmd.Flags()
		workdir, _ := flags.GetString("workdir")
		if workdir == "" {
			log.Fatalf("No workdir provided - set it with --workdir")
		}
		client, err := atomic.NewClient(workdir)
		if err != nil {
			log.Fatalf("Failed to create client: %s", err)
		}
		err = client.GenerateTasks()
		if err != nil {
			log.Fatalf("Failed to generate tasks: %s", err)
		}
	},
}

var taskTemplatesCmd = &cobra.Command{
	Use:   "task-templates",
	Short: "Task templates",
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
		query, err := getTaskQuery(flags)
		if err != nil {
			log.Fatalf("Failed to get task query: %s", err)
		}
		templates, err := client.ListTaskTemplates(query)
		if err != nil {
			log.Fatalf("Failed to list task templates: %s", err)
		}
		for _, template := range templates {
			blob, err := json.Marshal(template)
			if err != nil {
				log.Fatalf("Failed to marshal task template: %s", err)
			}
			fmt.Println(string(blob))
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
		query, err := getTaskQuery(flags)
		if err != nil {
			log.Fatalf("Failed to get task query: %s", err)
		}
		templates, err := client.ListTaskTemplates(query)
		if err != nil {
			log.Fatalf("Failed to list task templates: %s", err)
		}
		total := len(templates)
		fmt.Println(total)
	},
}

var generateTaskTemplatesCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate task templates",
	Run: func(cmd *cobra.Command, args []string) {
		flags := cmd.Flags()
		client, err := getClient(flags)
		if err != nil {
			log.Fatalf("Failed to create client: %s", err)
		}
		atomicsPath, _ := flags.GetString("atomics-path")
		if atomicsPath == "" {
			log.Fatalf("No atomics path provided - set it with $ATOMICS_PATH or --atomics-path")
		}
		err = client.GenerateTaskTemplates(atomicsPath)
		if err != nil {
			log.Fatalf("Failed to generate tasks templates: %s", err)
		}
	},
}

func getClient(flags *pflag.FlagSet) (*atomic.Client, error) {
	workdir, _ := flags.GetString("workdir")
	client, err := atomic.NewClient(workdir)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func getTaskQuery(flags *pflag.FlagSet) (*atomic.TaskQuery, error) {
	taskIds, _ := flags.GetStringArray("task-id")
	taskTemplateIds, _ := flags.GetStringArray("task-template-id")
	platforms, _ := flags.GetStringArray("platform")

	query := &atomic.TaskQuery{
		TaskIds:         taskIds,
		TaskTemplateIds: taskTemplateIds,
		Platforms:       platforms,
		Tags:            nil,
	}
	auto, _ := flags.GetBool("auto")
	if auto {
		if !slices.Contains(query.Platforms, runtime.GOOS) {
			query.Platforms = append(query.Platforms, runtime.GOOS)
		}
		elevated := atomic.IsElevated()
		if !elevated {
			elevationRequired := false
			query.ElevationRequired = &elevationRequired
		}
	}
	if flags.Changed("elevation-required") {
		elevationRequired, _ := flags.GetBool("elevation-required")
		query.ElevationRequired = &elevationRequired
	}
	return query, nil
}

func init() {
	rootCmd.PersistentFlags().String("workdir", "", "Path to working directory")

	rootCmd.AddCommand(tasksCmd)
	tasksCmd.AddCommand(generateTasksCmd, countTasksCmd, listTasksCmd, runTasksCmd)

	taskFlags := tasksCmd.PersistentFlags()
	taskFlags.StringArray("task-id", nil, "Task ID")
	taskFlags.StringArray("task-template-id", nil, "Task template ID")
	taskFlags.StringArrayP("platform", "p", nil, "(windows,linux,darwin)")
	taskFlags.Bool("elevation-required", false, "")
	taskFlags.Bool("auto", false, "Select tests which can be run by the current user")

	rootCmd.AddCommand(taskTemplatesCmd)
	taskTemplatesCmd.AddCommand(generateTaskTemplatesCmd, countTaskTemplatesCmd, listTaskTemplatesCmd)

	generateTasksCmd.PersistentFlags().String("atomics-path", atomic.DefaultAtomicsPath, "Path to atomics file/directory")
	generateTaskTemplatesCmd.PersistentFlags().String("atomics-path", atomic.DefaultAtomicsPath, "Path to atomics file/directory")

	taskTemplateFlags := taskTemplatesCmd.PersistentFlags()
	taskTemplateFlags.StringArray("task-id", nil, "Task ID")
	taskTemplateFlags.StringArray("task-template-id", nil, "Task template ID")
	taskTemplateFlags.StringArrayP("platform", "p", nil, "(windows,linux,darwin)")
	taskTemplateFlags.Bool("elevation-required", false, "")
	taskTemplateFlags.Bool("auto", false, "Select tests which can be run by the current user")
}
