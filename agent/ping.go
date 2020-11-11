package agent

var pingStatuses = map[int32]string{
	0: "ok",
}

type PingMonitor struct {
}

func (p *PingMonitor) GetStatuses() map[int32]string {
	return pingStatuses
}

func (p *PingMonitor) GetCurrentStatus() int32 {
	return 0
}
