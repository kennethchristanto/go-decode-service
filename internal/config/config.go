package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Port         string `yaml:"port"`
		ReadTimeout  int    `yaml:"read_timeout"`
		WriteTimeout int    `yaml:"write_timeout"`
	} `yaml:"server"`

	Security struct {
		SecretKey      string `yaml:"secret_key"`
		EnableCORS     bool   `yaml:"enable_cors"`
		MaxRequestSize int64  `yaml:"max_request_size"`
	} `yaml:"security"`
}

func LoadConfig(configPath string) (*Config, error) {
	viper.SetConfigFile(configPath)
	viper.AutomaticEnv()

	// Set defaults
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.read_timeout", 30)
	viper.SetDefault("server.write_timeout", 30)
	viper.SetDefault("security.enable_cors", true)
	viper.SetDefault("security.max_request_size", 1048576) // 1MB

	// Try to read config file
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

func GetConfig() *Config {
	config, err := LoadConfig("configs/config.yaml")
	if err != nil {
		// Use defaults if config file doesn't exist
		return &Config{
			Server: struct {
				Port         string `yaml:"port"`
				ReadTimeout  int    `yaml:"read_timeout"`
				WriteTimeout int    `yaml:"write_timeout"`
			}{
				Port:         "8080",
				ReadTimeout:  30,
				WriteTimeout: 30,
			},
			Security: struct {
				SecretKey      string `yaml:"secret_key"`
				EnableCORS     bool   `yaml:"enable_cors"`
				MaxRequestSize int64  `yaml:"max_request_size"`
			}{
				SecretKey:      "secretkey",
				EnableCORS:     true,
				MaxRequestSize: 1048576,
			},
		}
	}
	return config
}
