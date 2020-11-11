package sender

import (
	"fmt"
)

// Slack should be initialized with valid webhook for posting messages
type Mock struct {
	C chan (string)
}

// SendStatus to Slack webhook
func (m *Mock) SendStatus(service, status string) error {
	m.C <- fmt.Sprintf("%s %s", service, status)
	return nil
}

// SendAvailability to Slack webhook
func (m *Mock) SendAvailability(service string, available bool) error {
	m.C <- fmt.Sprintf("%s %v", service, available)
	return nil
}
