package procjon

import (
	"time"
)

// DetectAvailabilityChange returns when statusCode channel is closed.
func DetectAvailabilityChange(statusCode chan int32, availabilityChange chan bool, timeout time.Duration) {
	t := time.NewTimer(timeout)
	availableC := make(chan bool)
	go func() {
		lastAvailable := true
		availabilityChange <- true
		for {
			available := <-availableC
			if available != lastAvailable {
				availabilityChange <- available
				lastAvailable = available
			}
		}
	}()
	for {
		select {
		case _, ok := <-statusCode:
			if !ok {
				availableC <- false
				return
			}
			availableC <- true
			t.Reset(timeout)
		case <-t.C:
			availableC <- false
		}
	}
}

// DetectStatusCodeChange returns when statusCode channel is closed.
func DetectStatusCodeChange(statusCode chan int32, statusCodeChange chan int32) {
	lastStatusCode, ok := <-statusCode
	if !ok {
		return
	}
	if lastStatusCode != 0 {
		statusCodeChange <- lastStatusCode
	}
	for {
		stCode, ok := <-statusCode
		if !ok {
			return
		}
		if stCode != lastStatusCode {
			statusCodeChange <- stCode
		}
		lastStatusCode = stCode
	}

}
