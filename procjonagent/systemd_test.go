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
		statuses:   systemdUnitStatuses,
		connection: conn,
		unitName:   "dbus.service",
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
		statuses: map[int32]string{
			0: "active",
			1: "reloading",
			2: "inactive",
			3: "failed",
			4: "activating",
			5: "deactivating",
			6: "unknown",
		},
		connection: conn,
		unitName:   "dbusss.service",
	}
	status := dut.GetCurrentStatus()
	if status != 2 {
		t.Errorf("Got: %d, wanted: %d", status, 2)
	}
}
