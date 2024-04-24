package confiq_test

import (
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/greencoda/confiq"
	"github.com/greencoda/confiq/mocks"
	"github.com/stretchr/testify/suite"
)

var (
	errCannotParseInt   = errors.New("cannot parse int")
	errValueIsNotString = errors.New("value is not a string")
)

type unmarshalerNumber uint8

func (o *unmarshalerNumber) UnmarshalText(raw []byte) error {
	n, err := strconv.ParseInt(string(raw), 10, 8)
	if err != nil {
		return fmt.Errorf("%w: %w", errCannotParseInt, err)
	}

	*o = unmarshalerNumber(n)

	return nil
}

type customDecoderStruct struct {
	InternalValue string
}

func (t *customDecoderStruct) Decode(value any) error {
	if stringValue, ok := value.(string); !ok {
		return errValueIsNotString
	} else {
		t.InternalValue = stringValue

		return nil
	}
}

type DecodeTestSuite struct {
	suite.Suite

	configSet      *confiq.ConfigSet
	valueContainer *mocks.IValueContainer
}

func Test_Decoder(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(DecodeTestSuite))
}

func (s *DecodeTestSuite) SetupTest() {
	s.configSet = confiq.New(
		confiq.WithTag("cfg"),
	)
	s.valueContainer = mocks.NewIValueContainer(s.T())

	s.Require().NotNil(s.T(), s.configSet)
}

func (s *DecodeTestSuite) Test_Decode_UnsupportedPrimitive() {
	s.valueContainer.On("Errors").Return([]error{})
	s.valueContainer.On("Get").Return([]any{map[string]any{"test_int": 64}})

	loadErr := s.configSet.Load(s.valueContainer)
	s.Require().NoError(loadErr)

	type targetStruct struct {
		TestComplex complex64 `cfg:"test_int"`
	}

	var target targetStruct

	decodeErr := s.configSet.Decode(&target)

	s.Error(decodeErr)
}

func (s *DecodeTestSuite) Test_Decode_Slice() {
	s.valueContainer.On("Errors").Return([]error{})
	s.valueContainer.On("Get").Return([]any{
		map[string]any{
			"test_slice": []any{
				"uno",
				"dos",
				"tres",
			},
		},
	})

	loadErr := s.configSet.Load(s.valueContainer)
	s.Require().NoError(loadErr)

	type targetStruct struct {
		TestSlice []string `cfg:"test_slice"`
	}

	var (
		target   targetStruct
		expected = []string{
			"uno",
			"dos",
			"tres",
		}
	)

	decodeErr := s.configSet.Decode(&target)

	s.Equal(expected, target.TestSlice)
	s.NoError(decodeErr)
}

func (s *DecodeTestSuite) Test_Decode_Slice_FromInvalidFormat() {
	s.valueContainer.On("Errors").Return([]error{})
	s.valueContainer.On("Get").Return([]any{
		map[string]any{
			"test_map": map[string]any{
				"test_map_key_1": 1,
				"test_map_key_2": 2,
				"test_map_key_3": 3,
			},
		},
	})

	loadErr := s.configSet.Load(s.valueContainer)
	s.Require().NoError(loadErr)

	type targetStruct struct {
		TestSliceString []string `cfg:"test_map"`
	}

	var target targetStruct

	decodeErr := s.configSet.Decode(&target)

	s.Error(decodeErr)
}

func (s *DecodeTestSuite) Test_Decode_Slice_FromInvalidConfigValueType() {
	s.valueContainer.On("Errors").Return([]error{})
	s.valueContainer.On("Get").Return([]any{
		map[string]any{
			"test_slice": []any{
				"uno",
				"dos",
				"tres",
			},
		},
	})

	loadErr := s.configSet.Load(s.valueContainer)
	s.Require().NoError(loadErr)

	type targetStruct struct {
		TestSliceString []int `cfg:"test_slice,strict"`
	}

	var target targetStruct

	decodeErr := s.configSet.Decode(&target)

	s.Error(decodeErr)
}

func (s *DecodeTestSuite) Test_Decode_Map() {
	s.valueContainer.On("Errors").Return([]error{})
	s.valueContainer.On("Get").Return([]any{
		map[string]any{
			"test_map": map[string]any{
				"test_map_key_1": 1,
				"test_map_key_2": 2,
				"test_map_key_3": 3,
			},
		},
	})

	loadErr := s.configSet.Load(s.valueContainer)
	s.Require().NoError(loadErr)

	type targetStruct struct {
		TestMap map[string]int `cfg:"test_map"`
	}

	var (
		target   targetStruct
		expected = map[string]int{
			"test_map_key_1": 1,
			"test_map_key_2": 2,
			"test_map_key_3": 3,
		}
	)

	decodeErr := s.configSet.Decode(&target)

	s.Equal(expected, target.TestMap)
	s.NoError(decodeErr)
}

func (s *DecodeTestSuite) Test_Decode_Map_FromInvalidKeyFormat() {
	s.valueContainer.On("Errors").Return([]error{})
	s.valueContainer.On("Get").Return([]any{
		map[string]any{
			"test_map": map[string]any{
				"test_map_key_1": 1,
				"test_map_key_2": 2,
				"test_map_key_3": 3,
			},
		},
	})

	loadErr := s.configSet.Load(s.valueContainer)
	s.Require().NoError(loadErr)

	type targetStruct struct {
		TestMapInvalidKey map[bool]int `cfg:"test_map,strict"`
	}

	var target targetStruct

	decodeErr := s.configSet.Decode(&target)

	s.Error(decodeErr)
}

func (s *DecodeTestSuite) Test_Decode_Map_FromInvalidValueFormat() {
	s.valueContainer.On("Errors").Return([]error{})
	s.valueContainer.On("Get").Return([]any{
		map[string]any{
			"test_map": map[string]any{
				"test_map_key_1": 1,
				"test_map_key_2": 2,
				"test_map_key_3": 3,
			},
		},
	})

	loadErr := s.configSet.Load(s.valueContainer)
	s.Require().NoError(loadErr)

	type targetStruct struct {
		TestMapInvalidValue map[string]bool `cfg:"test_map,strict"`
	}

	var target targetStruct

	decodeErr := s.configSet.Decode(&target)

	s.Error(decodeErr)
}

func (s *DecodeTestSuite) Test_Decode_PointerField() {
	s.valueContainer.On("Errors").Return([]error{})
	s.valueContainer.On("Get").Return([]any{map[string]any{"test_bool": true}})

	loadErr := s.configSet.Load(s.valueContainer)
	s.Require().NoError(loadErr)

	type targetStruct struct {
		TestBoolPointer *bool `cfg:"test_bool"`
	}

	var (
		target   targetStruct
		expected = true
	)

	decodeErr := s.configSet.Decode(&target)

	s.Equal(&expected, target.TestBoolPointer)
	s.NoError(decodeErr)
}

func (s *DecodeTestSuite) Test_Decode_PointerField_FromInvalidFormat() {
	s.valueContainer.On("Errors").Return([]error{})
	s.valueContainer.On("Get").Return([]any{
		map[string]any{
			"test_section": map[string]any{
				"test_string": "efes",
			},
		},
	})

	loadErr := s.configSet.Load(s.valueContainer)
	s.Require().NoError(loadErr)

	type targetSubStruct struct {
		TestBoolPointerStrict *bool `cfg:"test_string,strict"`
	}

	type targetStruct struct {
		TestStruct *targetSubStruct `cfg:"test_section"`
	}

	var target targetStruct

	decodeErr := s.configSet.Decode(&target)

	s.Error(decodeErr)
}

func (s *DecodeTestSuite) Test_Decode_TextUnmarshalablePrimitive() {
	s.valueContainer.On("Errors").Return([]error{})
	s.valueContainer.On("Get").Return([]any{map[string]any{"test_int_string": "64"}})

	loadErr := s.configSet.Load(s.valueContainer)
	s.Require().NoError(loadErr)

	type targetStruct struct {
		TestInt unmarshalerNumber `cfg:"test_int_string"`
	}

	var (
		target   targetStruct
		expected = unmarshalerNumber(64)
	)

	decodeErr := s.configSet.Decode(&target)

	s.Equal(expected, target.TestInt)
	s.NoError(decodeErr)
}

func (s *DecodeTestSuite) Test_Decode_TextUnmarshalablePrimitive_FromInvalidFormat() {
	s.valueContainer.On("Errors").Return([]error{})
	s.valueContainer.On("Get").Return([]any{map[string]any{"test_bool": true}})

	loadErr := s.configSet.Load(s.valueContainer)
	s.Require().NoError(loadErr)

	type targetStruct struct {
		TestInt unmarshalerNumber `cfg:"test_bool"`
	}

	var target targetStruct

	decodeErr := s.configSet.Decode(&target)

	s.Error(decodeErr)
}

func (s *DecodeTestSuite) Test_Decode_CustomDecoder() {
	s.valueContainer.On("Errors").Return([]error{})
	s.valueContainer.On("Get").Return([]any{
		map[string]any{
			"test_section": map[string]any{
				"test_string": "efes",
			},
		},
	})

	loadErr := s.configSet.Load(s.valueContainer)
	s.Require().NoError(loadErr)

	type targetStruct struct {
		CustomStruct customDecoderStruct `cfg:"test_section.test_string"`
	}

	var (
		target   targetStruct
		expected = "efes"
	)

	decodeErr := s.configSet.Decode(&target)

	s.Equal(expected, target.CustomStruct.InternalValue)
	s.NoError(decodeErr)
}

func (s *DecodeTestSuite) Test_Decode_CustomDecoder_FromInvalidFormat() {
	s.valueContainer.On("Errors").Return([]error{})
	s.valueContainer.On("Get").Return([]any{
		map[string]any{
			"test_section": map[string]any{
				"test_string":       "efes",
				"test_string_array": []any{"aleph", "beth", "gimel"},
			},
		},
	})

	loadErr := s.configSet.Load(s.valueContainer)
	s.Require().NoError(loadErr)

	type targetStruct struct {
		CustomStruct customDecoderStruct `cfg:"test_section"`
	}

	var target targetStruct

	decodeErr := s.configSet.Decode(&target)

	s.Error(decodeErr)
}

func (s *DecodeTestSuite) Test_Decode_UnexportedField() {
	s.valueContainer.On("Errors").Return([]error{})
	s.valueContainer.On("Get").Return([]any{map[string]any{"test_string": "test"}})

	loadErr := s.configSet.Load(s.valueContainer)
	s.Require().NoError(loadErr)

	type targetStruct struct {
		TestExportedString   string `cfg:"test_string"`
		testUnexportedString string `cfg:"test_string"`
	}

	var (
		target             targetStruct
		expectedExported   = "test"
		expectedUnexported = ""
	)

	decodeErr := s.configSet.Decode(&target)

	s.Equal(expectedExported, target.TestExportedString)
	s.Equal(expectedUnexported, target.testUnexportedString)
	s.NoError(decodeErr)
}

func (s *DecodeTestSuite) Test_Decode_DefaultWithRequired() {
	s.valueContainer.On("Errors").Return([]error{})
	s.valueContainer.On("Get").Return([]any{map[string]any{"test_string": "test"}})

	loadErr := s.configSet.Load(s.valueContainer)
	s.Require().NoError(loadErr)

	type targetStruct struct {
		TestString string `cfg:"test_string,required,default=test"`
	}

	var target targetStruct

	decodeErr := s.configSet.Decode(&target)

	s.Error(decodeErr)
}

func (s *DecodeTestSuite) Test_Decode_RequiredFieldMissing() {
	s.valueContainer.On("Errors").Return([]error{})
	s.valueContainer.On("Get").Return([]any{map[string]any{}})

	loadErr := s.configSet.Load(s.valueContainer)
	s.Require().NoError(loadErr)

	type targetStruct struct {
		TestString string `cfg:"test_missing_string,required"`
	}

	var target targetStruct

	decodeErr := s.configSet.Decode(&target)

	s.Error(decodeErr)
}

func (s *DecodeTestSuite) Test_Decode_StringAsSlice() {
	s.valueContainer.On("Errors").Return([]error{})
	s.valueContainer.On("Get").Return([]any{map[string]any{"TEST_STRINGS": "test1;test2;test3"}})

	loadErr := s.configSet.Load(s.valueContainer)
	s.Require().NoError(loadErr)

	type targetStruct struct {
		TestStrings []string `cfg:"TEST_STRINGS"`
	}

	var (
		target   targetStruct
		expected = []string{
			"test1",
			"test2",
			"test3",
		}
	)

	decodeErr := s.configSet.Decode(&target)

	s.Equal(expected, target.TestStrings)
	s.NoError(decodeErr)
}

func (s *DecodeTestSuite) Test_Decode_InvalidAsSlice() {
	s.valueContainer.On("Errors").Return([]error{})
	s.valueContainer.On("Get").Return([]any{map[string]any{"test_bool": true}})

	loadErr := s.configSet.Load(s.valueContainer)
	s.Require().NoError(loadErr)

	type targetStruct struct {
		TestStrings []string `cfg:"test_bool"`
	}

	var target targetStruct

	decodeErr := s.configSet.Decode(&target)

	s.Error(decodeErr)
}

func (s *DecodeTestSuite) Test_Decode_DefaultField() {
	s.valueContainer.On("Errors").Return([]error{})
	s.valueContainer.On("Get").Return([]any{map[string]any{}})

	loadErr := s.configSet.Load(s.valueContainer)
	s.Require().NoError(loadErr)

	type targetStruct struct {
		TestString string `cfg:"test_missing_string,default=defaultValue"`
	}

	var (
		target   targetStruct
		expected = "defaultValue"
	)

	decodeErr := s.configSet.Decode(&target)

	s.Equal(expected, target.TestString)
	s.NoError(decodeErr)
}

func (s *DecodeTestSuite) Test_Decode_MissingField() {
	s.valueContainer.On("Errors").Return([]error{})
	s.valueContainer.On("Get").Return([]any{map[string]any{}})

	loadErr := s.configSet.Load(s.valueContainer)
	s.Require().NoError(loadErr)

	type targetStruct struct {
		TestString string `cfg:"test_missing_string"`
	}

	var target targetStruct

	decodeErr := s.configSet.Decode(&target)

	s.NoError(decodeErr)
}

func (s *DecodeTestSuite) Test_Decode_EmptyTag() {
	s.valueContainer.On("Errors").Return([]error{})
	s.valueContainer.On("Get").Return([]any{map[string]any{"test_string": "test"}})

	loadErr := s.configSet.Load(s.valueContainer)
	s.Require().NoError(loadErr)

	type targetSubStruct struct {
		TestString string `cfg:"test_string"`
	}

	type targetStruct struct {
		TestStruct targetSubStruct `cfg:""`
	}

	var (
		target   targetStruct
		expected = "test"
	)

	decodeErr := s.configSet.Decode(&target)

	s.Equal(expected, target.TestStruct.TestString)
	s.NoError(decodeErr)
}
