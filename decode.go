package confiq

import (
	"encoding"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

const sliceSplitChar = ";"

// Collection of decode errors.
var (
	ErrInvalidTarget        = errors.New("target must be non-nil pointer to a slice, map or struct that has at least one exported field with a the configured tag")
	ErrNoTargetFieldsAreSet = errors.New("none of the target fields were set from config values")
)

var (
	errCannotDecodeCustomTypeField       = errors.New("cannot decode field with custom decoder")
	errCannotDecodeNonRequiredField      = errors.New("cannot decode non-strict field")
	errCannotDecodeNonSliceValueToTarget = errors.New("cannot decode non-slice value to target")
	errCannotUnmarshalPrimitive          = errors.New("cannot unmarshal primitive as text")
	errCannotHaveDefaultForRequiredField = errors.New("cannot have default value for required field")
	errUnsupportedPrimitiveKind          = errors.New("unsupported primitive kind")
)

type (
	decoderFunc      func(targetField reflect.Value, value any) error
	fieldDecoderFunc func(targetField reflect.Value, value any, strict bool) (int, error)
)

type fieldOptions struct {
	path         string
	strict       bool
	required     bool
	defaultValue *string
}

type Decoder interface {
	Decode(value any) error
}

// Decode decodes the configuration values into the target struct.
func (c *ConfigSet) Decode(target interface{}, options ...decodeOption) error {
	return c.decode(target, options)
}

func (c *ConfigSet) decode(target interface{}, options []decodeOption) error {
	decodeSettings := &decodeSettings{
		strict: false,
		prefix: "",
	}

	for _, option := range options {
		option(decodeSettings)
	}

	targetValue := reflect.ValueOf(target)
	if targetValue.Kind() != reflect.Ptr || targetValue.IsNil() {
		return ErrInvalidTarget
	}

	targetValue = targetValue.Elem()

	decodedFieldCount, err := c.decodeField(targetValue, fieldOptions{
		path:         decodeSettings.prefix,
		strict:       decodeSettings.strict,
		required:     false,
		defaultValue: nil,
	})
	if err != nil {
		return err
	} else if decodedFieldCount == 0 {
		return ErrNoTargetFieldsAreSet
	}

	return nil
}

func (c *ConfigSet) getFieldConfigValue(fieldOpts fieldOptions) (any, error) {
	if fieldOpts.required && fieldOpts.defaultValue != nil {
		return nil, fmt.Errorf("%w: %s", errCannotHaveDefaultForRequiredField, fieldOpts.path)
	}

	configValue, err := c.getByPath(fieldOpts.path)
	if err != nil {
		if fieldOpts.required {
			return nil, fmt.Errorf("field is required: %w", err)
		}

		if fieldOpts.defaultValue != nil {
			return *fieldOpts.defaultValue, nil
		}

		return nil, errCannotDecodeNonRequiredField
	}

	return configValue, nil
}

func (c *ConfigSet) decodeField(targetValue reflect.Value, fieldOpts fieldOptions) (int, error) {
	fieldConfigValue, err := c.getFieldConfigValue(fieldOpts)
	if err != nil {
		if !errors.Is(err, errCannotDecodeNonRequiredField) {
			return 0, err
		}
	}

	var (
		decodedFields int
		decodeErr     error
		fieldDecoder  fieldDecoderFunc
	)

	if commonDecoder := getCommonDecoder(targetValue.Type()); commonDecoder != nil {
		return c.decodeCommon(commonDecoder, targetValue, fieldConfigValue, fieldOpts.strict)
	}

	if targetValue.Kind() == reflect.Ptr {
		if fieldConfigValue == nil {
			return 0, nil
		}

		dereferencedTargetValue := reflect.New(targetValue.Type().Elem()).Elem()

		decodedFields, err := c.decodeField(dereferencedTargetValue, fieldOpts)
		if err != nil {
			if fieldOpts.strict {
				return 0, fmt.Errorf("error decoding pointer value: %w", err)
			}

			return 0, nil
		}

		targetValue.Set(dereferencedTargetValue.Addr())

		return decodedFields, nil
	}

	// check if targetValue implements Decoder interface
	if decoder, ok := targetValue.Addr().Interface().(Decoder); ok {
		if err := decoder.Decode(fieldConfigValue); err != nil {
			return 0, fmt.Errorf("%w: %w", errCannotDecodeCustomTypeField, err)
		}

		return 1, nil
	}

	switch targetValue.Kind() {
	case reflect.Map:
		fieldDecoder = c.decodeMap
	case reflect.Slice:
		fieldDecoder = c.decodeSlice
	case reflect.Struct:
		fieldDecoder = c.decodeStruct
	default:
		fieldDecoder = c.decodePrimitiveType
	}

	decodedFields, decodeErr = fieldDecoder(targetValue, fieldConfigValue, fieldOpts.strict)
	if decodeErr != nil {
		return 0, decodeErr
	}

	return decodedFields, nil
}

func (c *ConfigSet) decodeCommon(commonDecoder decoderFunc, targetValue reflect.Value, fieldConfigValue any, strict bool) (int, error) {
	for targetValue.Kind() == reflect.Ptr {
		if fieldConfigValue == nil {
			return 0, nil
		}

		if targetValue.IsNil() {
			targetValue.Set(reflect.New(targetValue.Type().Elem()))
		}

		targetValue = targetValue.Elem()
	}

	if err := commonDecoder(targetValue, fieldConfigValue); err != nil {
		if strict {
			return 0, fmt.Errorf("error decoding field value: %w", err)
		}

		return 0, nil
	}

	return 1, nil
}

func (c *ConfigSet) decodeMap(targetMapValue reflect.Value, configValue any, strict bool) (int, error) {
	var (
		configMapValue     = reflect.ValueOf(configValue)
		targetMapValueType = targetMapValue.Type()
		targetKeyType      = targetMapValueType.Key()
		targetValueType    = targetMapValueType.Elem()
		setFieldCount      = 0
	)

	// setup empty map
	targetMapValue.Set(reflect.MakeMap(targetMapValueType))

	for _, key := range configMapValue.MapKeys() {
		var (
			k = reflect.New(targetKeyType).Elem()
			v = reflect.New(targetValueType).Elem()
		)

		// decode map key
		_, err := c.decodePrimitiveType(k, key.Interface(), strict)
		if err != nil {
			return 0, fmt.Errorf("error decoding map key: %w", err)
		}

		// decode map value
		decodedFieldCount, err := c.subValue(configValue).
			decodeField(v, fieldOptions{
				path:         keySegment(key.String()).String(),
				strict:       strict,
				required:     false,
				defaultValue: nil,
			})
		if err != nil {
			return 0, fmt.Errorf("error decoding map value: %w", err)
		}

		targetMapValue.SetMapIndex(k, v)

		setFieldCount += decodedFieldCount
	}

	return setFieldCount, nil
}

func (c *ConfigSet) decodeStruct(targetStructValue reflect.Value, configValue any, strict bool) (int, error) {
	var (
		targetStructType = targetStructValue.Type()
		setFieldCount    = 0
	)

	for i := range targetStructValue.NumField() {
		// get the struct field's tag and options
		targetStructFieldOpts := c.readTag(targetStructType.Field(i), c.decoder.tag)

		// get the struct field's reflection value
		targetStructFieldValue := targetStructValue.Field(i)

		// check if the field is exported
		if !targetStructFieldValue.CanSet() || !targetStructFieldValue.Addr().CanInterface() {
			continue
		}

		// set the field's strictness
		targetStructFieldOpts.strict = strict || targetStructFieldOpts.strict

		// decode the field
		decodedFieldCount, err := c.subValue(configValue).
			decodeField(targetStructFieldValue, targetStructFieldOpts)
		if err != nil {
			return 0, fmt.Errorf("error decoding struct field value: %w", err)
		}

		setFieldCount += decodedFieldCount
	}

	return setFieldCount, nil
}

func (c *ConfigSet) decodeSlice(targetSliceValue reflect.Value, configValue any, strict bool) (int, error) {
	var (
		configSliceValue     = reflect.ValueOf(configValue)
		configSliceValueKind = configSliceValue.Kind()
		setFieldCount        = 0
	)

	if configSliceValueKind != reflect.Slice {
		if configSliceValueKind != reflect.String {
			return 0, fmt.Errorf("%w: %v", errCannotDecodeNonSliceValueToTarget, configSliceValueKind)
		}

		return c.decodeSlice(targetSliceValue, strings.Split(configSliceValue.String(), sliceSplitChar), strict)
	}

	configSliceValueLength := configSliceValue.Len()

	// Set slice to the same length as the config slice
	targetSliceValue.Set(reflect.MakeSlice(targetSliceValue.Type(), configSliceValueLength, configSliceValueLength))

	// Decode each element based on its type
	for i := range configSliceValueLength {
		decodedFieldCount, err := c.subValue(configValue).
			decodeField(targetSliceValue.Index(i), fieldOptions{
				path:         indexSegment(i).String(),
				strict:       strict,
				required:     false,
				defaultValue: nil,
			})
		if err != nil {
			return setFieldCount, fmt.Errorf("error decoding slice element value: %w", err)
		}

		setFieldCount += decodedFieldCount
	}

	return setFieldCount, nil
}

func (c *ConfigSet) decodePrimitiveType(primitiveValue reflect.Value, configValue any, strict bool) (int, error) {
	primitiveInterface := primitiveValue.Addr().Interface()

	// check if primitive implements encoding.TextUnmarshaler interface
	if unmarshaler, ok := primitiveInterface.(encoding.TextUnmarshaler); ok {
		if err := unmarshaler.UnmarshalText(castToBytes(configValue)); err != nil {
			return 0, fmt.Errorf("%w: %w", errCannotUnmarshalPrimitive, err)
		}

		primitiveValue.Set(reflect.ValueOf(primitiveInterface).Elem())

		return 1, nil
	}

	// select the appropriate decoder function based on the primitive's kind
	var (
		primitiveDecoderFunc decoderFunc
		primitiveValueKind   = primitiveValue.Kind()
	)

	switch primitiveValueKind {
	case reflect.Bool:
		primitiveDecoderFunc = decodeBool
	case reflect.String:
		primitiveDecoderFunc = decodeString
	case reflect.Float32, reflect.Float64:
		primitiveDecoderFunc = decodeFloat
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		primitiveDecoderFunc = decodeInt
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		primitiveDecoderFunc = decodeUint
	default:
		return 0, fmt.Errorf("%w: %v", errUnsupportedPrimitiveKind, primitiveValueKind)
	}

	if err := primitiveDecoderFunc(primitiveValue, configValue); err != nil {
		if strict {
			return 0, fmt.Errorf("error decoding primitive value: %w", err)
		}

		return 0, nil
	}

	return 1, nil
}

func (c *ConfigSet) readTag(field reflect.StructField, tag string) fieldOptions {
	fieldOpts := fieldOptions{
		path:         "",
		strict:       false,
		required:     false,
		defaultValue: nil,
	}

	tagValue := field.Tag.Get(tag)
	if tagValue == "" {
		return fieldOpts
	}

	tagParts := strings.Split(tagValue, ",")

	fieldOpts.path = tagParts[0]

	// read the remaining tag parts
	for _, part := range tagParts[1:] {
		if part == "strict" {
			fieldOpts.strict = true

			continue
		}

		if part == "required" {
			fieldOpts.required = true

			continue
		}

		if strings.HasPrefix(part, "default=") {
			devaultValue := part[8:]
			fieldOpts.defaultValue = &devaultValue

			continue
		}
	}

	return fieldOpts
}

func (c *ConfigSet) subValue(value any) *ConfigSet {
	return &ConfigSet{
		value:   &value,
		decoder: c.decoder,
	}
}

func castToBytes(value any) []byte {
	if value == nil {
		return nil
	}

	return []byte(castToString(value))
}

func castToString(value any) string {
	if value == nil {
		return ""
	}

	if stringValue, ok := value.(string); ok {
		return stringValue
	}

	return fmt.Sprintf("%v", value)
}
