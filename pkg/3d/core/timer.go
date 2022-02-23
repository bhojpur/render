package core

// Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

import (
	"time"
)

// TimerManager manages multiple timers
type TimerManager struct {
	nextID int       // next timer id
	timers []timeout // list of timeouts
}

// TimerCallback is the type for timer callback functions
type TimerCallback func(interface{})

// Internal structure for each active timer
type timeout struct {
	id     int           // timeout id
	expire time.Time     // expiration time
	period time.Duration // period time
	cb     TimerCallback // callback function
	arg    interface{}   // callback function argument
}

// NewTimerManager creates and returns a new timer manager
func NewTimerManager() *TimerManager {

	tm := new(TimerManager)
	tm.Initialize()
	return tm
}

// Initialize initializes the timer manager.
// It is normally used when the TimerManager is embedded in another type.
func (tm *TimerManager) Initialize() {

	tm.nextID = 1
	tm.timers = make([]timeout, 0)
}

// SetTimeout sets a timeout with the specified duration and callback
// The function returns the timeout id which can be used to cancel the timeout
func (tm *TimerManager) SetTimeout(td time.Duration, arg interface{}, cb TimerCallback) int {

	return tm.setTimer(td, false, arg, cb)
}

// SetInterval sets a periodic timeout with the specified duration and callback
// The function returns the timeout id which can be used to cancel the timeout
func (tm *TimerManager) SetInterval(td time.Duration, arg interface{}, cb TimerCallback) int {

	return tm.setTimer(td, true, arg, cb)
}

// ClearTimeout clears the timeout specified by the id.
// Returns true if the timeout is found.
func (tm *TimerManager) ClearTimeout(id int) bool {

	for pos, t := range tm.timers {
		if t.id == id {
			copy(tm.timers[pos:], tm.timers[pos+1:])
			tm.timers[len(tm.timers)-1] = timeout{}
			tm.timers = tm.timers[:len(tm.timers)-1]
			return true
		}
	}
	return false
}

// ProcessTimers should be called periodically to process the timers
func (tm *TimerManager) ProcessTimers() {

	now := time.Now()
	for pos, t := range tm.timers {
		// If empty entry, ignore
		if t.id == 0 {
			continue
		}
		// Checks if entry expired
		if now.After(t.expire) {
			if t.period == 0 {
				tm.timers[pos] = timeout{}
			} else {
				tm.timers[pos].expire = now.Add(t.period)
			}
			t.cb(t.arg)
		}
	}
}

// setTimer sets a new timer with the specified duration
func (tm *TimerManager) setTimer(td time.Duration, periodic bool, arg interface{}, cb TimerCallback) int {

	// Creates timeout entry
	t := timeout{
		id:     tm.nextID,
		expire: time.Now().Add(td),
		cb:     cb,
		arg:    arg,
		period: 0,
	}
	if periodic {
		t.period = td
	}
	tm.nextID++

	// Look for empty entry
	for pos, ct := range tm.timers {
		if ct.id == 0 {
			tm.timers[pos] = t
			return t.id
		}
	}

	// If no empty entry found, add to end of array
	tm.timers = append(tm.timers, t)
	return t.id
}
