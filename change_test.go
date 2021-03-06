package procjon

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAvailabilityOnePing(t *testing.T) {
	c := make(chan bool)
	dut := NewAvailability(time.Millisecond*50, func(available bool) {
		c <- available
		t.Logf("%v", available)
	})
	go dut.Run()
	dut.Ping()
	avChange := <-c
	assert.Equal(t, avChange, true)
	ti := time.Now().UnixNano()
	avChange = <-c
	assert.Equal(t, avChange, false)
	assert.Greater(t, time.Now().UnixNano()-ti, int64(5e7))
}

func TestAvailabilityManyPing(t *testing.T) {
	c := make(chan bool)
	dut := NewAvailability(time.Millisecond*50, func(available bool) {
		c <- available
	})
	go dut.Run()
	// time.Sleep(time.Second)
	dut.Ping()
	dut.Ping()
	dut.Ping()
	avChange := <-c
	assert.Equal(t, avChange, true)
	ti := time.Now().UnixNano()
	avChange = <-c
	assert.Equal(t, avChange, false)
	assert.Greater(t, time.Now().UnixNano()-ti, int64(5e7))
}

func TestAvailabilityRecover(t *testing.T) {
	c := make(chan bool)
	dut := NewAvailability(time.Millisecond*50, func(available bool) {
		c <- available
	})
	go dut.Run()
	dut.Ping()
	avChange := <-c
	assert.Equal(t, avChange, true)
	avChange = <-c
	assert.Equal(t, avChange, false)
	dut.Ping()
	avChange = <-c
	assert.Equal(t, avChange, true)
	avChange = <-c
	assert.Equal(t, avChange, false)
	dut.Ping()
	avChange = <-c
	assert.Equal(t, avChange, true)
	dut.Ping()
	avChange = <-c
	assert.Equal(t, avChange, false)
}

func TestStatusCode(t *testing.T) {
	dut := StatusCode{last: 0}
	assert.Equal(t, dut.HasChanged(0), false)
	assert.Equal(t, dut.HasChanged(1), true)
	assert.Equal(t, dut.HasChanged(1000), true)
	assert.Equal(t, dut.HasChanged(1000), false)
	assert.Equal(t, dut.HasChanged(0), true)
	assert.Equal(t, dut.HasChanged(1), true)
}
