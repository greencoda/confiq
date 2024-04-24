package confiq

type configSetOption func(*ConfigSet)

// WithTag sets the struct tag to be used by the decoder for reading configuration values of struct fields.
func WithTag(tag string) configSetOption {
	return func(s *ConfigSet) {
		s.decoder.tag = tag
	}
}

type loadOption func(*loader)

// WithPrefix sets the prefix to be used when loading configuration values into the ConfigSet.
func WithPrefix(prefix string) loadOption {
	return func(l *loader) {
		l.prefix = prefix
	}
}

type decodeOption func(*decodeSettings)

// AsStrict sets the decoder to decode the configuration values as if all fields are set to strict.
func AsStrict() decodeOption {
	return func(d *decodeSettings) {
		d.strict = true
	}
}

// FromPrefix sets the prefix to be used when decoding configuration values into the target struct.
func FromPrefix(prefix string) decodeOption {
	return func(d *decodeSettings) {
		d.prefix = prefix
	}
}
