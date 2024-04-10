package confiq_test

import (
	"testing"

	"github.com/greencoda/confiq"
	"github.com/stretchr/testify/suite"
)

type LoadRawTestSuite struct {
	suite.Suite

	configSet *confiq.ConfigSet
}

func Test_LoadRawTestSuite(t *testing.T) {
	suite.Run(t, new(LoadRawTestSuite))
}

func (s *LoadRawTestSuite) SetupTest() {
	s.configSet = confiq.New(
		confiq.WithTag("cfg"),
	)
}

func (s *LoadRawTestSuite) Test_LoadMaps_WithInvalidConfigValue() {
	loadErr := s.configSet.LoadRawValue(false)

	s.Error(loadErr)
}

func (s *LoadRawTestSuite) Test_LoadMaps_WithNilConfigValue() {
	loadErr := s.configSet.LoadRawValue(nil)

	s.Error(loadErr)
}
