package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const (
	configDir  = ".config" // Name of the configuration directory within the user's home directory.
	configFile = "llm-cli/config.json" // Path to the configuration file within the config directory.
)

// Config represents the overall structure of the application's configuration file.
// It holds the name of the currently active profile and a map of all defined profiles.
type Config struct {
	CurrentProfile string             `json:"current_profile"` // The name of the currently active profile.
	Profiles       map[string]Profile `json:"profiles"`        // A map of profile names to their respective configurations.
}

// Profile defines the settings for a specific LLM provider and model.
// It includes various parameters required to interact with different LLM services.
type Profile struct {
	Provider           string `json:"provider"`            // The name of the LLM provider (e.g., "ollama", "openai", "bedrock", "vertexai").
	Endpoint           string `json:"endpoint,omitempty"`        // The API endpoint URL for the LLM service.
	APIKey             string `json:"api_key,omitempty"`         // The API key for authentication with the LLM service.
	Model              string `json:"model"`               // The specific model name to use (e.g., "llama3", "gpt-4", "gemini-1.5-pro-001").
	AWSRegion          string `json:"aws_region,omitempty"`      // AWS region for Bedrock.
	AWSAccessKeyID     string `json:"aws_access_key_id,omitempty"` // AWS Access Key ID for Bedrock.
	AWSSecretAccessKey string `json:"aws_secret_access_key,omitempty"` // AWS Secret Access Key for Bedrock.
	ProjectID          string `json:"project_id,omitempty"`      // GCP Project ID for Vertex AI.
	Location           string `json:"location,omitempty"`        // GCP Location for Vertex AI.
	CredentialsFile    string `json:"credentials_file,omitempty"` // Path to GCP credentials file for Vertex AI (e.g., service account key).
}

// Load reads the configuration file from the user's config directory.
// If the file does not exist, it returns a default configuration.
func Load() (*Config, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			// If the config file does not exist, return a default configuration.
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

// Save writes the current configuration to the user's config directory.
// It creates the directory if it does not exist.
func (c *Config) Save() error {
	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}

	// Ensure the configuration directory exists.
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}

// GetConfigPath returns the absolute path to the configuration file.
// It constructs the path based on the user's home directory.
func GetConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, configDir, configFile), nil
}
