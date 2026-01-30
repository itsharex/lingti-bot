package cmd

import (
	"fmt"

	"github.com/pltanton/lingti-bot/internal/mcp"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("lingti-bot %s\n", mcp.ServerVersion)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
