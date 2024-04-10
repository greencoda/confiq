package confiq

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/url"
	"reflect"
	"time"
)

var (
	errCannotParseIP             = errors.New("cannot parse IP address")
	errCannotParseURL            = errors.New("cannot parse URL")
	errCannotParseDuration       = errors.New("cannot parse duration")
	errCannotParseTime           = errors.New("cannot parse time")
	errCannotParseNonStringTime  = errors.New("cannot parse time from non-string type")
	errCannotParseJSONRawMessage = errors.New("cannot marshal source value to JSON")
)

type typeDefinition struct {
	kind     reflect.Kind
	pkgPath  string
	typeName string
}

var commonDecoders = map[typeDefinition]decoderFunc{
	{reflect.Int64, "time", "Duration"}:            decodeDuration,
	{reflect.Slice, "net", "IP"}:                   decodeIP,
	{reflect.Slice, "encoding/json", "RawMessage"}: decodeJSONRawMessage,
	{reflect.Struct, "time", "Time"}:               decodeTime,
	{reflect.Struct, "net/url", "URL"}:             decodeURL,
}

func getCommonDecoder(targetValType reflect.Type) decoderFunc {
	if targetValType.Kind() == reflect.Ptr {
		targetValType = targetValType.Elem()
	}

	if decoder, decoderFound := commonDecoders[typeDefinition{targetValType.Kind(), targetValType.PkgPath(), targetValType.Name()}]; decoderFound {
		return decoder
	}

	return nil
}

func decodeDuration(targetValue reflect.Value, sourceValue any) error {
	parsedDuration, err := time.ParseDuration(castToString(sourceValue))
	if err != nil {
		return fmt.Errorf("%w: %w", errCannotParseDuration, err)
	}

	targetValue.SetInt(int64(parsedDuration))

	return nil
}

func decodeIP(targetValue reflect.Value, sourceValue any) error {
	parsedIP := net.ParseIP(castToString(sourceValue))
	if parsedIP == nil {
		return errCannotParseIP
	}

	targetValue.Set(reflect.ValueOf(parsedIP))

	return nil
}

func decodeJSONRawMessage(targetValue reflect.Value, sourceValue any) error {
	marshaledSourceValue, err := json.Marshal(sourceValue)
	if err != nil {
		return fmt.Errorf("%w: %w", errCannotParseJSONRawMessage, err)
	}

	targetValue.Set(reflect.ValueOf(marshaledSourceValue))

	return nil
}

func decodeURL(targetValue reflect.Value, sourceValue any) error {
	parsedURL, err := url.Parse(castToString(sourceValue))
	if err != nil {
		return fmt.Errorf("%w: %w", errCannotParseURL, err)
	}

	targetValue.Set(reflect.ValueOf(parsedURL))

	return nil
}

func decodeTime(targetValue reflect.Value, sourceValue any) error {
	if timeValue, isTimeValue := sourceValue.(time.Time); isTimeValue {
		targetValue.Set(reflect.ValueOf(timeValue))

		return nil
	}

	stringValue, isStringValue := sourceValue.(string)
	if !isStringValue {
		return fmt.Errorf("%w: %T", errCannotParseNonStringTime, sourceValue)
	}

	parsedTime, err := time.Parse(time.RFC3339, stringValue)
	if err != nil {
		return fmt.Errorf("%w: %w", errCannotParseTime, err)
	}

	targetValue.Set(reflect.ValueOf(parsedTime))

	return nil
}
