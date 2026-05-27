package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/xiaolonglong/harborctl/internal/client"
)

var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Manage Harbor projects",
	Long:  `Commands for managing Harbor projects.`,
}

var projectListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all projects",
	Long:  `List all Harbor projects with detailed information.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		page, _ := cmd.Flags().GetInt("page")
		pageSize, _ := cmd.Flags().GetInt("page-size")

		projects, err := harborClient.ListProjects(name, page, pageSize)
		if err != nil {
			return err
		}

		if len(projects) == 0 {
			fmt.Println("No projects found.")
			return nil
		}

		printProjects(projects)
		return nil
	},
}

var projectCreateCmd = &cobra.Command{
	Use:   "create <project-name>",
	Short: "Create a new project",
	Long:  `Create a new Harbor project.

Examples:
  harborctl project create my-project
  harborctl project create my-public-project --public`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		isPublic, _ := cmd.Flags().GetBool("public")

		project, err := harborClient.CreateProject(name, isPublic)
		if err != nil {
			return err
		}

		fmt.Printf("Project '%s' created successfully!\n", project.Name)
		fmt.Printf("  ID:     %d\n", project.ProjectID)
		fmt.Printf("  Public: %v\n", project.Public)
		return nil
	},
}

var projectDeleteCmd = &cobra.Command{
	Use:   "delete <project-name>",
	Short: "Delete a project",
	Long: `Delete a Harbor project by name or ID.

Use --force to force delete even if project has artifacts.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		nameOrID := args[0]
		force, _ := cmd.Flags().GetBool("force")

		deletable, msg, err := harborClient.GetProjectDeletable(nameOrID)
		if err != nil && !force {
			return err
		}
		if !deletable && !force {
			fmt.Printf("Project is not deletable: %s\n", msg)
			fmt.Println("Use --force to force delete.")
			return nil
		}

		if err := harborClient.DeleteProject(nameOrID, force); err != nil {
			return err
		}

		fmt.Printf("Project '%s' deleted successfully!\n", nameOrID)
		return nil
	},
}

var projectInspectCmd = &cobra.Command{
	Use:   "inspect <project-name>",
	Short: "Show project details",
	Long:  `Display detailed information about a project.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		nameOrID := args[0]

		project, err := harborClient.GetProject(nameOrID)
		if err != nil {
			return err
		}

		fmt.Printf("Project: %s\n", project.Name)
		fmt.Printf("  ID:       %d\n", project.ProjectID)
		fmt.Printf("  Public:   %v\n", project.Public)
		fmt.Printf("  Owner:    %s (ID: %d)\n", project.OwnerName, project.OwnerID)
		fmt.Printf("  Created:  %s\n", project.CreationTime.Format("2006-01-02 15:04:05"))
		fmt.Printf("  Updated:  %s\n", project.UpdateTime.Format("2006-01-02 15:04:05"))
		return nil
	},
}

func printProjects(projects []client.Project) {
	fmt.Printf("%-5s %-30s %-10s %-19s\n", "ID", "Name", "Type", "Created")
	fmt.Printf("%-5s %-30s %-10s %-19s\n", strings.Repeat("-", 5), strings.Repeat("-", 30), strings.Repeat("-", 10), strings.Repeat("-", 19))

	for _, p := range projects {
		projType := "private"
		if p.Public {
			projType = "public"
		}
		fmt.Printf("%-5d %-30s %-10s %-19s\n", p.ProjectID, p.Name, projType, p.CreationTime.Format("2006-01-02 15:04:05"))
	}
}

func init() {
	rootCmd.AddCommand(projectCmd)

	projectCmd.AddCommand(projectListCmd)
	projectCmd.AddCommand(projectCreateCmd)
	projectCmd.AddCommand(projectDeleteCmd)
	projectCmd.AddCommand(projectInspectCmd)

	projectListCmd.Flags().StringP("name", "n", "", "Filter projects by name")
	projectListCmd.Flags().IntP("page", "", 1, "Page number")
	projectListCmd.Flags().IntP("page-size", "s", 100, "Page size")

	projectCreateCmd.Flags().BoolP("public", "p", false, "Create as public project")

	projectDeleteCmd.Flags().BoolP("force", "f", false, "Force delete project with artifacts")
}