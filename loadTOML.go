package confiq

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml"
)

// LoadTOMLFromFile loads a TOML file from the given path.
func (c *ConfigSet) LoadTOMLFromFile(path string, options ...loadOption) error {
	inputBytes, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return fmt.Errorf("%w: %w", ErrCannotOpenConfig, err)
	}

	return c.LoadTOMLFromBytes(inputBytes, options...)
}

// LoadTOMLFromString loads a TOML file from the given string.
func (c *ConfigSet) LoadTOMLFromString(input string, options ...loadOption) error {
	inputBytes := []byte(input)

	return c.LoadTOMLFromBytes(inputBytes, options...)
}

// LoadTOMLFromReader loads a TOML file from a reader stream.
func (c *ConfigSet) LoadTOMLFromReader(reader io.Reader, options ...loadOption) error {
	bytes, err := readerToBytes(reader)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrCannotLoadConfig, err)
	}

	return c.LoadTOMLFromBytes(bytes, options...)
}

// LoadTOMLFromBytes loads a TOML file from the given bytes.
func (c *ConfigSet) LoadTOMLFromBytes(input []byte, options ...loadOption) error {
	value := new(any)

	err := toml.Unmarshal(input, value)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrCannotLoadConfig, err)
	}

	return c.applyValue(*value, options...)
}
