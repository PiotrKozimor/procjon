package agent

var pingStatuses = []string{
	"ok",
}

// Ping is simplest agenter. It can be used to monitor host. When host loses connection with procjon, procjon will send availability update to Slack.
type Ping struct {
}

// GetStatuses return only "ok" status.
func (p *Ping) GetStatuses() []string {
	return pingStatuses
}

// GetCurrentStatus returns always 0.
func (p *Ping) GetCurrentStatus() uint32 {
	return 0
}
