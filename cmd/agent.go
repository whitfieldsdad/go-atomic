package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/whitfieldsdad/go-atomic-red-team/pkg/atomic"
)

var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "Agent commands",
}

var runAgentCmd = &cobra.Command{
	Use:   "run",
	Short: "Run agent",
	Run: func(cmd *cobra.Command, args []string) {
		flags := cmd.Flags()
		workdir, _ := flags.GetString("workdir")
		agent := atomic.NewAgent(workdir)

		err := agent.Run(cmd.Context())
		if err != nil {
			log.Fatalf("Failed to run agent: %v", err)
		}
	},
}

func init() {
	agentCmd.AddCommand(runAgentCmd)
	rootCmd.AddCommand(agentCmd)

	runAgentCmd.Flags().StringP("workdir", "w", "", "Working directory")
	_ = runAgentCmd.MarkFlagRequired("workdir")
}
