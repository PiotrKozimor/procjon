package procjonagent

import (
	"testing"

	"github.com/coreos/go-systemd/v22/dbus"
)

func TestGetCurrentStatus(t *testing.T) {
	conn, err := dbus.New()
	if err != nil {
		t.Error(err)
	}
	dut := SystemdServiceMonitor{
		Connection: conn,
		UnitName:   "dbus.service",
	}
	status := dut.GetCurrentStatus()
	if status != 0 {
		t.Errorf("Got: %d, wanted: %d", status, 0)
	}
}

func TestGetCurrentStatus_InvalidName(t *testing.T) {
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
