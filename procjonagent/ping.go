package procjonagent

import (
	"github.com/sparrc/go-ping"
)

var pingStatuses = map[int32]string{
	0: "pinged",
	1: "unreachable",
}

type PingMonitor struct {
	Pinger ping.Pinger
}

func (p *PingMonitor) GetStatuses() map[int32]string {
	return pingStatuses
}

func (p *PingMonitor) GetCurrentStatus() int32 {
	p.Pinger.Run()
	stats := p.Pinger.Statistics()
	if stats.PacketsRecv == 0 {
		return 1
	} else {
		return 0
	}
}
