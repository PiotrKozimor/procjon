package procjon

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

type slack struct {
	webhook string
}

// SendStatus to Slack webhook
func (s *slack) SendStatus(service, status string) {
	err := sendSlackMessage(s.webhook, strings.NewReader(fmt.Sprintf("Service %s change it's status to: %s.", service, status)))
	if err != nil {
		log.Printf("Could not send status update to Slack: %v", err)
	}
}

// SendAvailability to Slack webhook
func (s *slack) SendAvailability(service string, available bool) {
	var err error
	if available {
		err = sendSlackMessage(s.webhook, strings.NewReader(fmt.Sprintf("Service %s is available.")))
	} else {
		err = sendSlackMessage(s.webhook, strings.NewReader(fmt.Sprintf("Service %s is not available.")))
	}
	if err != nil {
		log.Printf("Could not send availability update to Slack: %v", err)
	}
}

// SendSlackMessage for given webhook
func sendSlackMessage(webhook string, r io.Reader) error {
	_, err := http.Post(webhook, "application/json", r)
	return err
}
