package confiq_test

import (
	"errors"
	"testing"

	"github.com/greencoda/confiq"
	"github.com/stretchr/testify/suite"
)

var errFailedToRead = errors.New("failed to read")

type brokenReader struct{}

func (br brokenReader) Read(_ []byte) (int, error) {
	return 0, errFailedToRead
}

type LoaderTestSuite struct {
	suite.Suite

	configSet *confiq.ConfigSet
}

func Test_LoaderTestSuite(t *testing.T) {
	suite.Run(t, new(LoaderTestSuite))
}

func (s *LoaderTestSuite) SetupTest() {
	s.configSet = confiq.New(
		confiq.WithTag("cfg"),
	)
}

func (s *LoaderTestSuite) Test_LoadMaps_WithInvalidValue_NoPrefix() {
	s.configSet.OverrideValue("test")

	loadErr := s.configSet.LoadJSONFromString(`{"a":"test_value 1"}`)

	s.Error(loadErr)
}

func (s *LoaderTestSuite) Test_LoadMaps_WithInvalidValue_WithPrefix() {
	s.configSet.OverrideValue("test")

	loadErr := s.configSet.LoadJSONFromString(`{"a":"test_value 1"}`, confiq.WithPrefix("prefix"))

	s.Error(loadErr)
}

func (s *LoaderTestSuite) Test_LoadMaps_FollowedWithDifferentType() {
	loadErr := s.configSet.LoadJSONFromString(`{"a":"test_value 1"}`, confiq.WithPrefix("prefix"))
	s.Require().NoError(loadErr)

	loadErr = s.configSet.LoadJSONFromString(`["a", "b"]`, confiq.WithPrefix("prefix"))
	s.Error(loadErr)
}

func (s *LoaderTestSuite) Test_LoadMaps_FromSameFormat_NoPrefix() {
	loadErr1 := s.configSet.LoadJSONFromString(`{"a":"test_value 1"}`)
	loadErr2 := s.configSet.LoadJSONFromString(`{"b":"test_value 2"}`)

	s.Require().NoError(loadErr1)
	s.Require().NoError(loadErr2)

	value, getErr := s.configSet.Get("")

	s.Equal(map[string]any{"a": "test_value 1", "b": "test_value 2"}, value)
	s.NoError(getErr)
}

func (s *LoaderTestSuite) Test_LoadMaps_FromDifferentFormat_NoPrefix() {
	loadErr1 := s.configSet.LoadJSONFromString(`{"a":"test_value 1"}`)
	loadErr2 := s.configSet.LoadTOMLFromString(`b = "test_value 2"`)

	s.Require().NoError(loadErr1)
	s.Require().NoError(loadErr2)

	value, getErr := s.configSet.Get("")

	s.Equal(map[string]any{"a": "test_value 1", "b": "test_value 2"}, value)
	s.NoError(getErr)
}

func (s *LoaderTestSuite) Test_LoadMaps_FromSameFormat_WithSamePrefix() {
	loadErr1 := s.configSet.LoadJSONFromString(`{"a":"test_value 1"}`, confiq.WithPrefix("prefix"))
	loadErr2 := s.configSet.LoadJSONFromString(`{"b":"test_value 2"}`, confiq.WithPrefix("prefix"))

	s.Require().NoError(loadErr1)
	s.Require().NoError(loadErr2)

	value, getErr := s.configSet.Get("")

	s.Equal(map[string]any{"prefix": map[string]any{"a": "test_value 1", "b": "test_value 2"}}, value)
	s.NoError(getErr)
}

func (s *LoaderTestSuite) Test_LoadMaps_FromSameFormat_WithDifferentPrefix() {
	loadErr1 := s.configSet.LoadJSONFromString(`{"a":"test_value 1"}`, confiq.WithPrefix("prefix_a"))
	loadErr2 := s.configSet.LoadJSONFromString(`{"b":"test_value 2"}`, confiq.WithPrefix("prefix_b"))

	s.Require().NoError(loadErr1)
	s.Require().NoError(loadErr2)

	value, getErr := s.configSet.Get("")

	s.Equal(map[string]any{"prefix_a": map[string]any{"a": "test_value 1"}, "prefix_b": map[string]any{"b": "test_value 2"}}, value)
	s.NoError(getErr)
}

func (s *LoaderTestSuite) Test_LoadSlice_WithInvalidValue_NoPrefix() {
	s.configSet.OverrideValue("test")

	loadErr := s.configSet.LoadJSONFromString(`["a", "b"]`)

	s.Error(loadErr)
}

func (s *LoaderTestSuite) Test_LoadSlice_WithInvalidValue_WithPrefix() {
	s.configSet.OverrideValue("test")

	loadErr := s.configSet.LoadJSONFromString(`["a", "b"]`, confiq.WithPrefix("prefix"))

	s.Error(loadErr)
}

func (s *LoaderTestSuite) Test_LoadSlice_FollowedWithDifferentType() {
	loadErr := s.configSet.LoadJSONFromString(`["a", "b"]`, confiq.WithPrefix("prefix"))
	s.Require().NoError(loadErr)

	loadErr = s.configSet.LoadJSONFromString(`{"a":"test_value 1"}`, confiq.WithPrefix("prefix"))
	s.Error(loadErr)
}

func (s *LoaderTestSuite) Test_LoadSlices_FromSameFormat_NoPrefix() {
	loadErr1 := s.configSet.LoadJSONFromString(`["a", "b"]`)
	loadErr2 := s.configSet.LoadJSONFromString(`["c", "d"]`)

	s.Require().NoError(loadErr1)
	s.Require().NoError(loadErr2)

	value, getErr := s.configSet.Get("")

	s.Equal([]any{"a", "b", "c", "d"}, value)
	s.NoError(getErr)
}

func (s *LoaderTestSuite) Test_LoadSlices_FromSameFormat_WithSamePrefix() {
	loadErr1 := s.configSet.LoadJSONFromString(`["a", "b"]`, confiq.WithPrefix("prefix"))
	loadErr2 := s.configSet.LoadJSONFromString(`["c", "d"]`, confiq.WithPrefix("prefix"))

	s.Require().NoError(loadErr1)
	s.Require().NoError(loadErr2)

	value, getErr := s.configSet.Get("")

	s.Equal(map[string]any{"prefix": []any{"a", "b", "c", "d"}}, value)
	s.NoError(getErr)
}

func (s *LoaderTestSuite) Test_LoadSlices_FromSameFormat_WithDifferentPrefix() {
	loadErr1 := s.configSet.LoadJSONFromString(`["a", "b"]`, confiq.WithPrefix("prefix_a"))
	loadErr2 := s.configSet.LoadJSONFromString(`["c", "d"]`, confiq.WithPrefix("prefix_b"))

	s.Require().NoError(loadErr1)
	s.Require().NoError(loadErr2)

	value, getErr := s.configSet.Get("")

	s.Equal(map[string]any{"prefix_a": []any{"a", "b"}, "prefix_b": []any{"c", "d"}}, value)
	s.NoError(getErr)
}

func (s *LoaderTestSuite) Test_LoadSlices_FromDifferentFormat_NoPrefix() {
	loadErr1 := s.configSet.LoadJSONFromString(`["a", "b"]`)
	loadErr2 := s.configSet.LoadYAMLFromString("---\n- c\n- d")

	s.Require().NoError(loadErr1)
	s.Require().NoError(loadErr2)

	value, getErr := s.configSet.Get("")

	s.Equal([]any{"a", "b", "c", "d"}, value)
	s.NoError(getErr)
}
