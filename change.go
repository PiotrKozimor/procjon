package procjon

import (
	"time"
)

// When Run() is called, callback will be called with false when timeout has expired since last Ping.
// When pinged after timeout occured. callback will be called with true.
type Availability struct {
	timer     *time.Timer
	refresh   chan bool
	timeout   time.Duration
	available bool
	callback  func(bool)
}

// NewAvailability creates new Availability object with provided timeout and callback.
// Use av.Run() function to start detecting availability changes.
// Call av.Ping() to renew timeout.
func NewAvailability(timeout time.Duration, callback func(bool)) *Availability {
	a := Availability{
		refresh:   make(chan bool, 1),
		available: true,
		callback:  callback,
		timeout:   timeout,
	}
	return &a
}

// Ping to renew timeout. Must call Run() before in seperate goroutine.
func (a *Availability) Ping() {
	a.refresh <- true
	if !a.available {
		go a.callback(true)
		a.timer.Reset(a.timeout)
	}
}

// Run will detect availability changes. Should be run in seperate goroutine.
func (a *Availability) Run() {
	a.timer = time.NewTimer(a.timeout)
	for {
		select {
		case <-a.timer.C:
			a.available = false
			a.callback(false)
		case <-a.refresh:
			if !a.timer.Stop() {
				<-a.timer.C
			}
			a.timer.Reset(a.timeout)
		}
	}
}

type StatusCode struct {
	last int32
}

// HasChanged returns true if new is different than value from previous call.
func (stc *StatusCode) HasChanged(new int32) (changed bool) {
	if new != stc.last {
		stc.last = new
		return true
	}
	return false
}
