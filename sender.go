package procjon

// Sender is used to send availability and status (e.g. to Slack) when change is detected.
type Sender interface {
	SendAvailability(service string, availability bool) error
	SendStatus(service string, status string) error
}
