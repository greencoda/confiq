package confiq_test

import (
	"testing"

	"github.com/greencoda/confiq"
	"github.com/stretchr/testify/suite"
)

type LoadEnvTestSuite struct {
	suite.Suite

	configSet *confiq.ConfigSet
}

func Test_LoadEnvTestSuite(t *testing.T) {
	suite.Run(t, new(LoadEnvTestSuite))
}

func (s *LoadEnvTestSuite) SetupTest() {
	s.configSet = confiq.New(
		confiq.WithTag("cfg"),
	)
}

func (s *LoadEnvTestSuite) Test_LoadEnvFromEnv() {
	s.T().Setenv("TEST_VALUE", "test")

	err := s.configSet.LoadEnvFromEnvironment()
	s.Require().NoError(err)

	value, err := s.configSet.Get("TEST_VALUE")

	s.Equal("test", value)
	s.NoError(err)
}

func (s *LoadEnvTestSuite) Test_LoadEnvFromEnv_WithPrefix() {
	s.T().Setenv("TEST_VALUE", "test")

	err := s.configSet.LoadEnvFromEnvironment(confiq.WithPrefix("testKey"))
	s.Require().NoError(err)

	value, err := s.configSet.Get("testKey.TEST_VALUE")

	s.Equal("test", value)
	s.NoError(err)
}

func (s *LoadEnvTestSuite) Test_LoadEnvFromEnv_Empty() {
	err := s.configSet.LoadEnvFromEnvironment()
	s.Require().NoError(err)

	value, err := s.configSet.Get("TEST_VALUE")

	s.Nil(value)
	s.Error(err)
}

func (s *LoadEnvTestSuite) Test_LoadEnvFromFile_InvalidPath() {
	err := s.configSet.LoadEnvFromFile("testdata/nonexistent.env")
	s.ErrorIs(err, confiq.ErrCannotOpenConfig)
}

func (s *LoadEnvTestSuite) Test_LoadEnvFromFile_Invalid() {
	err := s.configSet.LoadEnvFromFile("testdata/invalid.env")
	s.ErrorIs(err, confiq.ErrCannotLoadConfig)
}

func (s *LoadEnvTestSuite) Test_LoadEnvFromFile_() {
	err := s.configSet.LoadEnvFromFile("testdata/primitive.env")
	s.NoError(err)
}

func (s *LoadEnvTestSuite) Test_LoadEnvFromString_Invalid() {
	err := s.configSet.LoadEnvFromString("TEST_BOOL")
	s.ErrorIs(err, confiq.ErrCannotLoadConfig)
}

func (s *LoadEnvTestSuite) Test_LoadEnvFromString() {
	err := s.configSet.LoadEnvFromString("TEST_BOOL=true")
	s.NoError(err)
}

func (s *LoadEnvTestSuite) Test_LoadEnvFromBytes_Invalid() {
	err := s.configSet.LoadEnvFromBytes([]byte("TEST_BOOL"))
	s.ErrorIs(err, confiq.ErrCannotLoadConfig)
}

func (s *LoadEnvTestSuite) Test_LoadEnvFromBytes() {
	err := s.configSet.LoadEnvFromBytes([]byte("TEST_BOOL=true"))
	s.NoError(err)
}
