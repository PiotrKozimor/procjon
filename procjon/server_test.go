package procjon

import (
	"context"
	"log"
	"net"
	"testing"
	"time"

	"github.com/PiotrKozimor/procjon/pb"
	"github.com/dgraph-io/badger/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

var (
	lis *bufconn.Listener
	avC = make(map[string]chan bool)
	stC = make(map[string]chan string)
)

const bufSize = 1024 * 1024

type MockSlackSender struct {
}

func (s *MockSlackSender) SendAvailability(service string, availability bool) error {
	log.Printf("Service: %s, availability: %t", service, availability)
	avC[service] <- availability
	return nil
}

func (s *MockSlackSender) SendStatus(service string, status string) error {
	log.Printf("Service: %s, status: %s", service, status)
	stC[service] <- status
	return nil
}

func init() {
	db, err := badger.Open(badger.DefaultOptions("").WithInMemory(true))
	var s = Server{
		Slack: &MockSlackSender{},
		DB:    db,
	}
	if err != nil {
		log.Fatal(err)
	}
	lis = bufconn.Listen(bufSize)
	grpcServer := grpc.NewServer()
	pb.RegisterProcjonServer(grpcServer, &s)
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()

}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func TestProcjon(t *testing.T) {
	avC["redis"] = make(chan bool)
	stC["redis"] = make(chan string)
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	client := pb.NewProcjonClient(conn)
	resp, err := client.RegisterService(ctx, &pb.Service{
		ServiceIdentifier: "redis",
		Statuses: map[int32]string{
			0: "ok",
			1: "nok",
		},
		Timeout: 1,
	})
	if err != nil {
		t.Errorf("RegisterService failed: %v", err)
	}
	log.Printf("Response: %+v", resp)
	stream, err := client.SendServiceStatus(context.Background())
	go func() {
		avTestVector := []bool{true, false, true, false}
		for i := 0; true; i++ {
			av, ok := <-avC["redis"]
			if !ok {
				if i != len(avTestVector) {
					t.Errorf("Not all availabilityTestsVector consumed, i: %d", i)
				}
				return
			}
			t.Logf("Got availability: %t", av)
			if av != avTestVector[i] {
				t.Errorf("Got: %t, wanted: %t", av, avTestVector[i])
			}
		}
	}()
	go func() {
		stTestVector := []string{"nok", "ok", "nok"}
		for i := 0; true; i++ {
			st, ok := <-stC["redis"]
			if !ok {
				if i != len(stTestVector) {
					t.Errorf("Not all stTestVector consumed, i: %d", i)
				}
				return
			}
			t.Logf("Got status: %s", st)
			if st != stTestVector[i] {
				t.Errorf("Got: %s, wanted: %s", st, stTestVector[i])
			}
		}
	}()
	inStatusCodes := []int32{0, 1, 4, 0, 0, 0, 1}
	inDelays := []int32{50, 50, 50, 50, 2000, 50, 500}
	for i, stC := range inStatusCodes {
		err = stream.Send(&pb.ServiceStatus{ServiceIdentifier: "redis", StatusCode: stC})
		if err != nil {
			t.Fatalf("Failed to send status: %v", err)
		}
		time.Sleep(time.Duration(inDelays[i]) * time.Millisecond)
	}
	err = stream.CloseSend()
	if err != nil {
		t.Fatalf("Failed to CloseSend: %v", err)
	}
	time.Sleep(time.Second)
	close(avC["redis"])
	close(stC["redis"])
}
