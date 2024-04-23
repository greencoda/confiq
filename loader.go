package confiq

import (
	"errors"
	"fmt"
	"maps"
)

const noPrefix = ""

// Collection of loader errors.
var (
	ErrCannotOpenConfig = errors.New("cannot open config")
	ErrCannotLoadConfig = errors.New("cannot load config")
)

var (
	errValueCannotBeNil           = errors.New("value cannot be nil")
	errCannotApplyValueOfThisType = errors.New("cannot apply value of this type")
	errCannotApplyMapValue        = errors.New("cannot apply map value")
	errCannotApplySliceValue      = errors.New("cannot apply slice value")
)

type loader struct {
	prefix string
}

func (c *ConfigSet) Load(valueContainer IValueContainer, options ...loadOption) error {
	if errs := valueContainer.Errors(); len(errs) > 0 {
		return fmt.Errorf("%w: %w", ErrCannotLoadConfig, errors.Join(errs...))
	}

	return c.applyValues(valueContainer.Get(), options...)
}

// LoadRawValue loads a raw value into the config set.
// The value must be a map[string]any or a slice of any.
func (c *ConfigSet) LoadRawValue(newValues []any, options ...loadOption) error {
	if newValues == nil {
		return errValueCannotBeNil
	}

	return c.applyValues(newValues, options...)
}

func newLoader() *loader {
	return &loader{
		prefix: noPrefix,
	}
}

func (c *ConfigSet) applyValues(newValues []any, options ...loadOption) error {
	loader := newLoader()

	for _, option := range options {
		option(loader)
	}

	for _, newValue := range newValues {
		switch v := newValue.(type) {
		case map[string]any:
			if err := c.applyMap(v, loader.prefix); err != nil {
				return err
			}
		case []any:
			if err := c.applySlice(v, loader.prefix); err != nil {
				return err
			}
		default:
			return fmt.Errorf("%w: %T", errCannotApplyValueOfThisType, newValue)
		}
	}

	return nil
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
