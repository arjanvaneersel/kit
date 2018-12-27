package envprovider

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

var (
	ErrNoEnvironmentVariables error = errors.New("No environment variables found")
	ErrPrefixNotFound         error = errors.New("Prefix not found")
)

// Provider provides a map with all configuration settings set in the OS environment with the provided prefix
// The prefix will be removed from the key name. Returns a map or error.
type Provider struct {
	Prefix string
}

// Provide implements the Provider interface
func (e Provider) Provide() (map[string]string, error) {
	// Create an empty map to store the configuration
	cfg := map[string]string{}

	// Get the environment variables
	env := os.Environ()
	if len(env) == 0 {
		return nil, ErrNoEnvironmentVariables
	}

	// Convert the EnvProvider's prefix to upper case
	prefix := fmt.Sprintf("%s_", strings.ToUpper(e.Prefix))

	// Loop over each environment variable and store the values matching the prefix. Strip the prefix from the key
	// name when storing the key/value pair in the map
	for _, v := range env {
		if strings.HasPrefix(v, prefix) {
			pair := strings.Split(v, "=")
			cfg[strings.TrimPrefix(pair[0], prefix)] = pair[1]
		}
	}

	// Check if we found anything
	if len(cfg) == 0 {
		return nil, ErrPrefixNotFound
	}
	return cfg, nil
}
