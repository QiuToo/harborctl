package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/xiaolonglong/harborctl/internal/client"
)

var memberCmd = &cobra.Command{
	Use:   "member",
	Short: "Manage project members",
	Long:  `Commands for managing Harbor project members.`,
}

var memberListCmd = &cobra.Command{
	Use:   "list <project>",
	Short: "List project members",
	Long:  `List all members of a Harbor project.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectName := args[0]

		members, err := harborClient.ListProjectMembers(projectName)
		if err != nil {
			return err
		}

		if len(members) == 0 {
			fmt.Printf("No members found in project '%s'.\n", projectName)
			return nil
		}

		fmt.Printf("Members of project '%s':\n", projectName)
		fmt.Printf("%-10s %-20s %-10s\n", "ID", "Name", "Role")
		fmt.Printf("%-10s %-20s %-10s\n", strings.Repeat("-", 10), strings.Repeat("-", 20), strings.Repeat("-", 10))

		for _, m := range members {
			fmt.Printf("%-10d %-20s %-10s\n", m.EntityID, m.EntityName, m.RoleName)
		}

		return nil
	},
}

var memberAddCmd = &cobra.Command{
	Use:   "add <project> <username>",
	Short: "Add member to project",
	Long:  `Add a user to a Harbor project.`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectName := args[0]
		username := args[1]

		users, err := harborClient.ListUsers()
		if err != nil {
			return err
		}

		var userID int
		for _, u := range users {
			if u.Username == username {
				userID = u.UserID
				break
			}
		}

		if userID == 0 {
			return fmt.Errorf("user '%s' not found", username)
		}

		if err := harborClient.AddProjectMember(projectName, &client.ProjectMember{EntityID: userID}); err != nil {
			return err
		}

		fmt.Printf("Added user '%s' to project '%s'\n", username, projectName)
		return nil
	},
}

var memberRemoveCmd = &cobra.Command{
	Use:   "remove <project> <username|user-id>",
	Short: "Remove member from project",
	Long:  `Remove a member from a Harbor project.`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectName := args[0]
		identifier := args[1]

		members, err := harborClient.ListProjectMembers(projectName)
		if err != nil {
			return err
		}

		var memberID int
		for _, m := range members {
			if m.EntityName == identifier || fmt.Sprintf("%d", m.EntityID) == identifier {
				memberID = m.EntityID
				break
			}
		}

		if memberID == 0 {
			return fmt.Errorf("member '%s' not found in project '%s'", identifier, projectName)
		}

		if err := harborClient.RemoveProjectMember(projectName, memberID); err != nil {
			return err
		}

		fmt.Printf("Removed member '%s' from project '%s'\n", identifier, projectName)
		return nil
	},
}

var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search projects and repositories",
	Long:  `Search for Harbor projects and repositories by name.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := args[0]

		result, err := harborClient.Search(query)
		if err != nil {
			return err
		}

		fmt.Println("Search Results:")
		fmt.Println()

		if len(result.Projects) > 0 {
			fmt.Println("Projects:")
			fmt.Printf("%-5s %-30s %-10s\n", "ID", "Name", "Public")
			fmt.Printf("%-5s %-30s %-10s\n", strings.Repeat("-", 5), strings.Repeat("-", 30), strings.Repeat("-", 10))

			for _, p := range result.Projects {
				public := "private"
				if p.Public {
					public = "public"
				}
				fmt.Printf("%-5d %-30s %-10s\n", p.ProjectID, p.Name, public)
			}
			fmt.Println()
		}

		if len(result.Repositories) > 0 {
			fmt.Println("Repositories:")
			for _, r := range result.Repositories {
				fmt.Printf("  %s\n", r.Name)
			}
		}

		if len(result.Projects) == 0 && len(result.Repositories) == 0 {
			fmt.Println("No results found.")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(memberCmd)
	rootCmd.AddCommand(searchCmd)

	memberCmd.AddCommand(memberListCmd)
	memberCmd.AddCommand(memberAddCmd)
	memberCmd.AddCommand(memberRemoveCmd)
}