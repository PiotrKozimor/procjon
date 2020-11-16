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
	GetCurrentStatus() uint32
	GetStatuses() []string
}

type ConnectionOpts struct {
	Endpoint     string
	RootCertPath string
	CertPath     string
	KeyCertPath  string
}

type Service struct {
	Indentifier     string
	TimeoutSec      uint32
	UpdatePeriodSec uint32
}

var DefaultOpts = ConnectionOpts{
	Endpoint:     "localhost:8080",
	CertPath:     "procjonagent.pem",
	KeyCertPath:  "procjonagent.key",
	RootCertPath: "ca.pem",
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
	cert, err := tls.LoadX509KeyPair(opts.CertPath, opts.KeyCertPath)
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
func (a *Service) Run(ar Agenter, conn *grpc.ClientConn) error {
	service := pb.Service{
		Identifier: a.Indentifier,
		Timeout:    a.TimeoutSec,
		Statuses:   ar.GetStatuses(),
	}
	serviceStatus := pb.ServiceStatus{
		Identifier: service.Identifier,
		StatusCode: 0,
	}
	defer conn.Close()
	cl := pb.NewProcjonClient(conn)
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
		time.Sleep(time.Second * time.Duration(a.UpdatePeriodSec))
	}
}
