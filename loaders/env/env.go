// Package confiqenv allows confiq values to be loaded from Env format.
package confiqenv

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/go-envparse"
)

const (
	envSplitChar     = "="
	envSplitElements = 2
)

var (
	ErrCannotGetBytesFromReader = errors.New("cannot get bytes from reader")
	ErrCannotOpenEnvFile        = errors.New("cannot open Env file")
	ErrCannotReadEnvData        = errors.New("cannot read Env data")
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

// FromEnvironment loads config data from the environment variables.
func (c *Container) FromEnvironment() *Container {
	var (
		envSlice = os.Environ()
		envMap   = make(map[string]any)
	)

	for _, envElement := range envSlice {
		envKeyValue := strings.SplitN(envElement, envSplitChar, envSplitElements)

		envMap[envKeyValue[0]] = envKeyValue[1]
	}

	c.values = append(c.values, envMap)

	return nil
}

// FromFile loads a Env file from the given path.
func (c *Container) FromFile(path string) *Container {
	inputBytes, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		c.errors = append(c.errors, fmt.Errorf("%w: %w", ErrCannotOpenEnvFile, err))

		return c
	}

	return c.FromBytes(inputBytes)
}

// FromString loads a Env file from the given string.
func (c *Container) FromString(input string) *Container {
	return c.FromBytes([]byte(input))
}

// FromBytes loads a Env file from the given bytes.
func (c *Container) FromBytes(input []byte) *Container {
	return c.FromReader(bytes.NewReader(input))
}

// FromReader loads a Env file from a reader stream.
func (c *Container) FromReader(reader io.Reader) *Container {
	if reader == nil {
		c.errors = append(c.errors, ErrCannotReadEnvData)

		return c
	}

	parsedEnvMap, err := envparse.Parse(reader)
	if err != nil {
		c.errors = append(c.errors, fmt.Errorf("%w: %w", ErrCannotReadEnvData, err))

		return c
	}

	c.values = append(c.values, envMapAsAny(parsedEnvMap))

	return c
}

func envMapAsAny(envMap map[string]string) map[string]any {
	anyMap := make(map[string]any)

	for key, value := range envMap {
		anyMap[key] = value
	}

	return anyMap
}
