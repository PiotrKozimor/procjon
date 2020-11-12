package sender

import (
	"fmt"
	"testing"
)

// Slack should be initialized with valid webhook for posting messages
type Mock struct {
	Status       chan (string)
	Availability chan (string)
	T            *testing.T
}

// SendStatus to Slack webhook
func (m *Mock) SendStatus(service, status string) error {
	m.T.Logf("Service: %s, status: %s", service, status)
	m.Status <- fmt.Sprintf("%s %s", service, status)
	return nil
}

// SendAvailability to Slack webhook
func (m *Mock) SendAvailability(service string, available bool) error {
	m.T.Logf("Service: %s, availability: %t", service, available)
	m.Availability <- fmt.Sprintf("%s %v", service, available)
	return nil
}
