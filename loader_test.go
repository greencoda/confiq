package confiq_test

import (
	"errors"
	"testing"

	"github.com/greencoda/confiq"
	"github.com/greencoda/confiq/mocks"
	"github.com/stretchr/testify/suite"
)

var errLoader = errors.New("loader error")

type LoaderTestSuite struct {
	suite.Suite

	configSet       *confiq.ConfigSet
	valueContainer1 *mocks.IValueContainer
	valueContainer2 *mocks.IValueContainer
}

func Test_LoaderTestSuite(t *testing.T) {
	suite.Run(t, new(LoaderTestSuite))
}

func (s *LoaderTestSuite) SetupTest() {
	s.configSet = confiq.New(
		confiq.WithTag("cfg"),
	)
	s.valueContainer1 = mocks.NewIValueContainer(s.T())
	s.valueContainer2 = mocks.NewIValueContainer(s.T())
}

func (s *LoaderTestSuite) Test_Load_WithErrors() {
	s.configSet.OverrideValue("test")

	s.valueContainer1.On("Errors").Return([]error{errLoader})

	loadErr := s.configSet.Load(s.valueContainer1)

	s.ErrorIs(loadErr, errLoader)
}

func (s *LoaderTestSuite) Test_LoadMaps_WithInvalidValue_NoPrefix() {
	s.configSet.OverrideValue("test")

	s.valueContainer1.On("Errors").Return([]error{})
	s.valueContainer1.On("Get").Return([]any{map[string]any{"a": "test_value 1"}})

	loadErr := s.configSet.Load(s.valueContainer1)

	s.Error(loadErr)
}

func (s *LoaderTestSuite) Test_LoadMaps_WithInvalidValue_WithPrefix() {
	s.configSet.OverrideValue("test")

	s.valueContainer1.On("Errors").Return([]error{})
	s.valueContainer1.On("Get").Return([]any{map[string]any{"a": "test_value 1"}})

	loadErr := s.configSet.Load(s.valueContainer1, confiq.WithPrefix("prefix"))
	s.Error(loadErr)
}

func (s *LoaderTestSuite) Test_LoadMaps_FollowedWithDifferentType() {
	s.valueContainer1.On("Errors").Return([]error{})
	s.valueContainer1.On("Get").Once().Return([]any{map[string]any{"a": "test_value 1"}})

	loadErr := s.configSet.Load(s.valueContainer1, confiq.WithPrefix("prefix"))
	s.Require().NoError(loadErr)

	s.valueContainer1.On("Errors").Return([]error{})
	s.valueContainer1.On("Get").Once().Return([]any{[]any{"a", "b"}})

	loadErr = s.configSet.Load(s.valueContainer1, confiq.WithPrefix("prefix"))
	s.Error(loadErr)
}

func (s *LoaderTestSuite) Test_LoadMaps_NoPrefix() {
	s.valueContainer1.On("Errors").Return([]error{})
	s.valueContainer1.On("Get").Once().Return([]any{
		map[string]any{"a": "test_value 1"},
	})
	s.valueContainer2.On("Errors").Return([]error{})
	s.valueContainer2.On("Get").Once().Return([]any{
		map[string]any{"b": "test_value 2"},
	})

	loadErr1 := s.configSet.Load(s.valueContainer1)
	s.Require().NoError(loadErr1)

	loadErr2 := s.configSet.Load(s.valueContainer2)
	s.Require().NoError(loadErr2)

	value, getErr := s.configSet.Get("")

	s.Equal(map[string]any{"a": "test_value 1", "b": "test_value 2"}, value)
	s.NoError(getErr)
}

func (s *LoaderTestSuite) Test_LoadMaps_WithSamePrefix() {
	s.valueContainer1.On("Errors").Return([]error{})
	s.valueContainer1.On("Get").Once().Return([]any{
		map[string]any{"a": "test_value 1"},
	})
	s.valueContainer2.On("Errors").Return([]error{})
	s.valueContainer2.On("Get").Once().Return([]any{
		map[string]any{"b": "test_value 2"},
	})

	loadErr1 := s.configSet.Load(s.valueContainer1, confiq.WithPrefix("prefix"))
	s.Require().NoError(loadErr1)

	loadErr2 := s.configSet.Load(s.valueContainer2, confiq.WithPrefix("prefix"))
	s.Require().NoError(loadErr2)

	value, getErr := s.configSet.Get("")

	s.Equal(map[string]any{"prefix": map[string]any{"a": "test_value 1", "b": "test_value 2"}}, value)
	s.NoError(getErr)
}

func (s *LoaderTestSuite) Test_LoadMaps_WithDifferentPrefix() {
	s.valueContainer1.On("Errors").Return([]error{})
	s.valueContainer1.On("Get").Once().Return([]any{
		map[string]any{"a": "test_value 1"},
	})
	s.valueContainer2.On("Errors").Return([]error{})
	s.valueContainer2.On("Get").Once().Return([]any{
		map[string]any{"b": "test_value 2"},
	})

	loadErr1 := s.configSet.Load(s.valueContainer1, confiq.WithPrefix("prefix_a"))
	s.Require().NoError(loadErr1)

	loadErr2 := s.configSet.Load(s.valueContainer2, confiq.WithPrefix("prefix_b"))
	s.Require().NoError(loadErr2)

	value, getErr := s.configSet.Get("")

	s.Equal(map[string]any{"prefix_a": map[string]any{"a": "test_value 1"}, "prefix_b": map[string]any{"b": "test_value 2"}}, value)
	s.NoError(getErr)
}

func (s *LoaderTestSuite) Test_LoadSlice_WithInvalidValue_NoPrefix() {
	s.configSet.OverrideValue("test")

	s.valueContainer1.On("Errors").Return([]error{})
	s.valueContainer1.On("Get").Once().Return([]any{
		[]any{"a", "b"},
	})

	loadErr := s.configSet.Load(s.valueContainer1)

	s.Error(loadErr)
}

func (s *LoaderTestSuite) Test_LoadSlice_WithInvalidValue_WithPrefix() {
	s.configSet.OverrideValue("test")

	s.valueContainer1.On("Errors").Return([]error{})
	s.valueContainer1.On("Get").Once().Return([]any{
		[]any{"a", "b"},
	})

	loadErr := s.configSet.Load(s.valueContainer1, confiq.WithPrefix("prefix"))

	s.Error(loadErr)
}

func (s *LoaderTestSuite) Test_LoadSlice_FollowedWithDifferentType() {
	s.valueContainer1.On("Errors").Return([]error{})
	s.valueContainer1.On("Get").Once().Return([]any{
		[]any{"a", "b"},
	})

	loadErr := s.configSet.Load(s.valueContainer1, confiq.WithPrefix("prefix"))
	s.Require().NoError(loadErr)

	s.valueContainer1.On("Errors").Return([]error{})
	s.valueContainer1.On("Get").Once().Return([]any{
		map[string]any{"a": "test_value 1"},
	})

	loadErr = s.configSet.Load(s.valueContainer1, confiq.WithPrefix("prefix"))
	s.Error(loadErr)
}

func (s *LoaderTestSuite) Test_LoadSlices_NoPrefix() {
	s.valueContainer1.On("Errors").Return([]error{})
	s.valueContainer1.On("Get").Once().Return([]any{
		[]any{"a", "b"},
	})
	s.valueContainer2.On("Errors").Return([]error{})
	s.valueContainer2.On("Get").Once().Return([]any{
		[]any{"c", "d"},
	})

	loadErr1 := s.configSet.Load(s.valueContainer1)
	loadErr2 := s.configSet.Load(s.valueContainer2)

	s.Require().NoError(loadErr1)
	s.Require().NoError(loadErr2)

	value, getErr := s.configSet.Get("")

	s.Equal([]any{"a", "b", "c", "d"}, value)
	s.NoError(getErr)
}

func (s *LoaderTestSuite) Test_LoadSlices_WithSamePrefix() {
	s.valueContainer1.On("Errors").Return([]error{})
	s.valueContainer1.On("Get").Once().Return([]any{
		[]any{"a", "b"},
	})
	s.valueContainer2.On("Errors").Return([]error{})
	s.valueContainer2.On("Get").Once().Return([]any{
		[]any{"c", "d"},
	})

	loadErr1 := s.configSet.Load(s.valueContainer1, confiq.WithPrefix("prefix"))
	s.Require().NoError(loadErr1)

	loadErr2 := s.configSet.Load(s.valueContainer2, confiq.WithPrefix("prefix"))
	s.Require().NoError(loadErr2)

	value, getErr := s.configSet.Get("")

	s.Equal(map[string]any{"prefix": []any{"a", "b", "c", "d"}}, value)
	s.NoError(getErr)
}

func (s *LoaderTestSuite) Test_LoadSlices_WithDifferentPrefix() {
	s.valueContainer1.On("Errors").Return([]error{})
	s.valueContainer1.On("Get").Once().Return([]any{
		[]any{"a", "b"},
	})
	s.valueContainer2.On("Errors").Return([]error{})
	s.valueContainer2.On("Get").Once().Return([]any{
		[]any{"c", "d"},
	})

	loadErr1 := s.configSet.Load(s.valueContainer1, confiq.WithPrefix("prefix_a"))
	s.Require().NoError(loadErr1)

	loadErr2 := s.configSet.Load(s.valueContainer2, confiq.WithPrefix("prefix_b"))
	s.Require().NoError(loadErr2)

	value, getErr := s.configSet.Get("")

	s.Equal(map[string]any{"prefix_a": []any{"a", "b"}, "prefix_b": []any{"c", "d"}}, value)
	s.NoError(getErr)
}

func (s *LoaderTestSuite) Test_LoadMaps_WithInvalidConfigValue() {
	loadErr := s.configSet.LoadRawValue([]any{false})

	s.Error(loadErr)
}

func (s *LoaderTestSuite) Test_LoadMaps_WithNilConfigValue() {
	loadErr := s.configSet.LoadRawValue(nil)

	s.Error(loadErr)
}
