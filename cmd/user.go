package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/xiaolonglong/harborctl/internal/client"
)

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Manage Harbor users",
	Long:  `Commands for managing Harbor users.`,
}

var userListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all users",
	Long:  `List all Harbor users. Only admin users can list users.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		users, err := harborClient.ListUsers()
		if err != nil {
			return err
		}

		if len(users) == 0 {
			fmt.Println("No users found.")
			return nil
		}

		printUsers(users)
		return nil
	},
}

var userCreateCmd = &cobra.Command{
	Use:   "create <username>",
	Short: "Create a new user",
	Long: `Create a new Harbor user.

Requires admin privileges.

Examples:
  harborctl user create john --email john@example.com --password Secret123!
  harborctl user create adminuser --admin`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		username := args[0]
		email, _ := cmd.Flags().GetString("email")
		password, _ := cmd.Flags().GetString("password")
		isAdmin, _ := cmd.Flags().GetBool("admin")
		realname, _ := cmd.Flags().GetString("realname")

		if email == "" {
			return fmt.Errorf("email is required")
		}
		if password == "" {
			return fmt.Errorf("password is required")
		}

		newUser := &client.User{
			Username: username,
			Email:    email,
			Password: password,
			Admin:    isAdmin,
			RealName: realname,
		}

		user, err := harborClient.CreateUser(newUser)
		if err != nil {
			return err
		}

		fmt.Printf("User '%s' created successfully!\n", user.Username)
		fmt.Printf("  ID:     %d\n", user.UserID)
		fmt.Printf("  Email:  %s\n", user.Email)
		fmt.Printf("  Admin:  %v\n", user.Admin)
		return nil
	},
}

var userDeleteCmd = &cobra.Command{
	Use:   "delete <username|user-id>",
	Short: "Delete a user",
	Long:  `Delete a Harbor user by username or ID.

Only admin users can delete users. Cannot delete yourself.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		identifier := args[0]

		userID, err := findUserID(identifier)
		if err != nil {
			return err
		}

		if err := harborClient.DeleteUser(userID); err != nil {
			return err
		}

		fmt.Printf("User '%s' deleted successfully!\n", identifier)
		return nil
	},
}

var userInspectCmd = &cobra.Command{
	Use:   "inspect <username|user-id>",
	Short: "Show user details",
	Long:  `Display detailed information about a user.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		identifier := args[0]

		user, err := findUser(identifier)
		if err != nil {
			return err
		}

		fmt.Printf("User: %s\n", user.Username)
		fmt.Printf("  ID:       %d\n", user.UserID)
		fmt.Printf("  Email:    %s\n", user.Email)
		fmt.Printf("  RealName: %s\n", user.RealName)
		fmt.Printf("  Admin:    %v\n", user.Admin)
		return nil
	},
}

func findUser(identifier string) (*client.User, error) {
	if id, err := strconv.Atoi(identifier); err == nil {
		return harborClient.GetUser(id)
	}

	users, err := harborClient.ListUsers()
	if err != nil {
		return nil, err
	}

	for _, u := range users {
		if u.Username == identifier {
			return &u, nil
		}
	}

	return nil, fmt.Errorf("user '%s' not found", identifier)
}

func findUserID(identifier string) (int, error) {
	if id, err := strconv.Atoi(identifier); err == nil {
		return id, nil
	}

	users, err := harborClient.ListUsers()
	if err != nil {
		return 0, err
	}

	for _, u := range users {
		if u.Username == identifier {
			return u.UserID, nil
		}
	}

	return 0, fmt.Errorf("user '%s' not found", identifier)
}

func printUsers(users []client.User) {
	fmt.Printf("%-5s %-20s %-25s %-10s\n", "ID", "Username", "Email", "Admin")
	fmt.Printf("%-5s %-20s %-25s %-10s\n", strings.Repeat("-", 5), strings.Repeat("-", 20), strings.Repeat("-", 25), strings.Repeat("-", 10))

	for _, u := range users {
		adminStr := "false"
		if u.Admin {
			adminStr = "true"
		}
		fmt.Printf("%-5d %-20s %-25s %-10s\n", u.UserID, u.Username, u.Email, adminStr)
	}
}

func init() {
	rootCmd.AddCommand(userCmd)

	userCmd.AddCommand(userListCmd)
	userCmd.AddCommand(userCreateCmd)
	userCmd.AddCommand(userDeleteCmd)
	userCmd.AddCommand(userInspectCmd)

	userCreateCmd.Flags().StringP("email", "e", "", "User email (required)")
	userCreateCmd.Flags().StringP("password", "p", "", "User password (required)")
	userCreateCmd.Flags().StringP("realname", "r", "", "User real name")
	userCreateCmd.Flags().BoolP("admin", "a", false, "Grant admin privileges")
}