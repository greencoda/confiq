package confiq_test

import (
	"testing"

	"github.com/greencoda/confiq"
	"github.com/greencoda/confiq/mocks"
	"github.com/stretchr/testify/suite"
)

type ConfigSetTestSuite struct {
	suite.Suite

	configSet      *confiq.ConfigSet
	valueContainer *mocks.IValueContainer
}

func Test_ConfigSet(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(ConfigSetTestSuite))
}

func (s *ConfigSetTestSuite) SetupTest() {
	s.configSet = confiq.New(
		confiq.WithTag("cfg"),
	)
	s.valueContainer = mocks.NewIValueContainer(s.T())

	s.Require().NotNil(s.T(), s.configSet)
}

func (s *ConfigSetTestSuite) Test_Decode() {
	s.valueContainer.On("Errors").Return([]error{})
	s.valueContainer.On("Get").Return([]any{map[string]any{"test_string": "test"}})

	loadErr := s.configSet.Load(s.valueContainer)
	s.Require().NoError(loadErr)

	type JSONPrimitives struct {
		TestString string `cfg:"test_string"`
	}

	var jsonPrimitives JSONPrimitives

	decodeErr := s.configSet.Decode(&jsonPrimitives)

	s.Equal(JSONPrimitives{
		TestString: "test",
	}, jsonPrimitives)
	s.NoError(decodeErr)
}

func (s *ConfigSetTestSuite) Test_Decode_NotPointer() {
	s.valueContainer.On("Errors").Return([]error{})
	s.valueContainer.On("Get").Return([]any{map[string]any{"test_string": "test"}})

	loadErr := s.configSet.Load(s.valueContainer)
	s.Require().NoError(loadErr)

	type JSONPrimitives struct {
		TestString string `cfg:"test_string"`
	}

	var jsonPrimitives JSONPrimitives

	decodeErr := s.configSet.Decode(jsonPrimitives)

	s.Empty(jsonPrimitives)
	s.Error(decodeErr)
}

func (s *ConfigSetTestSuite) Test_Decode_KeySegmentIndex() {
	s.valueContainer.On("Errors").Return([]error{})
	s.valueContainer.On("Get").Return([]any{map[string]any{"test_string": "test"}})

	loadErr := s.configSet.Load(s.valueContainer)
	s.Require().NoError(loadErr)

	type JSONPrimitives struct {
		TestString string `cfg:"[test_string]"`
	}

	var jsonPrimitives JSONPrimitives

	decodeErr := s.configSet.Decode(&jsonPrimitives)

	s.Equal(JSONPrimitives{
		TestString: "test",
	}, jsonPrimitives)
	s.NoError(decodeErr)
}

func (s *ConfigSetTestSuite) Test_StrictDecode() {
	s.valueContainer.On("Errors").Return([]error{})
	s.valueContainer.On("Get").Return([]any{map[string]any{"test_string": "test"}})

	loadErr := s.configSet.Load(s.valueContainer)
	s.Require().NoError(loadErr)

	type JSONPrimitives struct {
		TestString string `cfg:"test_string"`
	}

	var jsonPrimitives JSONPrimitives

	decodeErr := s.configSet.StrictDecode(&jsonPrimitives)

	s.Equal(JSONPrimitives{
		TestString: "test",
	}, jsonPrimitives)
	s.NoError(decodeErr)
}

func (s *ConfigSetTestSuite) Test_Get_Primitives() {
	s.valueContainer.On("Errors").Return([]error{})
	s.valueContainer.On("Get").Return([]any{map[string]any{"test_string": "test"}})

	loadErr := s.configSet.Load(s.valueContainer)
	s.Require().NoError(loadErr)

	value, getErr := s.configSet.Get("test_string")

	s.Equal("test", value)
	s.NoError(getErr)
}

func (s *ConfigSetTestSuite) Test_Get_NestedPrimitives() {
	s.valueContainer.On("Errors").Return([]error{})
	s.valueContainer.On("Get").Return([]any{
		map[string]any{"test_section": map[string]any{
			"test_string_array": []any{"aleph", "beth", "gimel"},
		}},
	})

	loadErr := s.configSet.Load(s.valueContainer)
	s.Require().NoError(loadErr)

	value, getErr := s.configSet.Get("test_section.test_string_array[0]")

	s.Equal("aleph", value)
	s.NoError(getErr)
}

func (s *ConfigSetTestSuite) Test_Get_FromMap_WithInvalidKeyPath() {
	s.valueContainer.On("Errors").Return([]error{})
	s.valueContainer.On("Get").Return([]any{map[string]any{"test_string": "test"}})

	loadErr := s.configSet.Load(s.valueContainer)
	s.Require().NoError(loadErr)

	value, getErr := s.configSet.Get("noexistent_key")

	s.Empty(value)
	s.Error(getErr)
}

func (s *ConfigSetTestSuite) Test_Get_FromMap_WithInvalidIndexPath() {
	s.valueContainer.On("Errors").Return([]error{})
	s.valueContainer.On("Get").Return([]any{map[string]any{"test_string": "test"}})

	loadErr := s.configSet.Load(s.valueContainer)
	s.Require().NoError(loadErr)

	value, getErr := s.configSet.Get("[0]")

	s.Empty(value)
	s.Error(getErr)
}

func (s *ConfigSetTestSuite) Test_Get_FromArray_WithInvalidKeyPath() {
	s.valueContainer.On("Errors").Return([]error{})
	s.valueContainer.On("Get").Return([]any{
		[]any{
			map[string]any{"test_value": "test_1"},
			map[string]any{"test_value": "test_2"},
			map[string]any{"test_value": "test_3"},
		},
	})

	loadErr := s.configSet.Load(s.valueContainer)
	s.Require().NoError(loadErr)

	value, getErr := s.configSet.Get("noexistent_key")

	s.Empty(value)
	s.Error(getErr)
}

func (s *ConfigSetTestSuite) Test_Get_FromArray_WithInvalidIndexPath() {
	s.valueContainer.On("Errors").Return([]error{})
	s.valueContainer.On("Get").Return([]any{
		[]any{
			map[string]any{"test_value": "test_1"},
			map[string]any{"test_value": "test_2"},
			map[string]any{"test_value": "test_3"},
		},
	})

	loadErr := s.configSet.Load(s.valueContainer)
	s.Require().NoError(loadErr)

	value, getErr := s.configSet.Get("[-1]")

	s.Empty(value)
	s.Error(getErr)
}
