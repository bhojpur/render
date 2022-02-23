package stats

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
	"runtime"
	"time"

	"github.com/bhojpur/render/pkg/gls"
)

// Stats contains several statistics useful for performance evaluation
type Stats struct {
	gs           *gls.GLS  // Reference to OpenGL state
	Glstats      gls.Stats // GLS statistics structure
	UnilocHits   int       // Uniform location cache hits per frame
	UnilocMiss   int       // Uniform location cache misses per frame
	Unisets      int       // Uniform sets per frame
	Drawcalls    int       // Draw calls per frame
	Cgocalls     int       // Cgo calls per frame
	prevGls      gls.Stats // previous gls statistics
	prevCgocalls int64     // previous number of cgo calls
	frames       int       // frame counter
	last         time.Time // last update time
}

// NewStats creates and returns a pointer to a new statistics object
func NewStats(gs *gls.GLS) *Stats {

	s := new(Stats)
	s.gs = gs
	s.last = time.Now()
	return s
}

// Update should be called in the render loop with the desired update interval.
// Returns true when the interval has elapsed and the statistics has been updated.
func (s *Stats) Update(d time.Duration) bool {

	// If the specified time duration has not elapsed from previous update,
	// nothing to do.
	now := time.Now()
	s.frames++
	if s.last.Add(d).After(now) {
		return false
	}

	// Update GLS statistics
	s.gs.Stats(&s.Glstats)

	// Calculates uniform location cache hits per frame
	unilochits := s.Glstats.UnilocHits - s.prevGls.UnilocHits
	s.UnilocHits = int(float64(unilochits) / float64(s.frames))

	// Calculates uniform location cache hits per frame
	unilocmiss := s.Glstats.UnilocMiss - s.prevGls.UnilocMiss
	s.UnilocMiss = int(float64(unilocmiss) / float64(s.frames))

	// Calculates uniforms sets per frame
	unisets := s.Glstats.Unisets - s.prevGls.Unisets
	s.Unisets = int(float64(unisets) / float64(s.frames))

	// Calculates draw calls per frame
	drawcalls := s.Glstats.Drawcalls - s.prevGls.Drawcalls
	s.Drawcalls = int(float64(drawcalls) / float64(s.frames))

	// Calculates number of cgo calls per frame
	current := runtime.NumCgoCall()
	cgocalls := current - s.prevCgocalls
	s.Cgocalls = int(float64(cgocalls) / float64(s.frames))
	s.prevCgocalls = current

	s.prevGls = s.Glstats
	s.last = now
	s.frames = 0
	return true
}
