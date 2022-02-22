package gls

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

// It implements a loader of OpenGL functions for the platform
// and a Go binding for selected OpenGL functions. The binding maintains
// some cached state to minimize the number of C function calls.
// The OpenGL function loader is generated by the "renderglapi" tool by
// parsing the OpenGL "glcorearb.h" header file
//
// This package also contains abstractions for some OpenGL object such as Program,
// Uniform, VBO and others.

import (
	"math"
	"unsafe"

	"github.com/bhojpur/render/pkg/engine/util/logger"
)

// Package logger
var log = logger.New("GLS", logger.Default)

// Stats contains counters of WebGL resources being used as well
// the cumulative numbers of some WebGL calls for performance evaluation.
type Stats struct {
	Shaders    int    // Current number of shader programs
	Vaos       int    // Number of Vertex Array Objects
	Buffers    int    // Number of Buffer Objects
	Textures   int    // Number of Textures
	Caphits    uint64 // Cumulative number of hits for Enable/Disable
	UnilocHits uint64 // Cumulative number of uniform location cache hits
	UnilocMiss uint64 // Cumulative number of uniform location cache misses
	Unisets    uint64 // Cumulative number of uniform sets
	Drawcalls  uint64 // Cumulative number of draw calls
	Fbos       uint64 // Number of frame buffer objects
	Rbos       uint64 // Number of render buffer objects
}

const (
	capUndef    = 0
	capDisabled = 1
	capEnabled  = 2
	uintUndef   = math.MaxUint32
	intFalse    = 0
	intTrue     = 1
)

const (
	FloatSize = int32(unsafe.Sizeof(float32(0)))
)
