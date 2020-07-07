package procjonagent

import (
	"reflect"
	"testing"

	"github.com/coreos/go-systemd/v22/dbus"
)

func TestSystemdGetCurrentStatus(t *testing.T) {
	conn, err := dbus.New()
	if err != nil {
		t.Error(err)
	}
	dut := SystemdServiceMonitor{
		Connection: conn,
		UnitName:   "systemd-journald.service",
	}
	status := dut.GetCurrentStatus()
	if status != 0 {
		t.Errorf("Got: %d, wanted: %d", status, 0)
	}
}

func TestSystemdGetCurrentStatus_InvalidName(t *testing.T) {
	conn, err := dbus.New()
	if err != nil {
		t.Error(err)
	}
	dut := SystemdServiceMonitor{
		Connection: conn,
		UnitName:   "dbusss.service",
	}
	status := dut.GetCurrentStatus()
	if status != 2 {
		t.Errorf("Got: %d, wanted: %d", status, 2)
	}
}

func TestSystemdGetStatuses(t *testing.T) {
	e := SystemdServiceMonitor{}
	statuses := e.GetStatuses()
	if !reflect.DeepEqual(statuses, systemdUnitStatuses) {
		t.Errorf("Got: %+v, wanted: %+v", statuses, systemdUnitStatuses)
	}

}
