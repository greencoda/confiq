package confiq

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

var (
	errCannotParseBool  = errors.New("cannot parse bool")
	errCannotParseFloat = errors.New("cannot parse float")
	errCannotParseInt   = errors.New("cannot parse int")
	errCannotParseUint  = errors.New("cannot parse uint")
)

func decodeString(targetValue reflect.Value, sourceValue any) error {
	targetValue.SetString(castToString(sourceValue))

	return nil
}

func decodeBool(targetValue reflect.Value, sourceValue any) error {
	if boolValue, ok := sourceValue.(bool); ok {
		targetValue.SetBool(boolValue)

		return nil
	} else {
		parsedBool, err := strconv.ParseBool(castToString(sourceValue))
		if err != nil {
			return fmt.Errorf("%w: %w", errCannotParseBool, err)
		}

		targetValue.SetBool(parsedBool)
	}

	return nil
}

func decodeFloat(targetValue reflect.Value, sourceValue any) error {
	switch sV := sourceValue.(type) {
	case float32:
		targetValue.SetFloat(float64(sV))
	case float64:
		targetValue.SetFloat(sV)
	default:
		parsedFloat, err := strconv.ParseFloat(castToString(sourceValue), targetValue.Type().Bits())
		if err != nil {
			return fmt.Errorf("%w: %w", errCannotParseFloat, err)
		}

		targetValue.SetFloat(parsedFloat)
	}

	return nil
}

func decodeInt(targetValue reflect.Value, sourceValue any) error {
	switch sV := sourceValue.(type) {
	case int:
		targetValue.SetInt(int64(sV))
	case int8:
		targetValue.SetInt(int64(sV))
	case int16:
		targetValue.SetInt(int64(sV))
	case int32:
		targetValue.SetInt(int64(sV))
	case int64:
		targetValue.SetInt(sV)
	default:
		parsedInt, err := strconv.ParseInt(castToString(sourceValue), 0, targetValue.Type().Bits())
		if err != nil {
			return fmt.Errorf("%w: %w", errCannotParseInt, err)
		}

		targetValue.SetInt(parsedInt)
	}

	return nil
}

func decodeUint(targetValue reflect.Value, sourceValue any) error {
	switch sV := sourceValue.(type) {
	case uint:
		targetValue.SetUint(uint64(sV))
	case uint8:
		targetValue.SetUint(uint64(sV))
	case uint16:
		targetValue.SetUint(uint64(sV))
	case uint32:
		targetValue.SetUint(uint64(sV))
	case uint64:
		targetValue.SetUint(sV)
	default:
		parsedUint, err := strconv.ParseUint(castToString(sourceValue), 0, targetValue.Type().Bits())
		if err != nil {
			return fmt.Errorf("%w: %w", errCannotParseUint, err)
		}

		targetValue.SetUint(parsedUint)
	}

	return nil
}
