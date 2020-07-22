package procjonagent

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/PiotrKozimor/procjon/pb"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// ServiceMonitor can be used to define custom monitor. It is used in
// handleMonitor function.
type ServiceMonitor interface {
	GetCurrentStatus() int32
	GetStatuses() map[int32]string
}

// RootCmd is default command that can be consumed by procjonagent.
// Common flags are defined for this command
var RootCmd = &cobra.Command{}

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.Stamp})
	log.SetOutput(os.Stderr)
	RootCmd.Version = "v0.2.0-alpha"
	RootCmd.Flags().StringVarP(&endpoint, "endpoint", "e", "localhost:8080", "gRPC endpoint of procjon server")
	RootCmd.Flags().StringVarP(&identifier, "service", "s", "foo", "service identifier")
	RootCmd.Flags().Int32VarP(&timeout, "timeout", "t", 10, "procjon service timeout [s]")
	RootCmd.Flags().Int32VarP(&period, "period", "p", 4, "period for agent to sent status updates with [s]")
	RootCmd.Flags().StringVarP(&LogLevel, "loglevel", "l", "warning", "logrus log level")
	RootCmd.Flags().StringVar(&rootCertPath, "root-cert", "ca.pem", "root certificate path")
	RootCmd.Flags().StringVarP(&serverCertPath, "cert", "c", "procjon.pem", "certificate path")
	RootCmd.Flags().StringVarP(&serverKeyCertPath, "key-cert", "k", "procjon.key", "key certificate path")
}

var (
	endpoint   string
	identifier string
	timeout    int32
	// LogLevel according to logrus level naming convention.
	LogLevel          string
	period            int32
	rootCertPath      string
	serverKeyCertPath string
	serverCertPath    string
)

// HandleMonitor registers service and periodically send
// statusCode to procjon.
func HandleMonitor(m ServiceMonitor) error {
	service := pb.Service{
		ServiceIdentifier: identifier,
		Timeout:           timeout,
		Statuses:          m.GetStatuses(),
	}
	serviceStatus := pb.ServiceStatus{
		ServiceIdentifier: service.ServiceIdentifier,
		StatusCode:        0,
	}
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dir)
	b, err := ioutil.ReadFile(rootCertPath)
	if err != nil {
		return err
	}
	cp := x509.NewCertPool()
	if !cp.AppendCertsFromPEM(b) {
		return errors.New("credentials: failed to append certificates")
	}
	cert, err := tls.LoadX509KeyPair(serverCertPath, serverKeyCertPath)
	if err != nil {
		return err
	}
	config := &tls.Config{
		InsecureSkipVerify: false,
		RootCAs:            cp,
		Certificates:       []tls.Certificate{cert},
	}
	conn, err := grpc.Dial(endpoint, grpc.WithTransportCredentials(credentials.NewTLS(config)))
	if err != nil {
		return err
	}
	defer conn.Close()
	cl := pb.NewProcjonClient(conn)
	_, err = cl.RegisterService(context.Background(), &service)
	if err != nil {
		return err
	}
	stream, err := cl.SendServiceStatus(context.Background())
	if err != nil {
		return err
	}
	for {
		status := m.GetCurrentStatus()
		serviceStatus.StatusCode = status
		err = stream.Send(&serviceStatus)
		if err != nil {
			return err
		}
		time.Sleep(5 * time.Second)
	}
}
