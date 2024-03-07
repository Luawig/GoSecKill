package config

import (
	"github.com/spf13/viper"
)

var config *AppConfig

// AppConfig holds the application configuration
type AppConfig struct {
	// Application specific configurations

	// LogConfig holds the log configuration
	LogConfig LogConfig `mapstructure:"log"`
}

type LogConfig struct {
	DeveloperMode   bool `mapstructure:"developerMode"`
	OutputToConsole bool `mapstructure:"outputToConsole"`
}

// LoadConfig loads application configuration from file and environment variables.
func LoadConfig(configPath string) error {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(configPath)

	// Read configuration from the file
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	// Override with environment variables
	viper.AutomaticEnv()

	if err := viper.Unmarshal(&config); err != nil {
		return err
	}

	return nil
}
