package config

import (
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

const (
	configFileEnvKey  = "CONFIG_FILE"
	defaultConfigFile = "config.yml"
)

// Load loads the service configuration file. First it tries to load the file based on the CONFIG_FILE environment
// variable, if not set it defaults to `config.yml`.
func Load() (Service, error) {
	// Check the environment variable.
	configFile, ok := os.LookupEnv(configFileEnvKey)
	if !ok {
		// Use the default value if the environment value is not found.
		configFile = defaultConfigFile
	}

	// Open the file.
	file, err := os.Open(configFile)
	if err != nil {
		return Service{}, fmt.Errorf("config: failed to open file: %w", err)
	}
	defer file.Close()

	// Read the file data.
	data, err := io.ReadAll(file)
	if err != nil {
		return Service{}, fmt.Errorf("config: failed to read file: %w", err)
	}

	// Unmarshal the configuration.
	var config Service
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return Service{}, fmt.Errorf("config: failed to unmarshal file: %w", err)
	}

	return config, nil
}
