package confiq

type IValueContainer interface {
	Get() []any
	Errors() []error
}
