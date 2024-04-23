package confiqtoml_test

import (
	"errors"
	"strings"
	"testing"

	confiqtoml "github.com/greencoda/confiq/loaders/toml"
	"github.com/stretchr/testify/suite"
)

var errFailedToRead = errors.New("failed to read")

type brokenReader struct{}

func (bR brokenReader) Read(_ []byte) (int, error) {
	return 0, errFailedToRead
}

type TOMLTestSuite struct {
	c *confiqtoml.Container

	suite.Suite
}

func Test_TOMLTestSuite(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(TOMLTestSuite))
}

func (s *TOMLTestSuite) SetupTest() {
	s.c = confiqtoml.Load()
	s.Require().NotNil(s.c)
}

func (s *TOMLTestSuite) Test_Get() {
	s.c.FromBytes([]byte("test_string = \"test\""))

	s.Require().Len(s.c.Get(), 1)
	s.Require().Empty(s.c.Errors())

	valueMap, ok := s.c.Get()[0].(map[string]any)
	s.Require().True(ok)

	s.Contains(valueMap, "test_string")
	s.Equal("test", valueMap["test_string"])
}

func (s *TOMLTestSuite) Test_FromFile() {
	s.c.FromFile("testdata/valid.toml")

	s.Len(s.c.Get(), 1)
	s.Empty(s.c.Errors())
}

func (s *TOMLTestSuite) Test_FromFile_InvalidPath() {
	s.c.FromFile("testdata/nonexistent.toml")

	s.Empty(s.c.Get())
	s.Len(s.c.Errors(), 1)
	s.ErrorIs(s.c.Errors()[0], confiqtoml.ErrCannotOpenTOMLFile)
}

func (s *TOMLTestSuite) Test_FromFile_Invalid() {
	s.c.FromFile("testdata/invalid.toml")

	s.Empty(s.c.Get())
	s.Len(s.c.Errors(), 1)
	s.ErrorIs(s.c.Errors()[0], confiqtoml.ErrCannotReadTOMLBytes)
}

func (s *TOMLTestSuite) Test_FromString() {
	s.c.FromString("test_string = \"test\"")

	s.Len(s.c.Get(), 1)
	s.Empty(s.c.Errors())
}

func (s *TOMLTestSuite) Test_FromString_Invalid() {
	s.c.FromString("{")

	s.Empty(s.c.Get())
	s.Len(s.c.Errors(), 1)
	s.ErrorIs(s.c.Errors()[0], confiqtoml.ErrCannotReadTOMLBytes)
}

func (s *TOMLTestSuite) Test_FromReader() {
	s.c.FromReader(strings.NewReader("test_string = \"test\""))

	s.Len(s.c.Get(), 1)
	s.Empty(s.c.Errors())
}

func (s *TOMLTestSuite) Test_FromReader_Invalid() {
	s.c.FromReader(strings.NewReader("["))

	s.Empty(s.c.Get())
	s.Len(s.c.Errors(), 1)
	s.ErrorIs(s.c.Errors()[0], confiqtoml.ErrCannotReadTOMLBytes)
}

func (s *TOMLTestSuite) Test_FromReader_Nil() {
	s.c.FromReader(nil)

	s.Empty(s.c.Get())
	s.Len(s.c.Errors(), 1)
	s.ErrorIs(s.c.Errors()[0], confiqtoml.ErrCannotReadTOMLData)
}

func (s *TOMLTestSuite) Test_FromReader_BrokenReader() {
	s.c.FromReader(brokenReader{})

	s.Empty(s.c.Get())
	s.Len(s.c.Errors(), 1)
	s.ErrorIs(s.c.Errors()[0], confiqtoml.ErrCannotReadTOMLData)
}

func (s *TOMLTestSuite) Test_FromBytes() {
	s.c.FromBytes([]byte("test_string = \"test\""))

	s.Len(s.c.Get(), 1)
	s.Empty(s.c.Errors())
}

func (s *TOMLTestSuite) Test_FromBytes_Invalid() {
	s.c.FromBytes([]byte("["))

	s.Empty(s.c.Get())
	s.Len(s.c.Errors(), 1)
	s.ErrorIs(s.c.Errors()[0], confiqtoml.ErrCannotReadTOMLBytes)
}
