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
	endpoint         string
	rootCertPath     string
	agentCertPath    string
	agentKeyCertPath string
}

type Agent struct {
	conn         *grpc.ClientConn
	indentifier  string
	timeoutSec   int32
	updatePeriod time.Duration
}

var DefaultOpts = ConnectionOpts{
	endpoint:         "localhost:8080",
	agentCertPath:    "procjonagent.pem",
	agentKeyCertPath: "procjonagent.key",
	rootCertPath:     "ca.pem",
}

// NewConnection initializes connection to given endpoint.
// Certificates in ConnectionOpts must be provided.
// Connection must be closed. This is done in (*Agent) Run function.
func NewConnection(opts *ConnectionOpts) (*grpc.ClientConn, error) {
	b, err := ioutil.ReadFile(opts.rootCertPath)
	if err != nil {
		return nil, err
	}
	cp := x509.NewCertPool()
	if !cp.AppendCertsFromPEM(b) {
		return nil, errors.New("credentials: failed to append certificates")
	}
	cert, err := tls.LoadX509KeyPair(opts.agentCertPath, opts.agentKeyCertPath)
	if err != nil {
		return nil, err
	}
	config := &tls.Config{
		InsecureSkipVerify: false,
		RootCAs:            cp,
		Certificates:       []tls.Certificate{cert},
	}
	conn, err := grpc.Dial(opts.endpoint, grpc.WithTransportCredentials(credentials.NewTLS(config)))
	return conn, err
}

// Run registers service in running procjon server and periodically send
// statusCode to procjon server. Provide as argument any agent that implements Agenter interface.
func (a *Agent) Run(ar Agenter) error {
	service := pb.Service{
		ServiceIdentifier: a.indentifier,
		Timeout:           a.timeoutSec,
		Statuses:          ar.GetStatuses(),
	}
	serviceStatus := pb.ServiceStatus{
		ServiceIdentifier: service.ServiceIdentifier,
		StatusCode:        0,
	}
	defer a.conn.Close()
	cl := pb.NewProcjonClient(a.conn)
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
		time.Sleep(a.updatePeriod)
	}
}
