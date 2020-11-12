package agent

import (
	"testing"
	"time"

	"github.com/PiotrKozimor/procjon"
	"github.com/PiotrKozimor/procjon/sender"
	"github.com/stretchr/testify/assert"
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

func TestRun(t *testing.T) {
	mock := &sender.Mock{
		T:            t,
		Availability: make(chan string),
		Status:       make(chan string)}
	conn := procjon.MustConnectOnBuffer(mock)
	dut := Agent{
		Conn:         conn,
		Indentifier:  "redis",
		TimeoutSec:   2,
		UpdatePeriod: 1,
	}
	go dut.Run(&PingMonitor{})
	availability := <-mock.Availability
	assert.Equal(t, availability, "redis true")
	time.Sleep(time.Second * 3)
}

func TestNewConnection(t *testing.T) {
	conn, err := NewConnection(&DefaultOpts)
	if err != nil {
		t.Fatal(err)
	}
	conn.Close()
}
