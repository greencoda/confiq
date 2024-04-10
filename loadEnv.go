package confiq

import (
	"bytes"
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

// LoadEnvFromEnvironment loads config data from the environment variables.
func (c *ConfigSet) LoadEnvFromEnvironment(options ...loadOption) error {
	var (
		envSlice = os.Environ()
		envMap   = make(map[string]any)
	)

	for _, envElement := range envSlice {
		envKeyValue := strings.SplitN(envElement, envSplitChar, envSplitElements)

		envMap[envKeyValue[0]] = envKeyValue[1]
	}

	return c.applyValue(envMap, options...)
}

// LoadEnvFromFile loads an Env file from the given path.
func (c *ConfigSet) LoadEnvFromFile(path string, options ...loadOption) error {
	reader, err := os.Open(filepath.Clean(path))
	if err != nil {
		return fmt.Errorf("%w: %w", ErrCannotOpenConfig, err)
	}

	return c.LoadEnvFromReader(reader, options...)
}

// LoadEnvFromString loads an Env file from the given string.
func (c *ConfigSet) LoadEnvFromString(input string, options ...loadOption) error {
	reader := strings.NewReader(input)

	return c.LoadEnvFromReader(reader, options...)
}

// LoadEnvFromBytes loads an Env file from the given bytes.
func (c *ConfigSet) LoadEnvFromBytes(input []byte, options ...loadOption) error {
	reader := bytes.NewReader(input)

	return c.LoadEnvFromReader(reader, options...)
}

// LoadEnvFromReader loads an Env file from a reader stream.
func (c *ConfigSet) LoadEnvFromReader(input io.Reader, options ...loadOption) error {
	envMap, err := envparse.Parse(input)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrCannotLoadConfig, err)
	}

	return c.applyValue(c.envMapAsAny(envMap), options...)
}

func (c *ConfigSet) envMapAsAny(envMap map[string]string) map[string]any {
	anyMap := make(map[string]any)

	for key, value := range envMap {
		anyMap[key] = value
	}

	return anyMap
}
