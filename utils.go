package procjon

import (
	"context"
	"log"
	"net"

	"github.com/PiotrKozimor/procjon/pb"
	"github.com/dgraph-io/badger/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

// type MockSender struct {
// 	t   *testing.T
// 	avC map[string]chan bool
// 	stC map[string]chan string
// }

// func (s *MockSender) SendAvailability(service string, availability bool) error {
// 	s.t.Logf("Service: %s, availability: %t", service, availability)
// 	avC[service] <- availability
// 	return nil
// }

// func (s *MockSender) SendStatus(service string, status string) error {
// 	s.t.Logf("Service: %s, status: %s", service, status)
// 	stC[service] <- status
// 	return nil
// }

func MustConnectOnBuffer(sender Sender) *grpc.ClientConn {
	db, err := badger.Open(badger.DefaultOptions("").WithInMemory(true))
	if err != nil {
		log.Fatal(err)
	}
	var s = Server{
		Sender: sender,
		DB:     db,
	}
	lis := bufconn.Listen(bufSize)
	grpcServer := grpc.NewServer()
	pb.RegisterProcjonServer(grpcServer, &s)
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
	conn, err := grpc.DialContext(context.Background(), "bufnet", grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to dial bufnet: %v", err)
	}
	return conn
}
