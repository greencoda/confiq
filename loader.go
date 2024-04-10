package confiq

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"maps"
)

const noPrefix = ""

// Collection of loader errors.
var (
	ErrCannotOpenConfig = errors.New("cannot open config")
	ErrCannotLoadConfig = errors.New("cannot load config")
)

var (
	errCannotApplyValueOfThisType = errors.New("cannot apply value of this type")
	errCannotApplyMapValue        = errors.New("cannot apply map value")
	errCannotApplySliceValue      = errors.New("cannot apply slice value")
	errCannotGetBytesFromReader   = errors.New("cannot get bytes from reader")
)

type loader struct {
	prefix string
}

func newLoader() *loader {
	return &loader{
		prefix: noPrefix,
	}
}

func readerToBytes(reader io.Reader) ([]byte, error) {
	if reader == nil {
		return []byte{}, errCannotGetBytesFromReader
	}

	buffer := new(bytes.Buffer)

	if _, err := buffer.ReadFrom(reader); err != nil {
		return []byte{}, errCannotGetBytesFromReader
	}

	return buffer.Bytes(), nil
}

func (c *ConfigSet) applyValue(newValue any, options ...loadOption) error {
	loader := newLoader()

	for _, option := range options {
		option(loader)
	}

	switch v := newValue.(type) {
	case map[string]any:
		return c.applyMap(v, loader.prefix)
	case []any:
		return c.applySlice(v, loader.prefix)
	}

	return fmt.Errorf("%w: %T", errCannotApplyValueOfThisType, newValue)
}

func (c *ConfigSet) applyMap(newValue map[string]any, prefix string) error {
	if prefix == "" {
		return c.applyMapWithoutPrefix(newValue)
	}

	return c.applyMapWithPrefix(newValue, prefix)
}

func (c *ConfigSet) applySlice(newValue []any, prefix string) error {
	if prefix == "" {
		return c.applySliceWithoutPrefix(newValue)
	}

	return c.applySliceWithPrefix(newValue, prefix)
}

func (c *ConfigSet) applyMapWithoutPrefix(newValue map[string]any) error {
	if *c.value == nil {
		*c.value = newValue

		return nil
	}

	valueMap, ok := (*c.value).(map[string]any)
	if !ok {
		return errCannotApplyMapValue
	}

	maps.Copy(valueMap, newValue)

	return nil
}

func (c *ConfigSet) applyMapWithPrefix(newValue map[string]any, prefix string) error {
	if *c.value == nil {
		*c.value = map[string]any{}
	}

	valueMap, ok := (*c.value).(map[string]any)
	if !ok {
		return errCannotApplyMapValue
	}

	valueMapAtPath, ok := valueMap[prefix]
	if !ok {
		valueMap[prefix] = newValue

		return nil
	}

	valueMapAtPathMap, ok := valueMapAtPath.(map[string]any)
	if !ok {
		return errCannotApplyMapValue
	}

	maps.Copy(valueMapAtPathMap, newValue)

	return nil
}

func (c *ConfigSet) applySliceWithoutPrefix(newValue []any) error {
	if *c.value == nil {
		*c.value = newValue

		return nil
	}

	valueSlice, ok := (*c.value).([]any)
	if !ok {
		return errCannotApplySliceValue
	}

	valueSlice = append(valueSlice, newValue...)

	*c.value = valueSlice

	return nil
}

func (c *ConfigSet) applySliceWithPrefix(newValue []any, prefix string) error {
	if *c.value == nil {
		*c.value = map[string]any{}
	}

	valueMap, ok := (*c.value).(map[string]any)
	if !ok {
		return errCannotApplySliceValue
	}

	valueMapAtPath, ok := valueMap[prefix]
	if !ok {
		valueMap[prefix] = newValue

		return nil
	}

	valueMapAtPathSlice, ok := valueMapAtPath.([]any)
	if !ok {
		return errCannotApplySliceValue
	}

	valueMapAtPathSlice = append(valueMapAtPathSlice, newValue...)

	valueMap[prefix] = valueMapAtPathSlice

	return nil
}
