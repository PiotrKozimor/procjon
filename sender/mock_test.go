package sender

import "testing"

func TestMock(t *testing.T) {
	dut := Mock{
		Status:       make(chan string),
		Availability: make(chan string),
		T:            t,
	}
	go func() {
		dut.SendStatus("foo", "bar")
	}()
	msg := <-dut.Status
	if msg != "foo bar" {
		t.Errorf("Expected 'foo bar' status, got %s", msg)
	}
	go func() {
		dut.SendAvailability("foo", false)
	}()
	msg = <-dut.Availability
	if msg != "foo false" {
		t.Errorf("Expected 'foo false' status, got %s", msg)
	}
}
