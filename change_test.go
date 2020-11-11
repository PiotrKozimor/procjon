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
	})
	go func() {
		dut.Run()
	}()
	dut.Ping()
	ti := time.Now().UnixNano()
	avChange := <-c
	assert.Equal(t, avChange, false)
	assert.Greater(t, time.Now().UnixNano()-ti, int64(5e7))
}

func TestAvailabilityManyPing(t *testing.T) {
	c := make(chan bool)
	dut := NewAvailability(time.Millisecond*50, func(available bool) {
		c <- available
	})
	go func() {
		dut.Run()
	}()
	dut.Ping()
	dut.Ping()
	dut.Ping()
	ti := time.Now().UnixNano()
	avChange := <-c
	assert.Equal(t, avChange, false)
	assert.Greater(t, time.Now().UnixNano()-ti, int64(5e7))
}

func TestAvailabilityRecover(t *testing.T) {
	c := make(chan bool)
	dut := NewAvailability(time.Millisecond*50, func(available bool) {
		c <- available
	})
	go func() {
		dut.Run()
	}()
	dut.Ping()
	avChange := <-c
	assert.Equal(t, avChange, false)
	dut.Ping()
	avChange = <-c
	assert.Equal(t, avChange, true)
	avChange = <-c
	assert.Equal(t, avChange, false)
	dut.Ping()
	avChange = <-c
	assert.Equal(t, avChange, true)
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
