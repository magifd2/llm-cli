package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const (
	configDir  = ".config"
	configFile = "llm-cli/config.json"
)

// Config represents the structure of the configuration file.
	type Config struct {
	CurrentProfile string             `json:"current_profile"`
	Profiles       map[string]Profile `json:"profiles"`
}

// Profile defines the settings for a specific LLM provider and model.
type Profile struct {
	Provider string `json:"provider"`
	Endpoint string `json:"endpoint,omitempty"`
	APIKey    string `json:"api_key,omitempty"`
	Model          string `json:"model"`
	AWSRegion      string `json:"aws_region,omitempty"`
	AWSAccessKeyID string `json:"aws_access_key_id,omitempty"`
	AWSSecretAccessKey string `json:"aws_secret_access_key,omitempty"`
}

// Load reads the configuration file from the user's config directory.
func Load() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configPath := filepath.Join(home, configDir, configFile)
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Return a default config if the file doesn't exist
			return &Config{
				CurrentProfile: "default",
				Profiles: map[string]Profile{
					"default": {
						Provider: "ollama",
						Model:    "llama3",
					},
				},
			}, nil
		}
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// Save writes the configuration to the user's config directory.
func (c *Config) Save() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configPath := filepath.Join(home, configDir, configFile)
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}
