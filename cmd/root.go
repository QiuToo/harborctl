package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/xiaolonglong/harborctl/internal/client"
	"github.com/xiaolonglong/harborctl/internal/config"
)

var (
	cfgFile string
	addr   string
	user   string
	pass   string
)

var rootCmd = &cobra.Command{
	Use:   "harborctl",
	Short: "Harbor CLI - Command-line tool for managing Harbor Registry",
	Long: `Harbor CLI is a command-line tool for managing Harbor Registry.

Examples:
  harborctl -config /etc/harbor/harbor.yaml project list
  harborctl -addr http://192.168.2.222:80 -u admin -p Harbor123456 info
  harborctl project create my-project --public`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Skip initialization for version command
		if cmd.Name() == "version" {
			return nil
		}
		return initClient()
	},
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file path")
	rootCmd.PersistentFlags().StringVarP(&addr, "addr", "a", "", "Harbor server address")
	rootCmd.PersistentFlags().StringVarP(&user, "user", "u", "", "Harbor username")
	rootCmd.PersistentFlags().StringVarP(&pass, "pass", "p", "", "Harbor password")

	viper.BindPFlag("address", rootCmd.Flags().Lookup("addr"))
	viper.BindPFlag("username", rootCmd.Flags().Lookup("user"))
	viper.BindPFlag("password", rootCmd.Flags().Lookup("pass"))
}

const defaultConfigPath = "/etc/harbor/harbor.yaml"

func initConfig() {
	// Command-line flags have higher priority, so set them first
	if addr != "" {
		viper.Set("address", addr)
	}
	if user != "" {
		viper.Set("username", user)
	}
	if pass != "" {
		viper.Set("password", pass)
	}

	// If config file is specified, read from it
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
		if err := viper.ReadInConfig(); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to read config: %v\n", err)
		}
	} else {
		// Try default config path if no CLI flags provided
		if addr == "" && user == "" && pass == "" {
			if _, err := os.Stat(defaultConfigPath); err == nil {
				viper.SetConfigFile(defaultConfigPath)
				if err := viper.ReadInConfig(); err != nil {
					fmt.Fprintf(os.Stderr, "Warning: failed to read default config: %v\n", err)
				}
			}
		}
	}
}

var harborClient *client.Client

func initClient() error {
	cfg := &config.Config{
		Address:  viper.GetString("address"),
		Username: viper.GetString("username"),
		Password: viper.GetString("password"),
		Scheme:   viper.GetString("scheme"),
		Insecure: viper.GetBool("insecure"),
	}

	var err error
	harborClient, err = client.NewClient(cfg)
	return err
}

func Execute() error {
	return rootCmd.Execute()
}