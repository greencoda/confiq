// Package confiqjson allows confiq values to be loaded from JSON format.
package confiqjson

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

var (
	ErrCannotOpenJSONFile  = errors.New("cannot open JSON file")
	ErrCannotReadJSONData  = errors.New("cannot read JSON data")
	ErrCannotReadJSONBytes = errors.New("cannot read JSON bytes")
)

// Container is a struct that holds the loaded values.
type Container struct {
	values []any
	errors []error
}

// Get returns the loaded JSON values.
func (c *Container) Get() []any {
	return c.values
}

// Errors returns the errors that occurred during the loading process.
func (c *Container) Errors() []error {
	return c.errors
}

// Load creates an empty container, into which the JSON values can be loaded.
func Load() *Container {
	container := new(Container)

	return container
}

// FromFile loads a JSON file from the given path.
func (c *Container) FromFile(path string) *Container {
	bytes, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		c.errors = append(c.errors, fmt.Errorf("%w: %w", ErrCannotOpenJSONFile, err))

		return c
	}

	c.readFromBytes(bytes)

	return c
}

// FromString loads a JSON file from the given string.
func (c *Container) FromString(input string) *Container {
	c.readFromBytes([]byte(input))

	return c
}

// FromReader loads a JSON file from a reader stream.
func (c *Container) FromReader(reader io.Reader) *Container {
	if reader == nil {
		c.errors = append(c.errors, ErrCannotReadJSONData)

		return c
	}

	buffer := new(bytes.Buffer)

	if _, err := buffer.ReadFrom(reader); err != nil {
		c.errors = append(c.errors, fmt.Errorf("%w: %w", ErrCannotReadJSONData, err))

		return c
	}

	c.readFromBytes(buffer.Bytes())

	return c
}

// FromBytes loads a JSON file from the given bytes.
func (c *Container) FromBytes(input []byte) *Container {
	c.readFromBytes(input)

	return c
}

func (c *Container) readFromBytes(input []byte) {
	var value any

	if err := json.Unmarshal(input, &value); err != nil {
		c.errors = append(c.errors, fmt.Errorf("%w: %w", ErrCannotReadJSONBytes, err))

		return
	}

	c.values = append(c.values, value)
}
