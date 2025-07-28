package confiqenv_test

import (
	"errors"
	"strings"
	"testing"

	confiqenv "github.com/greencoda/confiq/loaders/env"
	"github.com/stretchr/testify/suite"
)

var errFailedToRead = errors.New("failed to read")

type brokenReader struct{}

func (bR brokenReader) Read(_ []byte) (int, error) {
	return 0, errFailedToRead
}

type EnvTestSuite struct {
	suite.Suite

	c *confiqenv.Container
}

func Test_EnvTestSuite(t *testing.T) {
	suite.Run(t, new(EnvTestSuite))
}

func (s *EnvTestSuite) SetupTest() {
	s.c = confiqenv.Load()
	s.Require().NotNil(s.c)
}

func (s *EnvTestSuite) Test_FromEnvironment() {
	s.T().Setenv("TEST_VALUE", "test")

	s.c.FromEnvironment()

	s.Len(s.c.Get(), 1)

	valueMap, ok := s.c.Get()[0].(map[string]any)
	s.Require().True(ok)

	s.Contains(valueMap, "TEST_VALUE")
	s.Equal("test", valueMap["TEST_VALUE"])
	s.Empty(s.c.Errors())
}

func (s *EnvTestSuite) Test_FromEnvironment_Empty() {
	s.c.FromEnvironment()

	s.Len(s.c.Get(), 1)

	valueMap, ok := s.c.Get()[0].(map[string]any)
	s.Require().True(ok)

	s.NotContains(valueMap, "TEST_VALUE")
	s.Empty(s.c.Errors())
}

func (s *EnvTestSuite) Test_FromFile() {
	s.c.FromFile("testdata/valid.env")

	s.Len(s.c.Get(), 1)
	s.Empty(s.c.Errors())
}

func (s *EnvTestSuite) Test_FromFile_InvalidPath() {
	s.c.FromFile("testdata/nonexistent.env")

	s.Empty(s.c.Get())
	s.Len(s.c.Errors(), 1)
	s.ErrorIs(s.c.Errors()[0], confiqenv.ErrCannotOpenEnvFile)
}

func (s *EnvTestSuite) Test_FromFile_Invalid() {
	s.c.FromFile("testdata/invalid.env")

	s.Empty(s.c.Get())
	s.Len(s.c.Errors(), 1)
	s.ErrorIs(s.c.Errors()[0], confiqenv.ErrCannotReadEnvData)
}

func (s *EnvTestSuite) Test_FromString() {
	s.c.FromString("TEST_BOOL=true")

	s.Len(s.c.Get(), 1)
	s.Empty(s.c.Errors())
}

func (s *EnvTestSuite) Test_FromString_Invalid() {
	s.c.FromString("TEST_BOOL")

	s.Empty(s.c.Get())
	s.Len(s.c.Errors(), 1)
	s.ErrorIs(s.c.Errors()[0], confiqenv.ErrCannotReadEnvData)
}

func (s *EnvTestSuite) Test_FromReader() {
	s.c.FromReader(strings.NewReader("TEST_BOOL=true"))

	s.Len(s.c.Get(), 1)
	s.Empty(s.c.Errors())
}

func (s *EnvTestSuite) Test_FromReader_Invalid() {
	s.c.FromReader(strings.NewReader("{"))

	s.Empty(s.c.Get())
	s.Len(s.c.Errors(), 1)
	s.ErrorIs(s.c.Errors()[0], confiqenv.ErrCannotReadEnvData)
}

func (s *EnvTestSuite) Test_FromReader_Nil() {
	s.c.FromReader(nil)

	s.Empty(s.c.Get())
	s.Len(s.c.Errors(), 1)
	s.ErrorIs(s.c.Errors()[0], confiqenv.ErrCannotReadEnvData)
}

func (s *EnvTestSuite) Test_FromReader_BrokenReader() {
	s.c.FromReader(brokenReader{})

	s.Empty(s.c.Get())
	s.Len(s.c.Errors(), 1)
	s.ErrorIs(s.c.Errors()[0], confiqenv.ErrCannotReadEnvData)
}

func (s *EnvTestSuite) Test_FromBytes() {
	s.c.FromBytes([]byte("TEST_BOOL=true"))

	s.Len(s.c.Get(), 1)
	s.Empty(s.c.Errors())
}

func (s *EnvTestSuite) Test_FromBytes_Invalid() {
	s.c.FromBytes([]byte("TEST_BOOL"))

	s.Empty(s.c.Get())
	s.Len(s.c.Errors(), 1)
	s.ErrorIs(s.c.Errors()[0], confiqenv.ErrCannotReadEnvData)
}
