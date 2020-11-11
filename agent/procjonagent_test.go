package agent

import (
	"os"
	"testing"
	"time"

	"github.com/PiotrKozimor/procjon/procjon"
)

type MockMonitor struct {
}

func (m *MockMonitor) GetCurrentStatus() int32 {
	return 0
}
func (m *MockMonitor) GetStatuses() map[int32]string {
	statuses := map[int32]string{0: "ok", 1: "nok"}
	return statuses
}

func TestHandleMonitor(t *testing.T) {
	if os.Getenv("SKIP_HANDLE_MONITOR") == "true" {
		t.Skip("Skipping TestHandleMonitor- conflict for listening on localhost.")
	}
	go func() {
		procjon.RootCmd.Execute()
	}()
	time.Sleep(time.Second * 1)
	m := MockMonitor{}
	go func() {
		err := HandleMonitor(&m)
		t.Error(err)
	}()
	time.Sleep(time.Second * 10)
}

func TestNewConnection(t *testing.T) {
	dut := New
}
