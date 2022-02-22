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

import (
	"fmt"
)

// Uniform represents an OpenGL uniform.
type Uniform struct {
	name      string // base name
	nameIdx   string // cached indexed name
	handle    uint32 // program handle
	location  int32  // last cached location
	lastIndex int32  // last index
}

// Init initializes this uniform location cache and sets its name.
func (u *Uniform) Init(name string) {

	u.name = name
	u.handle = 0     // invalid program handle
	u.location = -1  // invalid location
	u.lastIndex = -1 // invalid index
}

// Name returns the uniform name.
func (u *Uniform) Name() string {

	return u.name
}

// Location returns the location of this uniform for the current shader program.
// The returned location can be -1 if not found.
func (u *Uniform) Location(gs *GLS) int32 {

	handle := gs.prog.Handle()
	if handle != u.handle {
		u.location = gs.prog.GetUniformLocation(u.name)
		u.handle = handle
	}
	return u.location
}

// LocationIdx returns the location of this indexed uniform for the current shader program.
// The returned location can be -1 if not found.
func (u *Uniform) LocationIdx(gs *GLS, idx int32) int32 {

	if idx != u.lastIndex {
		u.nameIdx = fmt.Sprintf("%s[%d]", u.name, idx)
		u.lastIndex = idx
		u.handle = 0
	}
	handle := gs.prog.Handle()
	if handle != u.handle {
		u.location = gs.prog.GetUniformLocation(u.nameIdx)
		u.handle = handle
	}
	return u.location
}
