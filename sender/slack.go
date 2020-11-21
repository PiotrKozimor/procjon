package sender

import (
	"fmt"

	"github.com/slack-go/slack"
)

// Slack should be initialized with valid webhook for posting messages.
// See https://api.slack.com/messaging/webhooks for instructions on how to setup incoming webhooks.
type Slack struct {
	Webhook string
}

// SendStatus nicely formatted to Slack webhook.
func (s *Slack) SendStatus(service, status string) error {
	return slack.PostWebhook(s.Webhook, &slack.WebhookMessage{
		Text: fmt.Sprintf("Service %s changed it's status to: %s", service, status),
	})
}

// SendAvailability nicely formatted to Slack webhook.
func (s *Slack) SendAvailability(service string, available bool) error {
	var text string
	if available {
		text = fmt.Sprintf("Agent %s is available.", service)
	} else {
		text = fmt.Sprintf("Agent %s is not available.", service)
	}
	return slack.PostWebhook(s.Webhook, &slack.WebhookMessage{
		Text: text,
	})
}
