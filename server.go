package main

import (
	"context"
	"log"
	"os"

	"github.com/PiotrKozimor/procjon/pb"
	"github.com/PiotrKozimor/procjon/procjon"
	"github.com/dgraph-io/badger/v2"
)

type server struct {
	slack procjon.Slack
	db    *badger.DB
}

func (s *server) SendServiceStatus(stream pb.Procjon_SendServiceStatusServer) error {
	var service pb.Service
	status, err := stream.Recv()
	if err != nil {
		return err
	}
	err = procjon.LoadService(s.db, &service, status)
	if err != nil {
		return err
	}
	statusCodes := make(chan int32, 2)
	go procjon.ProcessAvailabilityAndStatus(s.slack, statusCodes, &service)
	statusCodes <- status.StatusCode
	for {
		status, err := stream.Recv()
		if err != nil {
			close(statusCodes)
			return err
		}
		statusCodes <- status.StatusCode
	}
}

func (s *server) RegisterService(ctx context.Context, service *pb.Service) (*pb.Empty, error) {
	err := procjon.SaveService(s.db, service)
	return &pb.Empty{}, err
}

func main() {
	db, err := badger.Open(badger.DefaultOptions("services"))
	if err != nil {
		log.Fatal(err)
	}
	var s = server{
		slack: procjon.Slack{Webhook: os.Getenv("PROCJON_SLACK_WEBHOOK")},
		db:    db,
	}

}
