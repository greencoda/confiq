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
