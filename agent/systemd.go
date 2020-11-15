package agent

import (
	"github.com/coreos/go-systemd/v22/dbus"
	log "github.com/sirupsen/logrus"
)

var systemdUnitStatuses = []string{
	"active",
	"reloading",
	"inactive",
	"failed",
	"activating",
	"deactivating",
	"unknown",
}

// SystemdServiceMonitor hold unit name to monitor and Connection
// to talk to dbus.
type SystemdServiceMonitor struct {
	UnitName   string
	Connection listUnits
}

type listUnits interface {
	ListUnitsByNames(units []string) ([]dbus.UnitStatus, error)
}

type systemdUnitStatus struct {
	ActiveStatus string
}

// GetCurrentStatus of SystemdServiceMonitor.Unit from dbus.
func (m *SystemdServiceMonitor) GetCurrentStatus() uint32 {
	statuses, err := m.Connection.ListUnitsByNames([]string{m.UnitName})
	if err != nil {
		log.Error(err)
		return 6
	}
	for code, status := range systemdUnitStatuses {
		if status == statuses[0].ActiveState {
			return uint32(code)
		}
	}
	log.Errorf("Could not find received status in statuses!")
	return 6
}

// GetStatuses statuses defined for SystemdServiceMonitor.
func (m *SystemdServiceMonitor) GetStatuses() []string {
	return systemdUnitStatuses
}
