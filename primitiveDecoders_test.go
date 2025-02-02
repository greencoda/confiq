package confiq_test

import (
	"testing"

	"github.com/greencoda/confiq"
	"github.com/greencoda/confiq/mocks"
	"github.com/stretchr/testify/suite"
)

type PrimitiveDecodersTestSuite struct {
	suite.Suite

	configSet      *confiq.ConfigSet
	valueContainer *mocks.IValueContainer
}

func Test_PrimitiveDecoders(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(PrimitiveDecodersTestSuite))
}

func (s *PrimitiveDecodersTestSuite) SetupTest() {
	s.configSet = confiq.New(
		confiq.WithTag("cfg"),
	)
	s.valueContainer = mocks.NewIValueContainer(s.T())

	s.Require().NotNil(s.T(), s.configSet)
}

func (s *PrimitiveDecodersTestSuite) Test_DecodeString() {
	s.valueContainer.On("Errors").Return([]error{})
	s.valueContainer.On("Get").Return([]any{map[string]any{"test_string": "test"}})

	loadErr := s.configSet.Load(s.valueContainer)
	s.Require().NoError(loadErr)

	type targetStruct struct {
		TestString string `cfg:"test_string"`
	}

	var (
		target   targetStruct
		expected = "test"
	)

	decodeErr := s.configSet.Decode(&target)

	s.Equal(expected, target.TestString)
	s.NoError(decodeErr)
}

func (s *PrimitiveDecodersTestSuite) Test_DecodeString_Ptr() {
	s.valueContainer.On("Errors").Return([]error{})
	s.valueContainer.On("Get").Return([]any{map[string]any{"test_string": "test"}})

	loadErr := s.configSet.Load(s.valueContainer)
	s.Require().NoError(loadErr)

	type targetStruct struct {
		TestString *string `cfg:"test_string"`
	}

	var (
		target   targetStruct
		expected = "test"
	)

	decodeErr := s.configSet.Decode(&target)

	s.Require().NotNil(target.TestString)
	s.Equal(expected, *target.TestString)
	s.NoError(decodeErr)
}

func (s *PrimitiveDecodersTestSuite) Test_DecodeString_Ptr_Nil() {
	s.valueContainer.On("Errors").Return([]error{})
	s.valueContainer.On("Get").Return([]any{map[string]any{}})

	loadErr := s.configSet.Load(s.valueContainer)
	s.Require().NoError(loadErr)

	type targetStruct struct {
		TestString *string `cfg:"test_string"`
	}

	var target targetStruct

	decodeErr := s.configSet.Decode(&target)

	s.Nil(target.TestString)
	s.Error(decodeErr)
}

func (s *PrimitiveDecodersTestSuite) Test_DecodeBool() {
	s.valueContainer.On("Errors").Return([]error{})
	s.valueContainer.On("Get").Return([]any{map[string]any{"test_bool": true}})

	loadErr := s.configSet.Load(s.valueContainer)
	s.Require().NoError(loadErr)

	type targetStruct struct {
		TestBool bool `cfg:"test_bool"`
	}

	var (
		target   targetStruct
		expected = true
	)

	decodeErr := s.configSet.Decode(&target)

	s.Equal(expected, target.TestBool)
	s.NoError(decodeErr)
}

func (s *PrimitiveDecodersTestSuite) Test_DecodeBool_FromString() {
	s.valueContainer.On("Errors").Return([]error{})
	s.valueContainer.On("Get").Return([]any{map[string]any{"test_bool_string": "true"}})

	loadErr := s.configSet.Load(s.valueContainer)
	s.Require().NoError(loadErr)

	type targetStruct struct {
		TestBool bool `cfg:"test_bool_string"`
	}

	var (
		target   targetStruct
		expected = true
	)

	decodeErr := s.configSet.Decode(&target)

	s.Equal(expected, target.TestBool)
	s.NoError(decodeErr)
}

func (s *PrimitiveDecodersTestSuite) Test_DecodeBool_FromInvalidFormat() {
	s.valueContainer.On("Errors").Return([]error{})
	s.valueContainer.On("Get").Return([]any{map[string]any{"test_bool_invalid_format": 1.1}})

	loadErr := s.configSet.Load(s.valueContainer)
	s.Require().NoError(loadErr)

	type targetStruct struct {
		TestDuration bool `cfg:"test_bool_invalid_format"`
	}

	var target targetStruct

	decodeErr := s.configSet.Decode(&target)

	s.Error(decodeErr)
}

func (s *PrimitiveDecodersTestSuite) Test_DecodeFloat64() {
	s.valueContainer.On("Errors").Return([]error{})
	s.valueContainer.On("Get").Return([]any{map[string]any{"test_float": 0.1234567890123456}})

	loadErr := s.configSet.Load(s.valueContainer)
	s.Require().NoError(loadErr)

	type targetStruct struct {
		TestFloat float64 `cfg:"test_float"`
	}

	var (
		target   targetStruct
		expected = 0.1234567890123456
	)

	decodeErr := s.configSet.Decode(&target)

	s.InEpsilon(expected, target.TestFloat, 0.02)
	s.NoError(decodeErr)
}

func (s *PrimitiveDecodersTestSuite) Test_DecodeFloat_FromFloat32() {
	type targetStruct struct {
		TestFloat float32 `cfg:"test_float"`
	}

	var (
		target   targetStruct
		expected float32 = 0.1234567890123456
	)

	s.configSet.OverrideValue(map[string]any{"test_float": expected})

	decodeErr := s.configSet.Decode(&target)

	s.InEpsilon(expected, target.TestFloat, 0.02)
	s.NoError(decodeErr)
}

func (s *PrimitiveDecodersTestSuite) Test_DecodeFloat_FromString() {
	s.valueContainer.On("Errors").Return([]error{})
	s.valueContainer.On("Get").Return([]any{map[string]any{"test_float_string": "0.1234567890123456"}})

	loadErr := s.configSet.Load(s.valueContainer)
	s.Require().NoError(loadErr)

	type targetStruct struct {
		TestFloat float64 `cfg:"test_float_string"`
	}

	var (
		target   targetStruct
		expected = 0.1234567890123456
	)

	decodeErr := s.configSet.Decode(&target)

	s.InEpsilon(expected, target.TestFloat, 0.02)
	s.NoError(decodeErr)
}

func (s *PrimitiveDecodersTestSuite) Test_DecodeFloat_FromInvalidFormat() {
	s.valueContainer.On("Errors").Return([]error{})
	s.valueContainer.On("Get").Return([]any{map[string]any{"test_float_invalid_format": true}})

	loadErr := s.configSet.Load(s.valueContainer)
	s.Require().NoError(loadErr)

	type targetStruct struct {
		TestFloat float64 `cfg:"test_float_invalid_format"`
	}

	var target targetStruct

	decodeErr := s.configSet.Decode(&target)

	s.Error(decodeErr)
}

func (s *PrimitiveDecodersTestSuite) Test_DecodeInt() {
	type targetStruct struct {
		TestInt int `cfg:"test_int"`
	}

	var (
		target   targetStruct
		expected = 64
	)

	s.configSet.OverrideValue(map[string]any{"test_int": expected})

	decodeErr := s.configSet.Decode(&target)

	s.Equal(expected, target.TestInt)
	s.NoError(decodeErr)
}

func (s *PrimitiveDecodersTestSuite) Test_DecodeInt8() {
	type targetStruct struct {
		TestInt8 int8 `cfg:"test_int"`
	}

	var (
		target   targetStruct
		expected int8 = 64
	)

	s.configSet.OverrideValue(map[string]any{"test_int": expected})

	decodeErr := s.configSet.Decode(&target)

	s.Equal(expected, target.TestInt8)
	s.NoError(decodeErr)
}

func (s *PrimitiveDecodersTestSuite) Test_DecodeInt16() {
	type targetStruct struct {
		TestInt16 int16 `cfg:"test_int"`
	}

	var (
		target   targetStruct
		expected int16 = 64
	)

	s.configSet.OverrideValue(map[string]any{"test_int": expected})

	decodeErr := s.configSet.Decode(&target)

	s.Equal(expected, target.TestInt16)
	s.NoError(decodeErr)
}

func (s *PrimitiveDecodersTestSuite) Test_DecodeInt32() {
	type targetStruct struct {
		TestInt32 int32 `cfg:"test_int"`
	}

	var (
		target   targetStruct
		expected int32 = 64
	)

	s.configSet.OverrideValue(map[string]any{"test_int": expected})

	decodeErr := s.configSet.Decode(&target)

	s.Equal(expected, target.TestInt32)
	s.NoError(decodeErr)
}

func (s *PrimitiveDecodersTestSuite) Test_DecodeInt64() {
	type targetStruct struct {
		TestInt64 int64 `cfg:"test_int"`
	}

	var (
		target   targetStruct
		expected int64 = 64
	)

	s.configSet.OverrideValue(map[string]any{"test_int": expected})

	decodeErr := s.configSet.Decode(&target)

	s.Equal(expected, target.TestInt64)
	s.NoError(decodeErr)
}

func (s *PrimitiveDecodersTestSuite) Test_DecodeInt_FromString() {
	s.valueContainer.On("Errors").Return([]error{})
	s.valueContainer.On("Get").Return([]any{map[string]any{"test_int_string": "64"}})

	loadErr := s.configSet.Load(s.valueContainer)
	s.Require().NoError(loadErr)

	type targetStruct struct {
		TestInt64 int64 `cfg:"test_int_string"`
	}

	var (
		target   targetStruct
		expected int64 = 64
	)

	decodeErr := s.configSet.Decode(&target)

	s.Equal(expected, target.TestInt64)
	s.NoError(decodeErr)
}

func (s *PrimitiveDecodersTestSuite) Test_DecodeInt_FromInvalidFormat() {
	s.valueContainer.On("Errors").Return([]error{})
	s.valueContainer.On("Get").Return([]any{map[string]any{"test_int_invalid_format": true}})

	loadErr := s.configSet.Load(s.valueContainer)
	s.Require().NoError(loadErr)

	type targetStruct struct {
		TestInt64 int64 `cfg:"test_int_invalid_format"`
	}

	var target targetStruct

	decodeErr := s.configSet.Decode(&target)

	s.Error(decodeErr)
}

func (s *PrimitiveDecodersTestSuite) Test_DecodeUInt() {
	type targetStruct struct {
		TestInt uint `cfg:"test_int"`
	}

	var (
		target   targetStruct
		expected uint = 64
	)

	s.configSet.OverrideValue(map[string]any{"test_int": expected})

	decodeErr := s.configSet.Decode(&target)

	s.Equal(expected, target.TestInt)
	s.NoError(decodeErr)
}

func (s *PrimitiveDecodersTestSuite) Test_DecodeUInt8() {
	type targetStruct struct {
		TestInt8 uint8 `cfg:"test_int"`
	}

	var (
		target   targetStruct
		expected uint8 = 64
	)

	s.configSet.OverrideValue(map[string]any{"test_int": expected})

	decodeErr := s.configSet.Decode(&target)

	s.Equal(expected, target.TestInt8)
	s.NoError(decodeErr)
}

func (s *PrimitiveDecodersTestSuite) Test_DecodeUInt16() {
	type targetStruct struct {
		TestInt16 uint16 `cfg:"test_int"`
	}

	var (
		target   targetStruct
		expected uint16 = 64
	)

	s.configSet.OverrideValue(map[string]any{"test_int": expected})

	decodeErr := s.configSet.Decode(&target)

	s.Equal(expected, target.TestInt16)
	s.NoError(decodeErr)
}

func (s *PrimitiveDecodersTestSuite) Test_DecodeUInt32() {
	type targetStruct struct {
		TestInt32 uint32 `cfg:"test_int"`
	}

	var (
		target   targetStruct
		expected uint32 = 64
	)

	s.configSet.OverrideValue(map[string]any{"test_int": expected})

	decodeErr := s.configSet.Decode(&target)

	s.Equal(expected, target.TestInt32)
	s.NoError(decodeErr)
}

func (s *PrimitiveDecodersTestSuite) Test_DecodeUInt64() {
	type targetStruct struct {
		TestInt64 uint64 `cfg:"test_int"`
	}

	var (
		target   targetStruct
		expected uint64 = 64
	)

	s.configSet.OverrideValue(map[string]any{"test_int": expected})

	decodeErr := s.configSet.Decode(&target)

	s.Equal(expected, target.TestInt64)
	s.NoError(decodeErr)
}

func (s *PrimitiveDecodersTestSuite) Test_DecodeUInt_FromString() {
	s.valueContainer.On("Errors").Return([]error{})
	s.valueContainer.On("Get").Return([]any{map[string]any{"test_int_string": "64"}})

	loadErr := s.configSet.Load(s.valueContainer)
	s.Require().NoError(loadErr)

	type targetStruct struct {
		TestInt64 uint64 `cfg:"test_int_string"`
	}

	var (
		target   targetStruct
		expected uint64 = 64
	)

	decodeErr := s.configSet.Decode(&target)

	s.Equal(expected, target.TestInt64)
	s.NoError(decodeErr)
}

func (s *PrimitiveDecodersTestSuite) Test_DecodeUInt_FromInvalidFormat() {
	s.valueContainer.On("Errors").Return([]error{})
	s.valueContainer.On("Get").Return([]any{map[string]any{"test_int_invalid_format": true}})

	loadErr := s.configSet.Load(s.valueContainer)
	s.Require().NoError(loadErr)

	type targetStruct struct {
		TestInt64 uint64 `cfg:"test_int_invalid_format"`
	}

	var target targetStruct

	decodeErr := s.configSet.Decode(&target)

	s.Error(decodeErr)
}
