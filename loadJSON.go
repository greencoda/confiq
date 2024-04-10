package confiq

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// LoadJSONFromFile loads a JSON file from the given path.
func (c *ConfigSet) LoadJSONFromFile(path string, options ...loadOption) error {
	bytes, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return fmt.Errorf("%w: %w", ErrCannotOpenConfig, err)
	}

	return c.LoadJSONFromBytes(bytes, options...)
}

// LoadJSONFromString loads a JSON file from the given string.
func (c *ConfigSet) LoadJSONFromString(input string, options ...loadOption) error {
	bytes := []byte(input)

	return c.LoadJSONFromBytes(bytes, options...)
}

// LoadJSONFromReader loads a JSON file from a reader stream.
func (c *ConfigSet) LoadJSONFromReader(reader io.Reader, options ...loadOption) error {
	bytes, err := readerToBytes(reader)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrCannotLoadConfig, err)
	}

	return c.LoadJSONFromBytes(bytes, options...)
}

// LoadJSONFromBytes loads a JSON file from the given bytes.
func (c *ConfigSet) LoadJSONFromBytes(input []byte, options ...loadOption) error {
	value := new(any)

	err := json.Unmarshal(input, value)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrCannotLoadConfig, err)
	}

	return c.applyValue(*value, options...)
}
