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

// SystemdUnit monitors any systemd unit. Please refer to https://www.freedesktop.org/software/systemd/man/org.freedesktop.systemd1.html#Properties1
// for list of possible unit statuses.
type SystemdUnit struct {
	Name       string
	Connection listUnits
}

type listUnits interface {
	ListUnitsByNames(units []string) ([]dbus.UnitStatus, error)
}

type systemdUnitStatus struct {
	ActiveStatus string
}

// GetCurrentStatus of SystemdUnit.Name from dbus.
func (m *SystemdUnit) GetCurrentStatus() uint32 {
	statuses, err := m.Connection.ListUnitsByNames([]string{m.Name})
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

// GetStatuses returns possible statuses defined for SystemdUnit.
func (m *SystemdUnit) GetStatuses() []string {
	return systemdUnitStatuses
}
