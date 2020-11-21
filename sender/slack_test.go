package sender

import (
	"testing"
)

func TestSendStatus(t *testing.T) {
	s := Slack{Webhook: "https://slack.com/api/api.test"}
	err := s.SendStatus("elastic-sls", "foo")
	if err != nil {
		t.Error(err)
	}
}

func TestSendAvailability(t *testing.T) {
	s := Slack{Webhook: "https://slack.com/api/api.test"}
	err := s.SendAvailability("elastic-sls", true)
	err = s.SendAvailability("elastic-sls", false)
	if err != nil {
		t.Error(err)
	}
}
