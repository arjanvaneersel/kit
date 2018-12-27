package txtprovider

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/arjanvaneersel/kit/cfg"
)

var (
	ErrEmptyFilename error = errors.New("Filename is empty")
)

// TxtProvider provides a map with all configuration settings set in a text file
// The prefix will be removed from the key name. Returns a map or error.
type TxtProvider struct {
	Filename  string
	Delimiter string
}

// Provide implements the Provider interface
func (t TxtProvider) Provide() (map[string]string, error) {
	// Create an empty map to store the configuration
	cfg := map[string]string{}

	// Check if the filename has been set
	if len(t.Filename) == 0 {
		return nil, ErrEmptyFilename
	}

	// Check if a delimiter has been set, if not set the default delimiter
	if len(t.Delimiter) == 0 {
		t.Delimiter = "="
	}

	// Open the file and ensure the file will be properly closed after all operations
	file, err := os.Open(t.Filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Create a scanner to loop over the lines of the file
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		v := scanner.Text()

		// Skip blank and ambigious lines
		if len(v) < 3 {
			continue
		}

		pair := strings.Split(v, "=")
		cfg[pair[0]] = pair[1]
	}

	return cfg, nil
}

// Save will store the provided configuration
func (t TxtProvider) Save(cfg *cfg.Config) error {
	// Check if the filename has been set
	if len(t.Filename) == 0 {
		return ErrEmptyFilename
	}

	// Check if a delimiter has been set, if not set the default delimiter
	if len(t.Delimiter) == 0 {
		t.Delimiter = "="
	}

	// Open or create the file and ensure the file will be properly closed after all operations
	file, err := os.Open(t.Filename)
	if err != nil {
		if os.IsNotExist(err) {
			file, err = os.Create(t.Filename)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	defer file.Close()

	// Get and loop over the map with configuration values and write each key/value pair in the file
	m := cfg.Map()
	for k, v := range m {
		file.WriteString(fmt.Sprintf("%s=%s\n", k, v))
	}

	return nil
}
