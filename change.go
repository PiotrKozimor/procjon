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
// Callback will be called after first Ping() call. We assume initial inavailability.
func NewAvailability(timeout time.Duration, callback func(bool)) *Availability {
	a := Availability{
		refresh:   make(chan bool),
		available: false,
		callback:  callback,
		timeout:   timeout,
		timer:     time.NewTimer(timeout),
	}
	return &a
}

// Ping to renew timeout. Must call Run() before in seperate goroutine.
func (a *Availability) Ping() {
	a.refresh <- true
}

// Run will detect availability changes. Should be run in seperate goroutine.
// Callback will be called in seperate goroutine.
func (a *Availability) Run() {
	for {
		select {
		case <-a.timer.C:
			a.available = false
			go a.callback(a.available)
		case <-a.refresh:
			if !a.available {
				a.available = true
				go a.callback(a.available)
			}
			a.timer = time.NewTimer(a.timeout)
		}
	}
}

type StatusCode struct {
	last uint32
}

// HasChanged returns true if new is different than value from previous call.
func (stc *StatusCode) HasChanged(new uint32) (changed bool) {
	if new != stc.last {
		stc.last = new
		return true
	}
	return false
}
