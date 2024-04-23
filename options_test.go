package confiq_test

import (
	"testing"

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
		String string `cfg:"test_string"`
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
		String string `alt:"test_string"`
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
