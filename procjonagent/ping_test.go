package procjonagent

import (
	"os"
	"testing"
	"time"

	"github.com/sparrc/go-ping"
)

func TestPing(t *testing.T) {
	pinger, err := ping.NewPinger(os.Getenv("PING"))
	if err != nil {
		t.Fatal(err)
	}
	pinger.Count = 3
	pinger.Timeout = time.Second * 4
	pinger.Interval = time.Second * 1
	pinger.SetPrivileged(true)
	dut := PingMonitor{Pinger: *pinger}
	status := dut.GetCurrentStatus()
	t.Log(status)
	if status != 0 {
		t.Errorf("Got %d, expected 0", status)
	}
}

func TestNoPing(t *testing.T) {
	pinger, err := ping.NewPinger(os.Getenv("NOPING"))
	if err != nil {
		t.Fatal(err)
	}
	pinger.Count = 3
	pinger.Timeout = time.Second * 7
	pinger.Interval = time.Second * 1
	pinger.SetPrivileged(true)
	dut := PingMonitor{Pinger: *pinger}
	status := dut.GetCurrentStatus()
	t.Log(status)
	// log.Println(status)
	if status != 1 {
		t.Errorf("Got %d, expected 1", status)
	}
}