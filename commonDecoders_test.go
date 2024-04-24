package confiq_test

import (
	"encoding/json"
	"net"
	"net/url"
	"testing"
	"time"

	"github.com/greencoda/confiq"
	"github.com/greencoda/confiq/mocks"
	"github.com/stretchr/testify/suite"
)

type CommonDecodersTestSuite struct {
	suite.Suite

	configSet      *confiq.ConfigSet
	valueContainer *mocks.IValueContainer
}

func Test_CommonDecoders(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(CommonDecodersTestSuite))
}

func (s *CommonDecodersTestSuite) SetupTest() {
	s.configSet = confiq.New(
		confiq.WithTag("cfg"),
	)
	s.valueContainer = mocks.NewIValueContainer(s.T())

	s.Require().NotNil(s.configSet)
}

func (s *CommonDecodersTestSuite) Test_Decode_Duration() {
	s.valueContainer.On("Errors").Return([]error{})
	s.valueContainer.On("Get").Return([]any{map[string]any{"test_duration": "15s"}})

	loadErr := s.configSet.Load(s.valueContainer)
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

func (s *CommonDecodersTestSuite) Test_Decode_Duration_FromNil() {
	s.valueContainer.On("Errors").Return([]error{})
	s.valueContainer.On("Get").Return([]any{map[string]any{"test_duration": nil}})

	loadErr := s.configSet.Load(s.valueContainer)
	s.Require().NoError(loadErr)

	type targetStruct struct {
		TestDuration time.Duration `cfg:"test_duration"`
	}

	var target targetStruct

	decodeErr := s.configSet.Decode(&target)

	s.Error(decodeErr)
}

func (s *CommonDecodersTestSuite) Test_Decode_Duration_FromInvalidFormat() {
	s.valueContainer.On("Errors").Return([]error{})
	s.valueContainer.On("Get").Return([]any{map[string]any{"test_duration_invalid_format": "fifteen seconds"}})

	loadErr := s.configSet.Load(s.valueContainer)
	s.Require().NoError(loadErr)

	type targetStruct struct {
		TestDuration time.Duration `cfg:"test_duration_invalid_format"`
	}

	var target targetStruct

	decodeErr := s.configSet.Decode(&target)

	s.Error(decodeErr)
}

func (s *CommonDecodersTestSuite) Test_Decode_IP() {
	s.valueContainer.On("Errors").Return([]error{})
	s.valueContainer.On("Get").Return([]any{map[string]any{"test_ip": "127.0.0.1"}})

	loadErr := s.configSet.Load(s.valueContainer)
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

func (s *CommonDecodersTestSuite) Test_Decode_IP_FromNil() {
	s.valueContainer.On("Errors").Return([]error{})
	s.valueContainer.On("Get").Return([]any{map[string]any{"test_ip": nil}})

	loadErr := s.configSet.Load(s.valueContainer)
	s.Require().NoError(loadErr)

	type targetStruct struct {
		TestIP net.IP `cfg:"test_ip"`
	}

	var target targetStruct

	decodeErr := s.configSet.Decode(&target)

	s.Error(decodeErr)
}

func (s *CommonDecodersTestSuite) Test_Decode_IP_FromInvalidFormat() {
	s.valueContainer.On("Errors").Return([]error{})
	s.valueContainer.On("Get").Return([]any{map[string]any{"test_ip_invalid_format": "127.0.0.0.1"}})

	loadErr := s.configSet.Load(s.valueContainer)
	s.Require().NoError(loadErr)

	type targetStruct struct {
		TestIP net.IP `cfg:"test_ip_invalid_format"`
	}

	var target targetStruct

	decodeErr := s.configSet.Decode(&target)

	s.Error(decodeErr)
}

func (s *CommonDecodersTestSuite) Test_Decode_JSONRawMessage() {
	s.valueContainer.On("Errors").Return([]error{})
	s.valueContainer.On("Get").Return([]any{map[string]any{"test_rawMessage": map[string]any{"rawMessage": "It's raw"}}})

	loadErr := s.configSet.Load(s.valueContainer)
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

func (s *CommonDecodersTestSuite) Test_Decode_JSONRawMessage_FromNil() {
	s.valueContainer.On("Errors").Return([]error{})
	s.valueContainer.On("Get").Return([]any{map[string]any{"test_rawMessage": nil}})

	loadErr := s.configSet.Load(s.valueContainer)
	s.Require().NoError(loadErr)

	type targetStruct struct {
		TestRawMessage json.RawMessage `cfg:"test_rawMessage"`
	}

	var target targetStruct

	decodeErr := s.configSet.Decode(&target)

	s.Error(decodeErr)
}

func (s *CommonDecodersTestSuite) Test_Decode_JSONRawMessage_WithUnmarshalable() {
	s.valueContainer.On("Errors").Return([]error{})
	s.valueContainer.On("Get").Return([]any{map[string]any{"test_rawMessage": map[string]any{"rawMessage": "It's raw"}}})

	loadErr := s.configSet.Load(s.valueContainer)
	s.Require().NoError(loadErr)

	type targetStruct struct {
		TestRawMessage json.RawMessage `cfg:"test_rawMessage"`
	}

	var target targetStruct

	s.configSet.OverrideValue(map[string]any{"test_rawMessage": map[bool]string{true: "It's raw"}})

	decodeErr := s.configSet.Decode(&target)

	s.Error(decodeErr)
}

func (s *CommonDecodersTestSuite) Test_Decode_Time() {
	s.valueContainer.On("Errors").Return([]error{})
	s.valueContainer.On("Get").Return([]any{map[string]any{"test_time": "2025-01-13T16:00:00+09:00"}})

	loadErr := s.configSet.Load(s.valueContainer)
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

func (s *CommonDecodersTestSuite) Test_Decode_Time_FromNil() {
	s.valueContainer.On("Errors").Return([]error{})
	s.valueContainer.On("Get").Return([]any{map[string]any{"test_time": nil}})

	loadErr := s.configSet.Load(s.valueContainer)
	s.Require().NoError(loadErr)

	type targetStruct struct {
		TestTime time.Time `cfg:"test_time"`
	}

	var target targetStruct

	decodeErr := s.configSet.Decode(&target)

	s.Error(decodeErr)
}

func (s *CommonDecodersTestSuite) Test_Decode_Time_FromTime() {
	timeValue := time.Date(2025, 1, 13, 16, 0, 0, 0, time.FixedZone("", 9*60*60))

	s.valueContainer.On("Errors").Return([]error{})
	s.valueContainer.On("Get").Return([]any{map[string]any{"test_time": timeValue}})

	loadErr := s.configSet.Load(s.valueContainer)
	s.Require().NoError(loadErr)

	type targetStruct struct {
		TestTime time.Time `cfg:"test_time"`
	}

	var (
		target   targetStruct
		expected = timeValue
	)

	decodeErr := s.configSet.Decode(&target)

	s.Equal(expected, target.TestTime)
	s.NoError(decodeErr)
}

func (s *CommonDecodersTestSuite) Test_Decode_Time_FromInvalidType() {
	s.valueContainer.On("Errors").Return([]error{})
	s.valueContainer.On("Get").Return([]any{map[string]any{"test_time_invalid_type": false}})

	loadErr := s.configSet.Load(s.valueContainer)
	s.Require().NoError(loadErr)

	type targetStruct struct {
		TestTime time.Time `cfg:"test_time_invalid_type"`
	}

	var target targetStruct

	decodeErr := s.configSet.Decode(&target)

	s.Error(decodeErr)
}

func (s *CommonDecodersTestSuite) Test_Decode_Time_FromInvalidFormat() {
	s.valueContainer.On("Errors").Return([]error{})
	s.valueContainer.On("Get").Return([]any{map[string]any{"test_time_invalid_format": "2025. 01. 13. 16:00:00 UTC"}})

	loadErr := s.configSet.Load(s.valueContainer)
	s.Require().NoError(loadErr)

	type targetStruct struct {
		TestTime time.Time `cfg:"test_time_invalid_format"`
	}

	var target targetStruct

	decodeErr := s.configSet.Decode(&target)

	s.Error(decodeErr)
}

func (s *CommonDecodersTestSuite) Test_Decode_URL() {
	s.valueContainer.On("Errors").Return([]error{})
	s.valueContainer.On("Get").Return([]any{map[string]any{"test_url": "http://www.test.com"}})

	loadErr := s.configSet.Load(s.valueContainer)
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

func (s *CommonDecodersTestSuite) Test_Decode_URL_FromNil() {
	s.valueContainer.On("Errors").Return([]error{})
	s.valueContainer.On("Get").Return([]any{map[string]any{"test_url": nil}})

	loadErr := s.configSet.Load(s.valueContainer)
	s.Require().NoError(loadErr)

	type targetStruct struct {
		TestURL *url.URL `cfg:"test_url"`
	}

	var target targetStruct

	decodeErr := s.configSet.Decode(&target)

	s.Error(decodeErr)
}

func (s *CommonDecodersTestSuite) Test_Decode_URL_FromInvalidFormat() {
	s.valueContainer.On("Errors").Return([]error{})
	s.valueContainer.On("Get").Return([]any{map[string]any{"test_url_invalid_format": "missing_protocol://test.com"}})

	loadErr := s.configSet.Load(s.valueContainer)
	s.Require().NoError(loadErr)

	type targetStruct struct {
		TestURL *url.URL `cfg:"test_url_invalid_format"`
	}

	var target targetStruct

	decodeErr := s.configSet.Decode(&target)

	s.Error(decodeErr)
}
