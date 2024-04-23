// Package confiqyaml allows confiq values to be loaded from YAML format.
package confiqyaml

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

var (
	ErrCannotOpenYAMLFile  = errors.New("cannot open YAML file")
	ErrCannotReadYAMLData  = errors.New("cannot read YAML data")
	ErrCannotReadYAMLBytes = errors.New("cannot read YAML bytes")
)

// Container is a struct that holds the loaded values.
type Container struct {
	values []any
	errors []error
}

// Get returns the loaded YAML values.
func (c *Container) Get() []any {
	return c.values
}

// Errors returns the errors that occurred during the loading process.
func (c *Container) Errors() []error {
	return c.errors
}

// Load creates an empty container, into which the YAML values can be loaded.
func Load() *Container {
	container := new(Container)

	return container
}

// FromFile loads a YAML file from the given path.
func (c *Container) FromFile(path string) *Container {
	bytes, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		c.errors = append(c.errors, fmt.Errorf("%w: %w", ErrCannotOpenYAMLFile, err))

		return c
	}

	c.readFromBytes(bytes)

	return c
}

// FromString loads a YAML file from the given string.
func (c *Container) FromString(input string) *Container {
	c.readFromBytes([]byte(input))

	return c
}

// FromReader loads a YAML file from a reader stream.
func (c *Container) FromReader(reader io.Reader) *Container {
	if reader == nil {
		c.errors = append(c.errors, ErrCannotReadYAMLData)

		return c
	}

	buffer := new(bytes.Buffer)

	if _, err := buffer.ReadFrom(reader); err != nil {
		c.errors = append(c.errors, fmt.Errorf("%w: %w", ErrCannotReadYAMLData, err))

		return c
	}

	c.readFromBytes(buffer.Bytes())

	return c
}

// FromBytes loads a YAML file from the given bytes.
func (c *Container) FromBytes(input []byte) *Container {
	c.readFromBytes(input)

	return c
}

func (c *Container) readFromBytes(input []byte) {
	var value any

	if err := yaml.Unmarshal(input, &value); err != nil {
		c.errors = append(c.errors, fmt.Errorf("%w: %w", ErrCannotReadYAMLBytes, err))

		return
	}

	c.values = append(c.values, value)
}
