package procjon

import (
	"os"
	"testing"
)

// func TestSendMessage()
func TestSendStatus(t *testing.T) {
	if os.Getenv("TRAVIS") == "true" {
		t.Skip("TRAVIS is set to true, skipping.")
	}
	s := Slack{Webhook: os.Getenv("PROCJON_SLACK_WEBHOOK")}
	err := s.SendStatus("elastic-sls", "foo")
	if err != nil {
		t.Error(err)
	}
}

func TestSendAvailability(t *testing.T) {
	if os.Getenv("TRAVIS") == "true" {
		t.Skip("TRAVIS is set to true, skipping.")
	}
	s := Slack{Webhook: os.Getenv("PROCJON_SLACK_WEBHOOK")}
	err := s.SendAvailability("elastic-sls", true)
	err = s.SendAvailability("elastic-sls", false)
	if err != nil {
		t.Error(err)
	}
}

func TestSendAvailabilities(t *testing.T) {
	if os.Getenv("TRAVIS") == "true" {
		t.Skip("TRAVIS is set to true, skipping.")
	}
	s := Slack{Webhook: os.Getenv("PROCJON_SLACK_WEBHOOK")}
	availabilities := make(chan bool)
	go SendAvailabilities(&s, "elastic-sls", availabilities)
	availabilities <- true
	availabilities <- false
}

func TestSendStatuses(t *testing.T) {
	if os.Getenv("TRAVIS") == "true" {
		t.Skip("TRAVIS is set to true, skipping.")
	}
	s := Slack{Webhook: os.Getenv("PROCJON_SLACK_WEBHOOK")}
	statuses := make(chan string)
	go SendStatuses(&s, "elastic-sls", statuses)
	statuses <- "foo"
	statuses <- "bar"
}
