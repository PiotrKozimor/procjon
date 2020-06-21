package procjon

import (
	"context"
	"log"
	"os"

	"github.com/PiotrKozimor/procjon"
	"github.com/dgraph-io/badger/v2"
)

type server struct {
	slack slack
	db    *badger.DB
}

func (s *server) SendServiceStatus(stream procjon.Procjon_SendServiceStatusServer) error {
	var service procjon.Service
	status, err := stream.Recv()
	if err != nil {
		return err
	}
	err = loadService(s.db, &service, status)
	if err != nil {
		return err
	}
	statusCodes := make(chan int32, 2)
	go processAvailabilityAndStatus(s.slack, statusCodes, &service)
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

func (s *server) RegisterService(ctx context.Context, service *procjon.Service) (*procjon.Empty, error) {
	err := saveService(s.db, service)
	return &procjon.Empty{}, err
}

func main() {
	db, err := badger.Open(badger.DefaultOptions("services"))
	if err != nil {
		log.Fatal(err)
	}
	var s = server{
		slack: slack{webhook: os.Getenv("PROCJON_SLACK_WEBHOOK")},
		db:    db,
	}

}
