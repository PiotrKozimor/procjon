package procjon

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Slack should be initialized with valid webhook for posting messages
type Slack struct {
	Webhook string
}

// SlackMessage contains text to send, is meant to be serialized to json
type SlackMessage struct {
	Text string `json:"text"`
}

// SendStatus to Slack webhook
func (s *Slack) SendStatus(service, status string) error {
	err := sendSlackMessage(s.Webhook, SlackMessage{Text: fmt.Sprintf("Service %s change it's status to: %s", service, status)})
	if err != nil {
		log.Printf("Could not send status update to Slack: %v", err)
	}
	return err
}

// SendAvailability to Slack webhook
func (s *Slack) SendAvailability(service string, available bool) error {
	var err error
	if available {
		err = sendSlackMessage(s.Webhook, SlackMessage{Text: fmt.Sprintf("Service %s is available.", service)})
	} else {
		err = sendSlackMessage(s.Webhook, SlackMessage{Text: fmt.Sprintf("Service %s is not available.", service)})
	}
	if err != nil {
		log.Printf("Could not send availability update to Slack: %v", err)
	}
	return err
}

// SendSlackMessage for given webhook
func sendSlackMessage(webhook string, message SlackMessage) error {
	marshalled, err := json.Marshal(message)
	r := bytes.NewReader(marshalled)
	resp, err := http.Post(webhook, "application/json", r)
	if resp.StatusCode != 200 {
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("cannot send message to Slack: %d, %s", resp.StatusCode, data)
	}
	return err
}
