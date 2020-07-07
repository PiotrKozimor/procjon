package procjonagent

import (
	"errors"
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

type MockedConn struct {
}

func (m MockedConn) ListUnitsByNames(units []string) ([]dbus.UnitStatus, error) {
	return nil, errors.New("mocked")
}
func TestSystemdGetCurrentStatus_ListUnitsByNamesError(t *testing.T) {
	conn := MockedConn{}
	dut := SystemdServiceMonitor{
		Connection: &conn,
		UnitName:   "dbusss.service",
	}
	status := dut.GetCurrentStatus()
	if status != 6 {
		t.Errorf("Got: %d, wanted: %d", status, 6)
	}
}

type MockedBasStatusConn struct {
}

func (m MockedBasStatusConn) ListUnitsByNames(units []string) ([]dbus.UnitStatus, error) {
	return []dbus.UnitStatus{dbus.UnitStatus{ActiveState: "foo"}}, nil
}
func TestSystemdGetCurrentStatus_ListUnitsByNamesBadStatus(t *testing.T) {
	conn := MockedBasStatusConn{}
	dut := SystemdServiceMonitor{
		Connection: &conn,
		UnitName:   "dbusss.service",
	}
	status := dut.GetCurrentStatus()
	if status != 6 {
		t.Errorf("Got: %d, wanted: %d", status, 6)
	}
}

func TestSystemdGetStatuses(t *testing.T) {
	e := SystemdServiceMonitor{}
	statuses := e.GetStatuses()
	if !reflect.DeepEqual(statuses, systemdUnitStatuses) {
		t.Errorf("Got: %+v, wanted: %+v", statuses, systemdUnitStatuses)
	}
}
