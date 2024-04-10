package confiq_test

import (
	"testing"

	"github.com/greencoda/confiq"
	"github.com/stretchr/testify/suite"
)

type OptionsTestSuite struct {
	suite.Suite
}

func Test_OptionsTestSuite(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(OptionsTestSuite))
}

func (s *OptionsTestSuite) Test_WithTag() {
	type ConfigStruct struct {
		String string `cfg:"test_string"`
	}

	configSet := confiq.New(
		confiq.WithTag("alt"),
	)

	loadErr := configSet.LoadJSONFromFile("./testdata/primitive.json")
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

	loadErr := configSet.LoadJSONFromString(`{"test_string":"test"}`)
	s.Require().NoError(loadErr)

	var cfg ConfigStruct

	decodeErr := configSet.Decode(&cfg)
	s.NoError(decodeErr)
}
