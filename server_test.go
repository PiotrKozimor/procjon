package procjon

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/PiotrKozimor/procjon/pb"
	"github.com/PiotrKozimor/procjon/sender"
	"github.com/stretchr/testify/assert"
)

func TestProcjon(t *testing.T) {
	mock := sender.Mock{
		Status:       make(chan string),
		Availability: make(chan string),
		T:            t,
	}
	conn := MustConnectOnBuffer(&mock)
	defer conn.Close()
	client := pb.NewProcjonClient(conn)
	resp, err := client.RegisterService(context.Background(), &pb.Service{
		Identifier: "redis",
		Statuses: []string{
			"ok",
			"nok",
		},
		Timeout: 1,
	})
	if err != nil {
		t.Errorf("RegisterService failed: %v", err)
	}
	log.Printf("Response: %+v", resp)
	stream, err := client.SendServiceStatus(context.Background())
	go func() {
		avTestVector := []string{"redis true", "redis false", "redis true", "redis false"}
		for i := 0; true; i++ {
			av, ok := <-mock.Availability
			if !ok {
				if i != len(avTestVector) {
					t.Errorf("Not all availabilityTestsVector consumed, i: %d", i)
				}
				return
			}
			assert.Equal(t, avTestVector[i], av)
		}
	}()
	go func() {
		stTestVector := []string{"redis nok", "redis ok", "redis nok"}
		for i := 0; true; i++ {
			st, ok := <-mock.Status
			if !ok {
				if i != len(stTestVector) {
					t.Errorf("Not all stTestVector consumed, i: %d", i)
				}
				return
			}
			assert.Equal(t, stTestVector[i], st)
		}
	}()
	inStatusCodes := []uint32{0, 1, 4, 0, 0, 0, 1}
	inDelays := []uint32{50, 50, 50, 50, 2000, 50, 500}
	for i, stC := range inStatusCodes {
		err = stream.Send(&pb.ServiceStatus{Identifier: "redis", StatusCode: stC})
		if err != nil {
			t.Fatalf("Failed to send status: %v", err)
		}
		t.Logf("Sent code: %d", stC)
		time.Sleep(time.Duration(inDelays[i]) * time.Millisecond)
	}
	err = stream.CloseSend()
	if err != nil {
		t.Fatalf("Failed to CloseSend: %v", err)
	}
	time.Sleep(time.Second * 2)
	close(mock.Availability)
	close(mock.Status)
	time.Sleep(time.Second * 1)
}
