package agent

var pingStatuses = []string{
	"ok",
}

type Ping struct {
}

func (p *Ping) GetStatuses() []string {
	return pingStatuses
}

func (p *Ping) GetCurrentStatus() uint32 {
	return 0
}
