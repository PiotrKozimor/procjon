package procjon

import (
	"context"
	"time"

	"github.com/PiotrKozimor/procjon/pb"
	"github.com/dgraph-io/badger/v2"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	pb.UnimplementedProcjonServer
	Sender Sender
	DB     *badger.DB
}

func (s *Server) SendServiceStatus(stream pb.Procjon_SendServiceStatusServer) error {
	var service pb.Service
	serviceStatus, err := stream.Recv()
	if err != nil {
		return status.Error(codes.Aborted, err.Error())
	}
	log.WithField("service", service.Identifier).Debugf("Received statusCode %d", serviceStatus.StatusCode)
	err = LoadService(s.DB, &service, serviceStatus)
	if err != nil {
		return status.Error(codes.NotFound, err.Error())
	}
	availability := NewAvailability(time.Duration(service.Timeout)*time.Second, func(available bool) {
		s.Sender.SendAvailability(service.Identifier, available)
	})
	go func() {
		availability.Run()
	}()
	statusC := StatusCode{last: 0}
	for {
		availability.Ping()
		if serviceStatus.StatusCode > uint32(len(service.Statuses)) {
			log.WithField("service", service.Identifier).Errorf("Got unregistered status code: %d", serviceStatus.StatusCode)
		} else {
			status := service.Statuses[serviceStatus.StatusCode]
			if statusC.HasChanged(serviceStatus.StatusCode) {
				s.Sender.SendStatus(service.Identifier, status)
			}
		}
		serviceStatus, err = stream.Recv()
		if err != nil {
			return err
		}
		log.WithField("service", service.Identifier).Debugf("Received statusCode %d", serviceStatus.StatusCode)
	}
}

func (s *Server) RegisterService(ctx context.Context, service *pb.Service) (*pb.Empty, error) {
	err := SaveService(s.DB, service)
	log.WithField("service", service.Identifier).Info("Registered service")
	return &pb.Empty{}, err
}
