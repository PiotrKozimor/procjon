package agent

import (
	"testing"
)

func TestPing(t *testing.T) {
	dut := Ping{}
	status := dut.GetCurrentStatus()
	if status != 0 {
		t.Errorf("Got %d, expected 0", status)
	}
}
