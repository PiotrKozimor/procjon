package procjon

import (
	"time"

	"github.com/PiotrKozimor/procjon/pb"
)

func ProcessAvailabilityAndStatus(s Slack, statusCode chan int32, service *pb.Service) {
	timerDuration := time.Second * time.Duration(service.Timeout)
	t := time.NewTimer(timerDuration)
	c := make(chan int32)
	go processStatus(s, c, service)
	select {
	case receivedStatusCode, ok := <-statusCode:
		if !ok {
			s.SendAvailability(service.ServiceIdentifier, false)
			break
		} else {
			t.Reset(timerDuration)
			c <- receivedStatusCode
		}
	case <-t.C:
		s.SendAvailability(service.ServiceIdentifier, false)
	}
}

func processStatus(s Slack, statusCode chan int32, service *pb.Service) {
	lastStatusCode := <-statusCode
	if lastStatusCode != 0 {
		s.SendStatus(service.ServiceIdentifier, service.Statuses[lastStatusCode])
	}
	for {
		stCode := <-statusCode
		if stCode != lastStatusCode {
			s.SendStatus(service.ServiceIdentifier, service.Statuses[stCode])
		}
		lastStatusCode = stCode
	}

}
