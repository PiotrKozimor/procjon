package procjonagent

import (
	"testing"
	"time"

	"github.com/sparrc/go-ping"
)

func TestPing(t *testing.T) {
	pinger, err := ping.NewPinger("google.com")
	if err != nil {
		t.Fatal(err)
	}
	pinger.Count = 3
	pinger.Timeout = time.Second * 4
	pinger.Interval = time.Second * 3
	pinger.SetPrivileged(true)
	dut := PingMonitor{Pinger: *pinger}
	status := dut.GetCurrentStatus()
	if status != 0 {
		t.Errorf("Got %d, expected 0", status)
	}
}

func TestNoPing(t *testing.T) {
	pinger, err := ping.NewPinger("10.0.0.0")
	if err != nil {
		t.Fatal(err)
	}
	pinger.Count = 3
	pinger.Timeout = time.Second * 4
	pinger.Interval = time.Second * 3
	pinger.SetPrivileged(true)
	dut := PingMonitor{Pinger: *pinger}
	status := dut.GetCurrentStatus()
	if status != 1 {
		t.Errorf("Got %d, expected 1", status)
	}
}
