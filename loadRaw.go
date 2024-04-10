package confiq

import "errors"

var errValueCannotBeNil = errors.New("value cannot be nil")

// LoadRawValue loads a raw value into the config set.
// The value must be a map[string]any or a slice of any.
func (c *ConfigSet) LoadRawValue(newValue any, options ...loadOption) error {
	if newValue == nil {
		return errValueCannotBeNil
	}

	return c.applyValue(newValue, options...)
}
