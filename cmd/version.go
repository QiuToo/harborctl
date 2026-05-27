package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version = "v1.1.2"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show harborctl version",
	Long:  `Display the version of harborctl.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("harborctl " + version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}