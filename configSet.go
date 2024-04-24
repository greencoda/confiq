// Package confiq is a Go package for populating structs from structured data formats such JSON, TOML, YAML or Env.
package confiq

import (
	"errors"
	"fmt"
	"reflect"
)

const defaultTag = "cfg"

var (
	errCannotGetKeyFromNonMap     = errors.New("cannot get key from non-map type")
	errKeyNotFound                = errors.New("key not found")
	errCannotGetIndexFromNonSlice = errors.New("cannot get index from non-slice type")
	errIndexOutOfBounds           = errors.New("index out of bounds")
)

type decoder struct {
	tag string
}

type decodeSettings struct {
	strict bool
	prefix string
}

// ConfigSet is a configuration set that can be used to load and decode configuration values into a struct.
type ConfigSet struct {
	value   *any
	decoder *decoder
}

// New creates a new ConfigSet with the given options.
func New(options ...configSetOption) *ConfigSet {
	var (
		value     any
		configSet = &ConfigSet{
			value:   &value,
			decoder: &decoder{tag: defaultTag},
		}
	)

	for _, option := range options {
		option(configSet)
	}

	return configSet
}

// Get returns the configuration value at the given path as an interface.
func (c *ConfigSet) Get(path string) (any, error) {
	return c.getByPath(path)
}

func (c *ConfigSet) getByPath(path string) (any, error) {
	currentValue := *c.value

	for path != "" {
		var currentSegment segment

		currentSegment, path = getNextSegment(path)

		switch v := currentSegment.(type) {
		case keySegment:
			segmentValue, err := getMapValue(currentValue, v.asString())
			if err != nil {
				return nil, err
			}

			currentValue = segmentValue
		case indexSegment:
			segmentValue, err := getSliceValue(currentValue, v.asInt())
			if err != nil {
				return nil, err
			}

			currentValue = segmentValue
		}
	}

	return currentValue, nil
}

func getMapValue(originMap any, key string) (any, error) {
	v := reflect.ValueOf(originMap)
	if v.Kind() != reflect.Map {
		return nil, fmt.Errorf("%w: %s(%T)", errCannotGetKeyFromNonMap, key, originMap)
	}

	value := v.MapIndex(reflect.ValueOf(key))
	if !value.IsValid() {
		return nil, fmt.Errorf("%w: %s", errKeyNotFound, key)
	}

	return value.Interface(), nil
}

func getSliceValue(originSlice any, index int) (any, error) {
	v := reflect.ValueOf(originSlice)
	if v.Kind() != reflect.Slice {
		return nil, fmt.Errorf("%w: %d(%T)", errCannotGetIndexFromNonSlice, index, originSlice)
	}

	if v.Len() <= index || index < 0 {
		return nil, fmt.Errorf("%w: %d", errIndexOutOfBounds, index)
	}

	return v.Index(index).Interface(), nil
}
