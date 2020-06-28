package main

import (
	"context"
	"net"
	"os"
	"time"

	"github.com/PiotrKozimor/procjon/pb"
	"github.com/PiotrKozimor/procjon/procjon"
	"github.com/dgraph-io/badger/v2"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var logger = log.New()

func init() {
	logger.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.Stamp})
	logger.SetOutput(os.Stderr)
	logger.SetLevel(log.InfoLevel)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.Stamp})
	log.SetOutput(os.Stderr)
	log.SetLevel(log.InfoLevel)
}

type server struct {
	pb.UnimplementedProcjonServer
	slack  procjon.Slack
	db     *badger.DB
	server pb.UnimplementedProcjonServer
}

func (s *server) SendServiceStatus(stream pb.Procjon_SendServiceStatusServer) error {
	var service pb.Service
	serviceStatus, err := stream.Recv()
	if err != nil {
		return status.Error(codes.Aborted, err.Error())
	}
	log.Debugf("Received statusCode %d", serviceStatus.StatusCode)
	err = procjon.LoadService(s.db, &service, serviceStatus)
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
			statusesToSend <- service.Statuses[stCode]
		}
	}()
	go procjon.SendStatuses(&s.slack, service.ServiceIdentifier, statusesToSend)
	go procjon.SendAvailabilities(&s.slack, service.ServiceIdentifier, availabilitiesToSend)
	go procjon.DetectAvailabilityChange(statusCodes1, availabilitiesToSend, time.Duration(service.Timeout)*time.Second)
	go procjon.DetectStatusCodeChange(statusCodes2, statusCodesToSend)
	statusCodes1 <- serviceStatus.StatusCode
	statusCodes2 <- serviceStatus.StatusCode
	for {
		serviceStatus, err := stream.Recv()
		if err != nil {
			close(statusCodes1)
			close(statusCodes2)
			return status.Error(codes.Aborted, err.Error())
		}
		log.Debugf("Received statusCode %d", serviceStatus.StatusCode)
		statusCodes1 <- serviceStatus.StatusCode
		statusCodes2 <- serviceStatus.StatusCode
	}
}

func (s *server) RegisterService(ctx context.Context, service *pb.Service) (*pb.Empty, error) {
	err := procjon.SaveService(s.db, service)
	log.WithField("service", service.ServiceIdentifier).Info("Registered service")
	return &pb.Empty{}, err
}

func main() {

	db, err := badger.Open(badger.DefaultOptions("services").WithLogger(logger))
	if err != nil {
		log.Fatal(err)
	}
	var s = server{
		slack: procjon.Slack{Webhook: os.Getenv("PROCJON_SLACK_WEBHOOK")},
		db:    db,
	}
	grpcServer := grpc.NewServer()
	pb.RegisterProcjonServer(grpcServer, &s)
	lis, err := net.Listen("unix", "procjon.sock")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer lis.Close()
	grpcServer.Serve(lis)
}
