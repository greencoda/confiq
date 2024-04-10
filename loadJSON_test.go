package confiq_test

import (
	"strings"
	"testing"

	"github.com/greencoda/confiq"
	"github.com/stretchr/testify/suite"
)

type LoadJSONTestSuite struct {
	suite.Suite

	configSet *confiq.ConfigSet
}

func Test_LoadJSONTestSuite(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(LoadJSONTestSuite))
}

func (s *LoadJSONTestSuite) SetupTest() {
	s.configSet = confiq.New(
		confiq.WithTag("cfg"),
	)
}

func (s *LoadJSONTestSuite) Test_LoadJSONFromFile() {
	err := s.configSet.LoadJSONFromFile("testdata/primitive.json")
	s.NoError(err)
}

func (s *LoadJSONTestSuite) Test_LoadJSONFromFile_WithPrefix() {
	err := s.configSet.LoadJSONFromFile("testdata/primitive.json", confiq.WithPrefix("prefix"))
	s.Require().NoError(err)

	value, err := s.configSet.Get("prefix")

	s.NotNil(value)
	s.NoError(err)
}

// LoadJSONFromFile.
func (s *LoadJSONTestSuite) Test_LoadJSONFromFile_InvalidPath() {
	err := s.configSet.LoadJSONFromFile("testdata/nonexistent.json")
	s.ErrorIs(err, confiq.ErrCannotOpenConfig)
}

func (s *LoadJSONTestSuite) Test_LoadJSONFromFile_Invalid() {
	err := s.configSet.LoadJSONFromFile("testdata/invalid.json")
	s.ErrorIs(err, confiq.ErrCannotLoadConfig)
}

// LoadJSONFromString.
func (s *LoadJSONTestSuite) Test_LoadJSONFromString() {
	err := s.configSet.LoadJSONFromString("{}")
	s.NoError(err)
}

func (s *LoadJSONTestSuite) Test_LoadJSONFromString_Invalid() {
	err := s.configSet.LoadJSONFromString("{")
	s.ErrorIs(err, confiq.ErrCannotLoadConfig)
}

// LoadJSONFromReader.
func (s *LoadJSONTestSuite) Test_LoadJSONFromReader() {
	err := s.configSet.LoadJSONFromReader(strings.NewReader("{}"))
	s.NoError(err)
}

func (s *LoadJSONTestSuite) Test_LoadJSONFromReader_Invalid() {
	err := s.configSet.LoadJSONFromReader(strings.NewReader("{"))
	s.ErrorIs(err, confiq.ErrCannotLoadConfig)
}

func (s *LoadJSONTestSuite) Test_LoadJSONFromReader_Nil() {
	err := s.configSet.LoadJSONFromReader(nil)
	s.ErrorIs(err, confiq.ErrCannotLoadConfig)
}

func (s *LoadJSONTestSuite) Test_LoadJSONFromReader_BrokenReader() {
	err := s.configSet.LoadJSONFromReader(brokenReader{})
	s.ErrorIs(err, confiq.ErrCannotLoadConfig)
}

// LoadJSONFromBytes.
func (s *LoadJSONTestSuite) Test_LoadJSONFromBytes() {
	err := s.configSet.LoadJSONFromBytes([]byte("{}"))
	s.NoError(err)
}

func (s *LoadJSONTestSuite) Test_LoadJSONFromBytes_Invalid() {
	err := s.configSet.LoadJSONFromBytes([]byte("{"))
	s.ErrorIs(err, confiq.ErrCannotLoadConfig)
}
