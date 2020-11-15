package agent

var pingStatuses = []string{
	"ok",
}

type PingMonitor struct {
}

func (p *PingMonitor) GetStatuses() []string {
	return pingStatuses
}

func (p *PingMonitor) GetCurrentStatus() uint32 {
	return 0
}
