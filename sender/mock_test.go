package sender

import "testing"

func TestMock(t *testing.T) {
	dut := Mock{
		C: make(chan string),
	}
	go func() {
		dut.SendStatus("foo", "bar")
	}()
	msg := <-dut.C
	if msg != "foo bar" {
		t.Errorf("Expected 'foo bar' status, got %s", msg)
	}
	go func() {
		dut.SendAvailability("foo", false)
	}()
	msg = <-dut.C
	if msg != "foo false" {
		t.Errorf("Expected 'foo false' status, got %s", msg)
	}
}
