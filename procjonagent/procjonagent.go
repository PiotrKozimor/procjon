package procjonagent

type StatusMonitor interface {
	GetCurrentStatus() int32
	GetStatuses() map[int32]string
}
