package jsonprovider

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/arjanvaneersel/kit/cfg"
)

var (
	ErrEmptyFilename error = errors.New("Filename is empty")
)

// Provider provides a map with all configuration settings set in a json file
// The prefix will be removed from the key name. Returns a map or error.
type Provider struct {
	Filename string
}

// Provide implements the Provider interface
func (j Provider) Provide() (map[string]string, error) {
	// Create an empty map to store the configuration
	cfg := map[string]string{}

	// Check if the filename has been set
	if len(j.Filename) == 0 {
		return nil, ErrEmptyFilename
	}

	// Open the file and ensure the file will be properly closed after all operations
	file, err := os.Open(j.Filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	dec := json.NewDecoder(file)
	if err := dec.Decode(&cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Save will store the provided configuration
func (j Provider) Save(cfg *cfg.Config) error {
	// Check if the filename has been set
	if len(j.Filename) == 0 {
		return ErrEmptyFilename
	}

	// Open or create the file and ensure the file will be properly closed after all operations
	file, err := os.Open(j.Filename)
	if err != nil {
		if os.IsNotExist(err) {
			file, err = os.Create(j.Filename)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	defer file.Close()

	// Encode and save the configuration map
	encoder := json.NewEncoder(file)
	if err := encoder.Encode(cfg.Map()); err != nil {
		return err
	}

	return nil
}
