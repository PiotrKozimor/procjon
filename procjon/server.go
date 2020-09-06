package procjon

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net"
	"os"
	"time"

	"github.com/PiotrKozimor/procjon/pb"
	"github.com/dgraph-io/badger/v2"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

func init() {
	logger.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.Stamp})
	logger.SetOutput(os.Stderr)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.Stamp})
	log.SetOutput(os.Stderr)
	RootCmd.Flags().StringVarP(&listenURL, "listen-url", "l", "localhost:8080", "gRPC URL address to listen")
	RootCmd.Flags().StringVar(&logLevel, "loglevel", "warning", "logrus log level")
	RootCmd.Flags().StringVar(&rootCertPath, "root-cert", "ca.pem", "root certificate path")
	RootCmd.Flags().StringVarP(&serverCertPath, "cert", "c", "procjon.pem", "certificate path")
	RootCmd.Flags().StringVarP(&serverKeyCertPath, "key-cert", "k", "procjon.key", "key certificate path")
}

type Server struct {
	pb.UnimplementedProcjonServer
	Slack AvailabilityStatusSender
	DB    *badger.DB
}

var (
	logger            = log.New()
	listenURL         string
	logLevel          string
	serverCertPath    string
	serverKeyCertPath string
	rootCertPath      string
)

var RootCmd = &cobra.Command{
	Version: "v0.3.1-alpha",
	Use:     "procjon",
	Short:   "procjon monitoring server",
	Long: `Procjon is simple monitoring tool that will report change in 
availability or status of registered services. Please refer to
https://github.com/PiotrKozimor/procjon for details.`,
	Run: func(cmd *cobra.Command, args []string) {
		l, err := log.ParseLevel(logLevel)
		if err != nil {
			log.Fatalln(err)
		}
		log.SetLevel(l)
		logger.SetLevel(l)
		db, err := badger.Open(badger.DefaultOptions("services").WithLogger(logger))
		if err != nil {
			log.Fatal(err)
		}
		var s = Server{
			Slack: &Slack{Webhook: os.Getenv("PROCJON_SLACK_WEBHOOK")},
			DB:    db,
		}
		b, err := ioutil.ReadFile(rootCertPath)
		if err != nil {
			log.Fatalln(err)
		}
		cp := x509.NewCertPool()
		if !cp.AppendCertsFromPEM(b) {
			log.Fatalln("credentials: failed to append certificates")
		}
		cert, err := tls.LoadX509KeyPair(serverCertPath, serverKeyCertPath)
		if err != nil {
			log.Fatalln(err)
		}
		config := tls.Config{
			Certificates: []tls.Certificate{cert},
			ClientAuth:   tls.RequireAndVerifyClientCert,
			ClientCAs:    cp,
		}
		creds := credentials.NewTLS(&config)
		grpcServer := grpc.NewServer(grpc.Creds(creds))
		pb.RegisterProcjonServer(grpcServer, &s)
		lis, err := net.Listen("tcp4", listenURL)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		defer lis.Close()
		err = grpcServer.Serve(lis)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
	},
}

func (s *Server) SendServiceStatus(stream pb.Procjon_SendServiceStatusServer) error {
	var service pb.Service
	serviceStatus, err := stream.Recv()
	if err != nil {
		return status.Error(codes.Aborted, err.Error())
	}
	log.WithField("service", service.ServiceIdentifier).Debugf("Received statusCode %d", serviceStatus.StatusCode)
	err = LoadService(s.DB, &service, serviceStatus)
	if err != nil {
		return status.Error(codes.NotFound, err.Error())
	}
	statusCodes1 := make(chan int32)
	statusCodes2 := make(chan int32)
	statusCodesToSend := make(chan int32)
	statusesToSend := make(chan string)
	availabilitiesToSend := make(chan bool)
	go func() {
		for {
			stCode := <-statusCodesToSend
			status, ok := service.Statuses[stCode]
			if !ok {
				log.WithField("service", service.ServiceIdentifier).Errorf("Got unregistered status code: %d", stCode)
			} else {
				statusesToSend <- status
			}
		}
	}()
	go SendStatuses(s.Slack, service.ServiceIdentifier, statusesToSend)
	go SendAvailabilities(s.Slack, service.ServiceIdentifier, availabilitiesToSend)
	go DetectAvailabilityChange(statusCodes1, availabilitiesToSend, time.Duration(service.Timeout)*time.Second)
	go DetectStatusCodeChange(statusCodes2, statusCodesToSend)
	statusCodes1 <- serviceStatus.StatusCode
	statusCodes2 <- serviceStatus.StatusCode
	for {
		serviceStatus, err := stream.Recv()
		if err != nil {
			close(statusCodes1)
			close(statusCodes2)
			return status.Error(codes.Aborted, err.Error())
		}
		log.WithField("service", service.ServiceIdentifier).Debugf("Received statusCode %d", serviceStatus.StatusCode)
		statusCodes1 <- serviceStatus.StatusCode
		statusCodes2 <- serviceStatus.StatusCode
	}
}

func (s *Server) RegisterService(ctx context.Context, service *pb.Service) (*pb.Empty, error) {
	err := SaveService(s.DB, service)
	log.WithField("service", service.ServiceIdentifier).Info("Registered service")
	return &pb.Empty{}, err
}
