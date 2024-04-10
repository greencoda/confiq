package confiq

// OverrideValue overrides the value of the ConfigSet during testing
// This is used to test certain decoding methods where the error handling
// depends on the value of the ConfigSet being loaded with certain invalid types
// which would not be possible via the usual value loading methods.
func (c *ConfigSet) OverrideValue(value any) {
	c.value = &value
}
