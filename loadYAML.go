package confiq

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// LoadYAMLFromFile loads a YAML file from the given path.
func (c *ConfigSet) LoadYAMLFromFile(path string, options ...loadOption) error {
	inputBytes, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return fmt.Errorf("%w: %w", ErrCannotOpenConfig, err)
	}

	return c.LoadYAMLFromBytes(inputBytes, options...)
}

// LoadYAMLFromString loads a YAML file from the given string.
func (c *ConfigSet) LoadYAMLFromString(input string, options ...loadOption) error {
	inputBytes := []byte(input)

	return c.LoadYAMLFromBytes(inputBytes, options...)
}

// LoadYAMLFromReader loads a YAML file from a reader stream.
func (c *ConfigSet) LoadYAMLFromReader(reader io.Reader, options ...loadOption) error {
	bytes, err := readerToBytes(reader)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrCannotLoadConfig, err)
	}

	return c.LoadYAMLFromBytes(bytes, options...)
}

// LoadYAMLFromBytes loads a YAML file from the given bytes.
func (c *ConfigSet) LoadYAMLFromBytes(input []byte, options ...loadOption) error {
	value := new(any)

	err := yaml.Unmarshal(input, value)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrCannotLoadConfig, err)
	}

	return c.applyValue(*value, options...)
}
