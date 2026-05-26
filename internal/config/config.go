package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// Config holds the Harbor connection configuration
type Config struct {
	Address  string `mapstructure:"address" json:"address"`
	Username string `mapstructure:"username" json:"username"`
	Password string `mapstructure:"password" json:"password"`
	Scheme  string `mapstructure:"scheme" json:"scheme"`
	Insecure bool  `mapstructure:"insecure" json:"insecure"`
}

// Load loads configuration from file
func Load(configPath string) (*Config, error) {
	configPath = getConfigPath(configPath)

	viper.SetConfigType("yaml")
	viper.SetConfigFile(configPath)

	viper.SetDefault("scheme", "http")
	viper.SetDefault("insecure", false)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config: %w", err)
		}
	}

	cfg := &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return cfg, nil
}

func getConfigPath(configPath string) string {
	if configPath == "" {
		defaultPath := "/etc/harbor/harbor.yaml"
		if _, err := os.Stat(defaultPath); err == nil {
			return defaultPath
		}
		return "./harbor.yaml"
	}
	return configPath
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.Address == "" {
		return fmt.Errorf("harbor address is required")
	}
	if c.Username == "" {
		return fmt.Errorf("harbor username is required")
	}
	if c.Password == "" {
		return fmt.Errorf("harbor password is required")
	}
	return nil
}

// GetBaseURL returns base URL
func (c *Config) GetBaseURL() string {
	scheme := c.Scheme
	if scheme == "" {
		scheme = "http"
	}
	return fmt.Sprintf("%s://%s/api/v2.0", scheme, c.Address)
}