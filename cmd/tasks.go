package cmd

import (
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/whitfieldsdad/go-atomic-red-team/pkg/atomic"
)

var taskTemplatesCmd = &cobra.Command{
	Use:   "task-templates",
	Short: "Task templates",
}

var generateTaskTemplatesCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate task templates",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		flags := cmd.Flags()
		client, err := getClient(flags)
		if err != nil {
			log.Fatalf("Failed to create client: %s", err)
		}
		i, _ := flags.GetString("input-path")
		o, _ := flags.GetString("output-path")
		err = client.GenerateTaskTemplates(i, o)
		if err != nil {
			log.Fatalf("Failed to generate task templates: %s", err)
		}
	},
}

func getClient(flags *pflag.FlagSet) (*atomic.Client, error) {
	client, err := atomic.NewClient()
	if err != nil {
		return nil, err
	}
	return client, nil
}

func init() {
	rootCmd.AddCommand(taskTemplatesCmd)
	taskTemplatesCmd.AddCommand(generateTaskTemplatesCmd)

	generateTaskTemplatesCmd.Flags().StringP("input-path", "i", "", "Input file/directory")
	generateTaskTemplatesCmd.Flags().StringP("output-path", "o", "", "Output directory")
	generateTaskTemplatesCmd.MarkFlagRequired("input-path")
	generateTaskTemplatesCmd.MarkFlagRequired("output-path")
}
