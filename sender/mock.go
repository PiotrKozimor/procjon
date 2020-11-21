package sender

import (
	"fmt"
	"testing"
)

// Mock sender is used for testing procjon server.
type Mock struct {
	Status       chan (string)
	Availability chan (string)
	T            *testing.T
}

// SendStatus sends status to Status channel formatted as "service status", eg. "redis ok".
// It also logs service and status to testing.T.
func (m *Mock) SendStatus(service, status string) error {
	m.T.Logf("Service: %s, status: %s", service, status)
	m.Status <- fmt.Sprintf("%s %s", service, status)
	return nil
}

// SendAvailability sends availability to Availability channel formatted as "service available", eg. "redis false".
// It also logs service and availability to testing.T.
func (m *Mock) SendAvailability(service string, available bool) error {
	m.T.Logf("Service: %s, availability: %t", service, available)
	m.Availability <- fmt.Sprintf("%s %v", service, available)
	return nil
}
