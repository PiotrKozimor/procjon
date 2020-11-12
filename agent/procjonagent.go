package agent

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
	"time"

	"github.com/PiotrKozimor/procjon/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Agenter is implemented by all of procjonagents.
type Agenter interface {
	GetCurrentStatus() int32
	GetStatuses() map[int32]string
}

type ConnectionOpts struct {
	Endpoint         string
	RootCertPath     string
	AgentCertPath    string
	AgentKeyCertPath string
}

type Agent struct {
	Conn         *grpc.ClientConn
	Indentifier  string
	TimeoutSec   int32
	UpdatePeriod time.Duration
}

var DefaultOpts = ConnectionOpts{
	Endpoint:         "localhost:8080",
	AgentCertPath:    "procjonagent.pem",
	AgentKeyCertPath: "procjonagent.key",
	RootCertPath:     "ca.pem",
}

// NewConnection initializes connection to given endpoint.
// Certificates in ConnectionOpts must be provided.
// Connection must be closed. This is done in (*Agent) Run function.
func NewConnection(opts *ConnectionOpts) (*grpc.ClientConn, error) {
	b, err := ioutil.ReadFile(opts.RootCertPath)
	if err != nil {
		return nil, err
	}
	cp := x509.NewCertPool()
	if !cp.AppendCertsFromPEM(b) {
		return nil, errors.New("credentials: failed to append certificates")
	}
	cert, err := tls.LoadX509KeyPair(opts.AgentCertPath, opts.AgentKeyCertPath)
	if err != nil {
		return nil, err
	}
	config := &tls.Config{
		InsecureSkipVerify: false,
		RootCAs:            cp,
		Certificates:       []tls.Certificate{cert},
	}
	conn, err := grpc.Dial(opts.Endpoint, grpc.WithTransportCredentials(credentials.NewTLS(config)))
	return conn, err
}

// Run registers service in running procjon server and periodically send
// statusCode to procjon server. Provide as argument any agent that implements Agenter interface.
func (a *Agent) Run(ar Agenter) error {
	service := pb.Service{
		ServiceIdentifier: a.Indentifier,
		Timeout:           a.TimeoutSec,
		Statuses:          ar.GetStatuses(),
	}
	serviceStatus := pb.ServiceStatus{
		ServiceIdentifier: service.ServiceIdentifier,
		StatusCode:        0,
	}
	defer a.Conn.Close()
	cl := pb.NewProcjonClient(a.Conn)
	_, err := cl.RegisterService(context.Background(), &service)
	if err != nil {
		return err
	}
	stream, err := cl.SendServiceStatus(context.Background())
	if err != nil {
		return err
	}
	for {
		status := ar.GetCurrentStatus()
		serviceStatus.StatusCode = status
		err = stream.Send(&serviceStatus)
		if err != nil {
			return err
		}
		time.Sleep(a.UpdatePeriod)
	}
}
