package confiq

// ConfigSetOptions is exposed so that functions which wrap the New function can make adding the WithTag option easier.
type ConfigSetOptions []loadOption

type configSetOption func(*ConfigSet)

// WithTag sets the struct tag to be used by the decoder for reading configuration values of struct fields.
func WithTag(tag string) configSetOption {
	return func(s *ConfigSet) {
		s.decoder.tag = tag
	}
}

// LoadOptions is exposed so that functions which wrap the Load function can make adding the WithPrefix option easier.
type LoadOptions []loadOption

type loadOption func(*loader)

// WithPrefix sets the prefix to be used when loading configuration values into the ConfigSet.
func WithPrefix(prefix string) loadOption {
	return func(l *loader) {
		l.prefix = prefix
	}
}

// DecodeOptions is exposed so that functions which wrap the Decode function can make adding the AsStrict and FromPrefix options easier.
type DecodeOptions []decodeOption

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
