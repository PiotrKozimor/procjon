package procjon

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// AvailabilitySender defines sending availability update for service.
type AvailabilitySender interface {
	SendAvailability(service string, availability bool) error
}

// StatusSender defines sending status update for service.
type StatusSender interface {
	SendStatus(service string, status string) error
}

// AvailabilityStatusSender defines sending status and availability
// update for service.
type AvailabilityStatusSender interface {
	AvailabilitySender
	StatusSender
}

// Slack should be initialized with valid webhook for posting messages
type Slack struct {
	Webhook string
}

// SlackMessage contains text to send, is meant to be serialized to json
type SlackMessage struct {
	Text string `json:"text"`
}

// SendStatuses forever from status channel
func SendStatuses(sender StatusSender, service string, status chan string) {
	for {
		statusToSend := <-status
		log.WithField("service", service).Infof("Sending status %s", statusToSend)
		if err := sender.SendStatus(service, statusToSend); err != nil {
			log.Error(err)
		}
	}
}

// SendAvailabilities forever from available channel
func SendAvailabilities(sender AvailabilitySender, service string, available chan bool) {
	for {
		availability := <-available
		log.WithField("service", service).Infof("Sending availability %t", availability)
		if err := sender.SendAvailability(service, availability); err != nil {
			log.Error(err)
		}
	}
}

// SendStatus to Slack webhook
func (s *Slack) SendStatus(service, status string) error {
	err := sendSlackMessage(s.Webhook, SlackMessage{Text: fmt.Sprintf("Service %s changed it's status to: %s", service, status)})
	if err != nil {
		log.Errorf("cannot not send status update to Slack: %v", err)
	}
	return err
}

// SendAvailability to Slack webhook
func (s *Slack) SendAvailability(service string, available bool) error {
	var err error
	if available {
		err = sendSlackMessage(s.Webhook, SlackMessage{Text: fmt.Sprintf("Agent %s is available.", service)})
	} else {
		err = sendSlackMessage(s.Webhook, SlackMessage{Text: fmt.Sprintf("Agent %s is not available.", service)})
	}
	if err != nil {
		log.Errorf("cannot not send availability update to Slack: %v", err)
	}
	return err
}

// SendSlackMessage for given webhook
func sendSlackMessage(webhook string, message SlackMessage) error {
	marshalled, err := json.Marshal(message)
	r := bytes.NewReader(marshalled)
	resp, err := http.Post(webhook, "application/json", r)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("cannot send message to Slack: %d, %s", resp.StatusCode, data)
	}
	return err
}
