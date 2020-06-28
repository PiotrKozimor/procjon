package procjonagent

import (
	"log"

	"github.com/coreos/go-systemd/v22/dbus"
)

var systemdUnitStatuses = map[int32]string{
	0: "active",
	1: "reloading",
	2: "inactive",
	3: "failed",
	4: "activating",
	5: "deactivating",
	6: "unknown",
}

type SystemdServiceMonitor struct {
	statuses   map[int32]string
	unitName   string
	connection *dbus.Conn
}

type systemdUnitStatus struct {
	ActiveStatus string
}

func (m *SystemdServiceMonitor) GetCurrentStatus() int32 {
	statuses, err := m.connection.ListUnitsByNames([]string{m.unitName})
	if err != nil {
		return 6
	}
	for code, status := range m.statuses {
		if status == statuses[0].ActiveState {
			return code
		}
	}
	log.Printf("Could not find received status in statuses!")
	return 6
}

func (m *SystemdServiceMonitor) GetStatuses() map[int32]string {
	return m.statuses
}
