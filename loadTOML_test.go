package confiq_test

import (
	"strings"
	"testing"

	"github.com/greencoda/confiq"
	"github.com/stretchr/testify/suite"
)

type LoadTOMLTestSuite struct {
	suite.Suite

	configSet *confiq.ConfigSet
}

func Test_LoadTOMLTestSuite(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(LoadTOMLTestSuite))
}

func (s *LoadTOMLTestSuite) SetupTest() {
	s.configSet = confiq.New(
		confiq.WithTag("cfg"),
	)
}

func (s *LoadTOMLTestSuite) Test_LoadTOMLFromFile() {
	err := s.configSet.LoadTOMLFromFile("testdata/primitive.toml")
	s.NoError(err)
}

func (s *LoadTOMLTestSuite) Test_LoadTOMLFromFile_WithPrefix() {
	err := s.configSet.LoadTOMLFromFile("testdata/primitive.toml", confiq.WithPrefix("prefix"))
	s.Require().NoError(err)

	value, err := s.configSet.Get("prefix")

	s.NotNil(value)
	s.NoError(err)
}

func (s *LoadTOMLTestSuite) Test_LoadTOMLFromFile_InvalidPath() {
	err := s.configSet.LoadTOMLFromFile("testdata/nonexistent.toml")
	s.ErrorIs(err, confiq.ErrCannotOpenConfig)
}

func (s *LoadTOMLTestSuite) Test_LoadTOMLFromFile_Invalid() {
	err := s.configSet.LoadTOMLFromFile("testdata/invalid.toml")
	s.ErrorIs(err, confiq.ErrCannotLoadConfig)
}

// LoadTOMLFromString.
func (s *LoadTOMLTestSuite) Test_LoadTOMLFromString() {
	err := s.configSet.LoadTOMLFromString("test_bool = true")
	s.NoError(err)
}

func (s *LoadTOMLTestSuite) Test_LoadTOMLFromString_Invalid() {
	err := s.configSet.LoadTOMLFromString("[")
	s.ErrorIs(err, confiq.ErrCannotLoadConfig)
}

// LoadTOMLFromReader.
func (s *LoadJSONTestSuite) Test_LoadTOMLFromReader() {
	err := s.configSet.LoadTOMLFromReader(strings.NewReader("test_bool = true"))
	s.NoError(err)
}

func (s *LoadJSONTestSuite) Test_LoadTOMLFromReader_Invalid() {
	err := s.configSet.LoadTOMLFromReader(strings.NewReader("["))
	s.ErrorIs(err, confiq.ErrCannotLoadConfig)
}

func (s *LoadJSONTestSuite) Test_LoadTOMLFromReader_Nil() {
	err := s.configSet.LoadTOMLFromReader(nil)
	s.ErrorIs(err, confiq.ErrCannotLoadConfig)
}

func (s *LoadJSONTestSuite) Test_LoadTOMLFromReader_BrokenReader() {
	err := s.configSet.LoadTOMLFromReader(brokenReader{})
	s.ErrorIs(err, confiq.ErrCannotLoadConfig)
}

// LoadTOMLFromBytes.
func (s *LoadTOMLTestSuite) Test_LoadTOMLFromBytes() {
	err := s.configSet.LoadTOMLFromBytes([]byte("test_bool = true"))
	s.NoError(err)
}

func (s *LoadTOMLTestSuite) Test_LoadTOMLFromBytes_Invalid() {
	err := s.configSet.LoadTOMLFromBytes([]byte("["))
	s.ErrorIs(err, confiq.ErrCannotLoadConfig)
}
