package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var gcCmd = &cobra.Command{
	Use:   "gc",
	Short: "Garbage collection management",
	Long:  `Commands for managing Harbor garbage collection.`,
}

var gcListCmd = &cobra.Command{
	Use:   "list",
	Short: "List GC execution history",
	Long:  `List all garbage collection execution history.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		jobs, err := harborClient.GetGCJobs()
		if err != nil {
			return err
		}

		if len(jobs) == 0 {
			fmt.Println("No GC jobs found.")
			return nil
		}

		fmt.Printf("%-5s %-15s %-15s %-19s %s\n", "ID", "Kind", "Status", "Created", "Details")
		fmt.Printf("%-5s %-15s %-15s %-19s %s\n", strings.Repeat("-", 5), strings.Repeat("-", 15), strings.Repeat("-", 15), strings.Repeat("-", 19), strings.Repeat("-", 20))

		for _, j := range jobs {
			fmt.Printf("%-5d %-15s %-15s %-19s %s\n", j.JobID, j.JobKind, j.JobStatus, j.CreationTime.Format("2006-01-02 15:04:05"), j.JobDetail)
		}

		return nil
	},
}

var gcRunCmd = &cobra.Command{
	Use:   "run [--dry-run]",
	Short: "Trigger garbage collection",
	Long:  `Trigger a garbage collection job.
	
Use --dry-run to test without actually deleting any images.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		dryRun, _ := cmd.Flags().GetBool("dry-run")

		job, err := harborClient.TriggerGC(dryRun)
		if err != nil {
			return err
		}

		fmt.Printf("GC Job triggered successfully!\n")
		fmt.Printf("  ID:     %d\n", job.JobID)
		fmt.Printf("  Status: %s\n", job.JobStatus)

		if dryRun {
			fmt.Println("\nNote: This was a DRY RUN. No images were actually deleted.")
		}

		return nil
	},
}

var gcScheduleCmd = &cobra.Command{
	Use:   "schedule",
	Short: "Show GC schedule",
	Long:  `Display current GC schedule configuration.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		schedule, err := harborClient.GetGCSchedule()
		if err != nil {
			return err
		}

		fmt.Println("GC Schedule:")
		fmt.Printf("  Type: %s\n", schedule.ScheduleType)
		fmt.Printf("  Cron: %s\n", schedule.Cron)

		if schedule.Parameters != nil {
			fmt.Printf("  Delete Untagged: %v\n", schedule.Parameters.DeleteUntagged)
			if schedule.Parameters.AgeDays > 0 {
				fmt.Printf("  Age Days: %d\n", schedule.Parameters.AgeDays)
			}
		}

		return nil
	},
}

var gcScheduleUpdateCmd = &cobra.Command{
	Use:   "schedule update --cron <cron> [--delete-untagged]",
	Short: "Update GC schedule",
	Long:  `Update GC schedule configuration.
	
Examples:
  harborctl gc schedule update --cron "0 2 * * *"           # Daily at 2am
  harborctl gc schedule update --cron "0 0 * * 0"           # Weekly on Sunday
  harborctl gc schedule update --cron "0 2 * * *" --delete-untagged`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cron, _ := cmd.Flags().GetString("cron")
		deleteUntagged, _ := cmd.Flags().GetBool("delete-untagged")

		if cron == "" {
			return fmt.Errorf("--cron is required")
		}

		// Determine schedule type based on cron
		scheduleType := "Manual"
		if strings.Contains(cron, "* * * * *") {
			scheduleType = "Hourly"
		} else if strings.Count(cron, "*") >= 4 {
			scheduleType = "Daily"
		} else if strings.HasPrefix(cron, "0 0") {
			scheduleType = "Weekly"
		}

		err := harborClient.UpdateGCSchedule(scheduleType, cron, deleteUntagged)
		if err != nil {
			return err
		}

		fmt.Printf("GC Schedule updated successfully!\n")
		fmt.Printf("  Type: %s\n", scheduleType)
		fmt.Printf("  Cron: %s\n", cron)
		fmt.Printf("  Delete Untagged: %v\n", deleteUntagged)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(gcCmd)

	gcCmd.AddCommand(gcListCmd)
	gcCmd.AddCommand(gcRunCmd)
	gcCmd.AddCommand(gcScheduleCmd)
	gcCmd.AddCommand(gcScheduleUpdateCmd)

	gcRunCmd.Flags().BoolP("dry-run", "d", false, "Dry run mode (don't actually delete)")

	gcScheduleUpdateCmd.Flags().StringP("cron", "c", "", "Cron expression (required)")
	gcScheduleUpdateCmd.Flags().BoolP("delete-untagged", "", false, "Delete untagged images")
}