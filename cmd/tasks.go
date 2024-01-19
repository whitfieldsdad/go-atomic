package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"strings"

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

	// Select tasks using ATT&CK tactic or technique IDs (e.g. TA0003, T1087)
	attackTechniques, _ := flags.GetStringArray("attack-technique-id")
	attackTactics, _ := flags.GetStringArray("attack-tactic-id")

	if len(attackTactics) > 0 {
		attackTacticsToTechniquesPath, _ := flags.GetString("attack-tactics-to-techniques-path")
		attackTacticsToTechniques, err := readMapOfAttackTacticsToTechniques(attackTacticsToTechniquesPath)
		if err != nil {
			log.Fatalf("Failed to read map of ATT&CK tactics to techniques: %s", err)
		}
		for _, tactic := range attackTactics {
			techniques, ok := attackTacticsToTechniques[tactic]
			if ok {
				log.Infof("Selected %d techniques related to %s: %s", len(techniques), tactic, strings.Join(techniques, ", "))
				attackTechniques = append(attackTechniques, techniques...)
			}
		}
	}

	query := &atomic.TaskQuery{
		TaskIds:            taskIds,
		AttackTechniqueIds: attackTechniques,
		Platforms:          platforms,
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

func readMapOfAttackTacticsToTechniques(path string) (map[string][]string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	m := make(map[string][]string)
	err = json.Unmarshal(b, &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func init() {
	rootCmd.AddCommand(testsCmd)
	testsCmd.AddCommand(listTasksCmd, countTasksCmd, executeTasksCmd)

	sharedFlags := testsCmd.PersistentFlags()
	sharedFlags.String("atomics-path", atomic.DefaultAtomicsPath, "Path to atomics file/directory")
	sharedFlags.String("attack-tactics-to-techniques-path", "", "Path to ATT&CK tactic to technique lookup table")
	sharedFlags.StringArray("task-id", nil, "Task ID")
	sharedFlags.StringArray("attack-tactic-id", nil, "ATT&CK tactic IDs")
	sharedFlags.StringArray("attack-technique-id", nil, "ATT&CK technique IDs")
	sharedFlags.StringArrayP("platform", "p", nil, "(windows,linux,darwin)")
	sharedFlags.Bool("elevation-required", false, "")
	sharedFlags.Bool("auto", false, "Select tests which can be run by the current user")
}
