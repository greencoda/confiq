package confiq_test

import (
	"encoding/json"
	"net"
	"net/url"
	"testing"
	"time"

	"github.com/greencoda/confiq"
	"github.com/stretchr/testify/suite"
)

type CommonDecodersTestSuite struct {
	suite.Suite

	configSet *confiq.ConfigSet
}

func Test_CommonDecoders(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(CommonDecodersTestSuite))
}

func (s *CommonDecodersTestSuite) SetupTest() {
	s.configSet = confiq.New(
		confiq.WithTag("cfg"),
	)

	s.Require().NotNil(s.configSet)
}

func (s *CommonDecodersTestSuite) Test_Decode_Duration() {
	loadErr := s.configSet.LoadJSONFromFile("./testdata/common.json")
	s.Require().NoError(loadErr)

	type targetStruct struct {
		TestDuration time.Duration `cfg:"test_duration"`
	}

	var (
		target   targetStruct
		expected = 15 * time.Second
	)

	decodeErr := s.configSet.Decode(&target)

	s.Equal(expected, target.TestDuration)
	s.NoError(decodeErr)
}

func (s *CommonDecodersTestSuite) Test_Decode_Duration_FromInvalidFormat() {
	loadErr := s.configSet.LoadJSONFromFile("./testdata/common.json")
	s.Require().NoError(loadErr)

	type targetStruct struct {
		TestDuration time.Duration `cfg:"test_duration_invalid_format"`
	}

	var target targetStruct

	decodeErr := s.configSet.Decode(&target)

	s.Error(decodeErr)
}

func (s *CommonDecodersTestSuite) Test_Decode_IP() {
	loadErr := s.configSet.LoadJSONFromFile("./testdata/common.json")
	s.Require().NoError(loadErr)

	type targetStruct struct {
		TestIP net.IP `cfg:"test_ip"`
	}

	var (
		target   targetStruct
		expected = net.ParseIP("127.0.0.1")
	)

	decodeErr := s.configSet.Decode(&target)

	s.Equal(expected, target.TestIP)
	s.NoError(decodeErr)
}

func (s *CommonDecodersTestSuite) Test_Decode_IP_FromInvalidFormat() {
	loadErr := s.configSet.LoadJSONFromFile("./testdata/common.json")
	s.Require().NoError(loadErr)

	type targetStruct struct {
		TestIP net.IP `cfg:"test_ip_invalid_format"`
	}

	var target targetStruct

	decodeErr := s.configSet.Decode(&target)

	s.Error(decodeErr)
}

func (s *CommonDecodersTestSuite) Test_Decode_JSONRawMessage() {
	loadErr := s.configSet.LoadJSONFromFile("./testdata/common.json")
	s.Require().NoError(loadErr)

	type targetStruct struct {
		TestRawMessage json.RawMessage `cfg:"test_rawMessage"`
	}

	var (
		target   targetStruct
		expected = json.RawMessage(`{"rawMessage":"It's raw"}`)
	)

	decodeErr := s.configSet.Decode(&target)

	s.Equal(expected, target.TestRawMessage)
	s.NoError(decodeErr)
}

func (s *CommonDecodersTestSuite) Test_Decode_JSONRawMessage_WithUnmarshalable() {
	loadErr := s.configSet.LoadJSONFromFile("./testdata/common.json")
	s.Require().NoError(loadErr)

	type targetStruct struct {
		TestRawMessage json.RawMessage `cfg:"test_rawMessage"`
	}

	var target targetStruct

	s.configSet.OverrideValue(map[string]any{"test_rawMessage": map[bool]string{true: "It's raw"}})

	decodeErr := s.configSet.Decode(&target)

	s.Error(decodeErr)
}

func (s *CommonDecodersTestSuite) Test_Decode_URL() {
	loadErr := s.configSet.LoadJSONFromFile("./testdata/common.json")
	s.Require().NoError(loadErr)

	type targetStruct struct {
		TestURL *url.URL `cfg:"test_url"`
	}

	var (
		target   targetStruct
		expected = &url.URL{
			Scheme:      "http",
			Opaque:      "",
			User:        nil,
			Host:        "www.test.com",
			Path:        "",
			RawPath:     "",
			OmitHost:    false,
			ForceQuery:  false,
			RawQuery:    "",
			Fragment:    "",
			RawFragment: "",
		}
	)

	decodeErr := s.configSet.Decode(&target)

	s.Equal(expected, target.TestURL)
	s.NoError(decodeErr)
}

func (s *CommonDecodersTestSuite) Test_Decode_URL_FromInvalidFormat() {
	loadErr := s.configSet.LoadJSONFromFile("./testdata/common.json")
	s.Require().NoError(loadErr)

	type targetStruct struct {
		TestURL *url.URL `cfg:"test_url_invalid_format"`
	}

	var target targetStruct

	decodeErr := s.configSet.Decode(&target)

	s.Error(decodeErr)
}

func (s *CommonDecodersTestSuite) Test_Decode_Time() {
	loadErr := s.configSet.LoadJSONFromFile("./testdata/common.json")
	s.Require().NoError(loadErr)

	type targetStruct struct {
		TestTime time.Time `cfg:"test_time"`
	}

	var (
		target   targetStruct
		expected = time.Date(2025, 1, 13, 16, 0, 0, 0, time.FixedZone("", 9*60*60))
	)

	decodeErr := s.configSet.Decode(&target)

	s.Equal(expected, target.TestTime)
	s.NoError(decodeErr)
}

func (s *CommonDecodersTestSuite) Test_Decode_Time_FromTime() {
	loadErr := s.configSet.LoadTOMLFromFile("./testdata/common.toml")
	s.Require().NoError(loadErr)

	type targetStruct struct {
		TestTime time.Time `cfg:"test_time"`
	}

	var (
		target   targetStruct
		expected = time.Date(2025, 1, 13, 16, 0, 0, 0, time.FixedZone("", 9*60*60))
	)

	decodeErr := s.configSet.Decode(&target)

	s.Equal(expected, target.TestTime)
	s.NoError(decodeErr)
}

func (s *CommonDecodersTestSuite) Test_Decode_Time_FromInvalidType() {
	loadErr := s.configSet.LoadJSONFromFile("./testdata/common.json")
	s.Require().NoError(loadErr)

	type targetStruct struct {
		TestTime time.Time `cfg:"test_time_invalid_type"`
	}

	var target targetStruct

	decodeErr := s.configSet.Decode(&target)

	s.Error(decodeErr)
}

func (s *CommonDecodersTestSuite) Test_Decode_Time_FromInvalidFormat() {
	loadErr := s.configSet.LoadJSONFromFile("./testdata/common.json")
	s.Require().NoError(loadErr)

	type targetStruct struct {
		TestTime time.Time `cfg:"test_time_invalid_format"`
	}

	var target targetStruct

	decodeErr := s.configSet.Decode(&target)

	s.Error(decodeErr)
}
