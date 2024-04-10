package confiq_test

import (
	"testing"

	"github.com/greencoda/confiq"
	"github.com/stretchr/testify/suite"
)

type ConfigSetTestSuite struct {
	suite.Suite

	configSet *confiq.ConfigSet
}

func Test_ConfigSet(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(ConfigSetTestSuite))
}

func (s *ConfigSetTestSuite) SetupTest() {
	s.configSet = confiq.New(
		confiq.WithTag("cfg"),
	)

	s.Require().NotNil(s.T(), s.configSet)
}

func (s *ConfigSetTestSuite) Test_Decode() {
	loadErr := s.configSet.LoadJSONFromFile("./testdata/primitive.json")
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
	loadErr := s.configSet.LoadJSONFromFile("./testdata/primitive.json")
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
	loadErr := s.configSet.LoadJSONFromFile("./testdata/primitive.json")
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
	loadErr := s.configSet.LoadJSONFromFile("./testdata/primitive.json")
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
	loadErr := s.configSet.LoadJSONFromFile("./testdata/primitive.json")
	s.Require().NoError(loadErr)

	value, getErr := s.configSet.Get("test_string")

	s.Equal("test", value)
	s.NoError(getErr)
}

func (s *ConfigSetTestSuite) Test_Get_NestedPrimitives() {
	loadErr := s.configSet.LoadJSONFromFile("./testdata/composite.json")
	s.Require().NoError(loadErr)

	value, getErr := s.configSet.Get("test_section.test_string_array[0]")

	s.Equal("aleph", value)
	s.NoError(getErr)
}

func (s *ConfigSetTestSuite) Test_Get_FromMap_WithInvalidKeyPath() {
	loadErr := s.configSet.LoadJSONFromFile("./testdata/primitive.json")
	s.Require().NoError(loadErr)

	value, getErr := s.configSet.Get("noexistent_key")

	s.Empty(value)
	s.Error(getErr)
}

func (s *ConfigSetTestSuite) Test_Get_FromMap_WithInvalidIndexPath() {
	loadErr := s.configSet.LoadJSONFromFile("./testdata/primitive.json")
	s.Require().NoError(loadErr)

	value, getErr := s.configSet.Get("[0]")

	s.Empty(value)
	s.Error(getErr)
}

func (s *ConfigSetTestSuite) Test_Get_FromArray_WithInvalidKeyPath() {
	loadErr := s.configSet.LoadJSONFromFile("./testdata/array.json")
	s.Require().NoError(loadErr)

	value, getErr := s.configSet.Get("noexistent_key")

	s.Empty(value)
	s.Error(getErr)
}

func (s *ConfigSetTestSuite) Test_Get_FromArray_WithInvalidIndexPath() {
	loadErr := s.configSet.LoadJSONFromFile("./testdata/array.json")
	s.Require().NoError(loadErr)

	value, getErr := s.configSet.Get("[-1]")

	s.Empty(value)
	s.Error(getErr)
}
