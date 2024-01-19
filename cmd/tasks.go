package cmd

import (
	"fmt"
	"runtime"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/whitfieldsdad/go-atomic-red-team/pkg/atomic"
	"golang.org/x/exp/slices"
)

var testsCmd = &cobra.Command{
	Use:   "tasks",
	Short: "Tasks",
}

var listTasksCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l", "all"},
	Short:   "List tasks",
	Run: func(cmd *cobra.Command, args []string) {
		flags := cmd.Flags()
		tasks := listTasks(flags)
		for _, task := range tasks {
			printJson(task)
		}
	},
}

var countTasksCmd = &cobra.Command{
	Use:     "count",
	Aliases: []string{"c", "n", "total"},
	Short:   "Count tasks",
	Run: func(cmd *cobra.Command, args []string) {
		flags := cmd.Flags()
		tasks := listTasks(flags)
		total := len(tasks)
		fmt.Println(total)
	},
}

var executeTasksCmd = &cobra.Command{
	Use:     "exec",
	Aliases: []string{"x", "execute", "run"},
	Short:   "Execute tasks",
	Run: func(cmd *cobra.Command, args []string) {
		flags := cmd.Flags()
		tasks := listTasks(flags)

		ctx := cmd.Context()
		for _, task := range tasks {
			result := task.Exec(ctx)
			printJson(result)
		}
	},
}

func listTasks(flags *pflag.FlagSet) []atomic.Task {
	atomicsPath, _ := flags.GetString("atomics-path")
	if atomicsPath == "" {
		log.Fatalf("No atomics path provided - set it with $ATOMICS_PATH or --atomics-path")
	}
	q := getTaskQuery(flags)
	client, err := atomic.NewClient(atomicsPath)
	if err != nil {
		log.Fatalf("Failed to create client: %s", err)
	}
	tasks, err := client.ListTasks(q)
	if err != nil {
		log.Fatalf("Failed to list tasks: %s", err)
	}
	return tasks
}

func getTaskQuery(flags *pflag.FlagSet) *atomic.TaskQuery {
	taskIds, _ := flags.GetStringArray("task-id")
	platforms, _ := flags.GetStringArray("platform")
	tags, _ := flags.GetStringArray("tag")

	query := &atomic.TaskQuery{
		TaskIds:   taskIds,
		Platforms: platforms,
		Tags:      tags,
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
	return query
}

func init() {
	rootCmd.AddCommand(testsCmd)
	testsCmd.AddCommand(listTasksCmd, countTasksCmd, executeTasksCmd)

	sharedFlags := testsCmd.PersistentFlags()
	sharedFlags.StringP("atomics-path", "i", atomic.DefaultAtomicsPath, "Path to atomics file/directory")

	sharedFlags.StringArray("task-id", nil, "Task ID")
	sharedFlags.StringArrayP("platform", "p", nil, "(windows,linux,darwin)")
	sharedFlags.StringArrayP("tag", "t", nil, "")
	sharedFlags.Bool("elevation-required", false, "")
	sharedFlags.Bool("auto", false, "Select tests which can be run by the current user")
}
