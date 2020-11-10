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

func TestSendAvailabilities(t *testing.T) {
	s := Slack{Webhook: os.Getenv("PROCJON_SLACK_WEBHOOK")}
	availabilities := make(chan bool)
	go SendAvailabilities(&s, "elastic-sls", availabilities)
	availabilities <- true
	availabilities <- false
}

func TestSendStatuses(t *testing.T) {
	s := Slack{Webhook: os.Getenv("PROCJON_SLACK_WEBHOOK")}
	statuses := make(chan string)
	go SendStatuses(&s, "elastic-sls", statuses)
	statuses <- "foo"
	statuses <- "bar"
}

func TestSendStatusesBadMethod(t *testing.T) {
	s := Slack{Webhook: "https://slack.com/foo"}
	statuses := make(chan string)
	go SendStatuses(&s, "elastic-sls", statuses)
	statuses <- "foo"
	statuses <- "bar"
}

func TestSendAvailabilitiesBadMethod(t *testing.T) {
	s := Slack{Webhook: "https://slack.com/foo"}
	availabilities := make(chan bool)
	go SendAvailabilities(&s, "elastic-sls", availabilities)
	availabilities <- true
	availabilities <- false
}

func TestSendStatusesBadURL(t *testing.T) {
	s := Slack{Webhook: "https://sladw.com/foo"}
	statuses := make(chan string)
	go SendStatuses(&s, "elastic-sls", statuses)
	statuses <- "foo"
	statuses <- "bar"
}

func TestSendAvailabilitiesBadURL(t *testing.T) {
	s := Slack{Webhook: "https://sladw.com/foo"}
	availabilities := make(chan bool)
	go SendAvailabilities(&s, "elastic-sls", availabilities)
	availabilities <- true
	availabilities <- false
}
