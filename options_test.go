package confiq_test

import (
	"testing"
	"time"

	"github.com/greencoda/confiq"
	"github.com/greencoda/confiq/mocks"
	"github.com/stretchr/testify/suite"
)

type OptionsTestSuite struct {
	suite.Suite

	valueContainer *mocks.IValueContainer
}

func Test_OptionsTestSuite(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(OptionsTestSuite))
}

func (s *OptionsTestSuite) SetupTest() {
	s.valueContainer = mocks.NewIValueContainer(s.T())
}

func (s *OptionsTestSuite) Test_WithTag() {
	type ConfigStruct struct {
		TestString string `cfg:"test_string"`
	}

	configSet := confiq.New(
		confiq.WithTag("alt"),
	)

	s.valueContainer.On("Errors").Return([]error{})
	s.valueContainer.On("Get").Return([]any{map[string]any{"test_string": "test"}})

	loadErr := configSet.Load(s.valueContainer)
	s.Require().NoError(loadErr)

	var cfg ConfigStruct

	decodeErr := configSet.Decode(&cfg)
	s.NoError(decodeErr)
}

func (s *OptionsTestSuite) Test_WithPrefix() {
	type ConfigStruct struct {
		TestString string `cfg:"test_string"`
	}

	configSet := confiq.New()

	s.valueContainer.On("Errors").Return([]error{})
	s.valueContainer.On("Get").Return([]any{map[string]any{"test_string": "test"}})

	loadErr := configSet.Load(s.valueContainer)
	s.Require().NoError(loadErr)

	var cfg ConfigStruct

	decodeErr := configSet.Decode(&cfg)
	s.NoError(decodeErr)
}

func (s *OptionsTestSuite) Test_AsStrict() {
	type ConfigStruct struct {
		TestString   string        `cfg:"test_string"`
		TestDuration time.Duration `cfg:"test_duration"`
	}

	configSet := confiq.New()

	s.valueContainer.On("Errors").Return([]error{})
	s.valueContainer.On("Get").Return([]any{map[string]any{"test_string": "test", "test_duration": "one hour"}})

	loadErr := configSet.Load(s.valueContainer)
	s.Require().NoError(loadErr)

	var cfg ConfigStruct

	decodeErr := configSet.Decode(&cfg, confiq.AsStrict())

	s.Error(decodeErr)
}

func (s *OptionsTestSuite) Test_FromPrefix() {
	type ConfigStruct struct {
		TestString string `cfg:"test_string"`
	}

	configSet := confiq.New()

	s.valueContainer.On("Errors").Return([]error{})
	s.valueContainer.On("Get").Return([]any{map[string]any{"test_struct": map[string]any{"test_string": "test"}}})

	loadErr := configSet.Load(s.valueContainer)
	s.Require().NoError(loadErr)

	var cfg ConfigStruct

	decodeErr := configSet.Decode(&cfg, confiq.FromPrefix("test_struct"))
	s.NoError(decodeErr)
}
