package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/xiaolonglong/harborctl/internal/client"
)

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show Harbor system information",
	Long:  `Display Harbor system information including version and health status.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		showDetails, _ := cmd.Flags().GetBool("details")

		sysInfo, err := harborClient.GetSystemInfo()
		if err != nil {
			return err
		}

		fmt.Println("╔═══════════════════════════════════════════════════════╗")
		fmt.Println("║           Harbor System Information            ║")
		fmt.Println("╚═══════════════════════════════════════════════════════╝")
		fmt.Printf("Version:          %s\n", sysInfo.HarborVersion)
		fmt.Printf("Database:         %s\n", sysInfo.DatabaseType)
		fmt.Printf("Self-Registration: %v\n", sysInfo.SelfRegistration)
		fmt.Printf("LDAP Enabled:     %v\n", sysInfo.LDAPEnabled)

		if showDetails {
			health, err := harborClient.GetHealth()
			if err == nil {
				fmt.Println("\n--- Component Health Status ---")
				printHealthStatus(health)
			}

			stats, err := harborClient.GetStatistics()
			if err == nil {
				fmt.Println("\n--- Statistics ---")
				fmt.Printf("Projects:      %d\n", stats.ProjectCount)
				fmt.Printf("Repositories:  %d\n", stats.RepoCount)
			}
		}

		return nil
	},
}

var healthCmd = &cobra.Command{
	Use:   "health",
	Short: "Check Harbor health status",
	Long:  `Display health status of all Harbor components.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		health, err := harborClient.GetHealth()
		if err != nil {
			return err
		}

		fmt.Println("Harbor Components Health Status:")
		printHealthStatus(health)
		return nil
	},
}

var statsCmd = &cobra.Command{
	Use:   "stat",
	Short: "Show system statistics",
	Long:  `Display Harbor system statistics.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		stats, err := harborClient.GetStatistics()
		if err != nil {
			return err
		}

		fmt.Println("--- Harbor Statistics ---")
		fmt.Printf("Projects:     %d\n", stats.ProjectCount)
		fmt.Printf("Repositories: %d\n", stats.RepoCount)
		return nil
	},
}

func printHealthStatus(health *client.OverallHealthStatus) {
	// Handle new format with Components array
	if len(health.Components) > 0 {
		for _, c := range health.Components {
			statusIcon := "✓"
			if c.Status != "healthy" {
				statusIcon = "✗"
			}
			fmt.Printf("  %s %s: %s\n", statusIcon, c.Name, c.Status)
		}
		return
	}

	// Fallback to old format
	components := []struct {
		Name   string
		Status string
	}{
		{"Harbor", health.Harbor.Status},
		{"Portal", health.Portal.Status},
		{"Core", health.Core.Status},
		{"Jobservice", health.Jobservice.Status},
		{"Registry", health.Registry.Status},
		{"Database", health.Database.Status},
		{"Redis", health.Redis.Status},
		{"Proxy", health.Proxy.Status},
	}

	for _, c := range components {
		if c.Status == "" {
			continue
		}
		statusIcon := "✓"
		if c.Status != "healthy" {
			statusIcon = "✗"
		}
		fmt.Printf("  %s %s: %s\n", statusIcon, c.Name, c.Status)
	}
}

func init() {
	rootCmd.AddCommand(infoCmd)
	rootCmd.AddCommand(healthCmd)
	rootCmd.AddCommand(statsCmd)

	infoCmd.Flags().BoolP("details", "d", false, "Show detailed information")
}