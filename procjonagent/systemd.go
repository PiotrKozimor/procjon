package procjonagent

import (
	"github.com/coreos/go-systemd/v22/dbus"
	log "github.com/sirupsen/logrus"
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
	UnitName   string
	Connection *dbus.Conn
}

type systemdUnitStatus struct {
	ActiveStatus string
}

func (m *SystemdServiceMonitor) GetCurrentStatus() int32 {
	statuses, err := m.Connection.ListUnitsByNames([]string{m.UnitName})
	if err != nil {
		return 6
	}
	for code, status := range systemdUnitStatuses {
		if status == statuses[0].ActiveState {
			return code
		}
	}
	log.Errorf("Could not find received status in statuses!")
	return 6
}

func (m *SystemdServiceMonitor) GetStatuses() map[int32]string {
	return systemdUnitStatuses
}
