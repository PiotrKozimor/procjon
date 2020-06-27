package procjon

import (
	"testing"
	"time"
)

func TestDetectStatusCodeChange_WithInitialZeroStatusCode(t *testing.T) {
	cCodes := make(chan int32)
	cCodeChanges := make(chan int32)
	inVector := []int32{0, 1, 1, 2, 2, 1, 3, 2}
	outVector := []int32{-1, 1, -1, 2, -1, 1, 3, 2}
	go DetectStatusCodeChange(cCodes, cCodeChanges)
	for i, inVal := range inVector {
		t.Logf("Send: %d", inVal)
		cCodes <- inVal
		out := outVector[i]
		if out != -1 {
			change := <-cCodeChanges
			t.Logf("Received: %d", change)
			if change != out {
				t.Errorf("Got: %d, expected: %d, index: %d", change, out, i)
			}
		}
	}
}

func TestDetectStatusCodeChange_WithInitialNonZeroStatusCode(t *testing.T) {
	cCodes := make(chan int32)
	cCodeChanges := make(chan int32)
	inVector := []int32{1, 1, 1, 2, 2, 1, 3, 2}
	outVector := []int32{1, -1, -1, 2, -1, 1, 3, 2}
	go DetectStatusCodeChange(cCodes, cCodeChanges)
	for i, inVal := range inVector {
		t.Logf("Send: %d", inVal)
		cCodes <- inVal
		out := outVector[i]
		if out != -1 {
			change := <-cCodeChanges
			t.Logf("Received: %d", change)
			if change != out {
				t.Errorf("Got: %d, expected: %d, index: %d", change, out, i)
			}
		}
	}
	close(cCodes)
}
func TestDetectStatusCodeWith_ChannelClose(t *testing.T) {
	cCodes := make(chan int32)
	cCodeChanges := make(chan int32)
	go DetectStatusCodeChange(cCodes, cCodeChanges)
	cCodes <- 1
	close(cCodes)
	recv := <-cCodeChanges
	if recv != 1 {
		t.Errorf("Got 1, wanted: %d", recv)
	}
}

func TestDetectAvailabilityChange(t *testing.T) {
	cCodes := make(chan int32)
	cAvail := make(chan bool)
	go DetectAvailabilityChange(cCodes, cAvail, time.Second)
	// statusCodes := [...]int32{1, 1, 0, 1, 2, 0}
	a := <-cAvail
	if a != true {
		t.Error("Initial availability not true")
	}
	cCodes <- 1
	cCodes <- 1
	close(cCodes)
	a = <-cAvail
	if a != false {
		t.Error("Availability change not detected")
	}
}

func TestDetectAvailabilityChange_WithChannelClose(t *testing.T) {
	cCodes := make(chan int32)
	cAvail := make(chan bool)
	go DetectAvailabilityChange(cCodes, cAvail, time.Second)
	// statusCodes := [...]int32{1, 1, 0, 1, 2, 0}
	a := <-cAvail
	if a != true {
		t.Error("Initial availability not true")
	}
	close(cCodes)
	a = <-cAvail
	if a != false {
		t.Error("Availability change not detected")
	}
}

func TestDetectAvailabilityChange_WithTimeout(t *testing.T) {
	cCodes := make(chan int32)
	cAvail := make(chan bool)
	go DetectAvailabilityChange(cCodes, cAvail, 50*time.Millisecond)
	// statusCodes := [...]int32{1, 1, 0, 1, 2, 0}
	go func() {
		cCodes <- 1
		time.Sleep(100 * time.Millisecond)
		cCodes <- 2
	}()
	a := <-cAvail
	if a != true {
		t.Error("Initial availability not true")
	}
	a = <-cAvail
	if a != false {
		t.Error("Availability change not detected")
	}
	a = <-cAvail
	if a != true {
		t.Error("Availability change not detected")
	}
}
