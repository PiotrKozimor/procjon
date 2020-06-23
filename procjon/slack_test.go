package procjon

import (
	"os"
	"testing"
)

// func TestSendMessage()
func TestSendStatus(t *testing.T) {
	s := Slack{Webhook: os.Getenv("PROCJON_SLACK_WEBHOOK")}
	err := s.SendStatus("elastic-sls", "foo")
	if err != nil {
		t.Error(err)
	}
}

func TestSendAvailability(t *testing.T) {
	s := Slack{Webhook: os.Getenv("PROCJON_SLACK_WEBHOOK")}
	err := s.SendAvailability("elastic-sls", true)
	err = s.SendAvailability("elastic-sls", false)
	if err != nil {
		t.Error(err)
	}
}
