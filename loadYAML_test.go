package confiq_test

import (
	"strings"
	"testing"

	"github.com/greencoda/confiq"
	"github.com/stretchr/testify/suite"
)

type LoadYAMLTestSuite struct {
	suite.Suite

	configSet *confiq.ConfigSet
}

func Test_LoadYAMLTestSuite(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(LoadYAMLTestSuite))
}

func (s *LoadYAMLTestSuite) SetupTest() {
	s.configSet = confiq.New(
		confiq.WithTag("cfg"),
	)
}

// LoadYAMLFromFile.
func (s *LoadYAMLTestSuite) Test_LoadYAMLFromFile() {
	err := s.configSet.LoadYAMLFromFile("testdata/primitive.yaml")
	s.NoError(err)
}

func (s *LoadYAMLTestSuite) Test_LoadYAMLFromFile_WithPrefix() {
	err := s.configSet.LoadYAMLFromFile("testdata/primitive.yaml", confiq.WithPrefix("prefix"))
	s.Require().NoError(err)

	value, err := s.configSet.Get("prefix")

	s.NotNil(value)
	s.NoError(err)
}

func (s *LoadYAMLTestSuite) Test_LoadYAMLFromFile_InvalidPath() {
	err := s.configSet.LoadYAMLFromFile("testdata/nonexistent.yaml")
	s.ErrorIs(err, confiq.ErrCannotOpenConfig)
}

func (s *LoadYAMLTestSuite) Test_LoadYAMLFromFile_Invalid() {
	err := s.configSet.LoadYAMLFromFile("testdata/invalid.yaml")
	s.ErrorIs(err, confiq.ErrCannotLoadConfig)
}

// LoadYAMLFromString.
func (s *LoadYAMLTestSuite) Test_LoadYAMLFromString() {
	err := s.configSet.LoadYAMLFromString("test_bool: true")
	s.NoError(err)
}

func (s *LoadYAMLTestSuite) Test_LoadYAMLFromString_Invalid() {
	err := s.configSet.LoadYAMLFromString("{test")
	s.ErrorIs(err, confiq.ErrCannotLoadConfig)
}

// LoadYAMLFromReader.
func (s *LoadJSONTestSuite) Test_LoadYAMLFromReader() {
	err := s.configSet.LoadYAMLFromReader(strings.NewReader("test_bool: true"))
	s.NoError(err)
}

func (s *LoadJSONTestSuite) Test_LoadYAMLFromReader_Invalid() {
	err := s.configSet.LoadYAMLFromReader(strings.NewReader("{test"))
	s.ErrorIs(err, confiq.ErrCannotLoadConfig)
}

func (s *LoadJSONTestSuite) Test_LoadYAMLFromReader_Nil() {
	err := s.configSet.LoadYAMLFromReader(nil)
	s.ErrorIs(err, confiq.ErrCannotLoadConfig)
}

func (s *LoadJSONTestSuite) Test_LoadYAMLFromReader_BrokenReader() {
	err := s.configSet.LoadYAMLFromReader(brokenReader{})
	s.ErrorIs(err, confiq.ErrCannotLoadConfig)
}

// LoadYAMLFromBytes.
func (s *LoadYAMLTestSuite) Test_LoadYAMLFromBytes() {
	err := s.configSet.LoadYAMLFromBytes([]byte("test_bool: true"))
	s.NoError(err)
}

func (s *LoadYAMLTestSuite) Test_LoadYAMLFromBytes_Invalid() {
	err := s.configSet.LoadYAMLFromBytes([]byte("{test"))
	s.ErrorIs(err, confiq.ErrCannotLoadConfig)
}
