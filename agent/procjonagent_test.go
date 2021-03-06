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
	dut := Service{
		Indentifier:     "redis",
		TimeoutSec:      2,
		UpdatePeriodSec: 1,
	}
	go dut.Run(&Ping{}, conn)
	availability := <-mock.Availability
	assert.Equal(t, availability, "redis true")
	time.Sleep(time.Second * 3)
}

func TestNewConnection(t *testing.T) {
	conn, err := NewConnection(&ConnectionOpts{
		CertPath:     "../.certs/procjonagent.pem",
		KeyCertPath:  "../.certs/procjonagent.key",
		RootCertPath: "../.certs/ca.pem",
		Endpoint:     "localhost:8080",
	})
	assert.Nil(t, err)
	conn.Close()
}

func TestNewConnectionNoCACert(t *testing.T) {
	_, err := NewConnection(&ConnectionOpts{
		CertPath:     "../.certs/procjonagent.pem",
		KeyCertPath:  "../.certs/procjonagent.key",
		RootCertPath: "ca.pem",
		Endpoint:     "localhost:8080",
	})
	assert.EqualError(t, err, "open ca.pem: no such file or directory")
}

func TestNewConnectionBadCACert(t *testing.T) {
	_, err := NewConnection(&ConnectionOpts{
		CertPath:     "../.certs/procjonagent.pem",
		KeyCertPath:  "../.certs/procjonagent.key",
		RootCertPath: "bad.pem",
		Endpoint:     "localhost:8080",
	})
	assert.EqualError(t, err, "credentials: failed to append certificates")
}

func TestNewConnectionBadAgentCert(t *testing.T) {
	_, err := NewConnection(&ConnectionOpts{
		CertPath:     "procjonagent.pem",
		KeyCertPath:  "../.certs/procjonagent.key",
		RootCertPath: "../.certs/ca.pem",
		Endpoint:     "localhost:8080",
	})
	assert.EqualError(t, err, "open procjonagent.pem: no such file or directory")
}
