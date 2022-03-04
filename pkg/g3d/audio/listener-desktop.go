//go:build !wasm
// +build !wasm

package audio

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
	"github.com/bhojpur/render/pkg/g3d/audio/al"
	"github.com/bhojpur/render/pkg/g3d/core"
	"github.com/bhojpur/render/pkg/gls"
	"github.com/bhojpur/render/pkg/math32"
)

// Listener is an audio listener positioned in space.
type Listener struct {
	core.Node
}

// NewListener creates a Listener object.
func NewListener() *Listener {

	l := new(Listener)
	l.Node.Init(l)
	return l
}

// SetVelocity sets the velocity of the listener with x, y, z components.
func (l *Listener) SetVelocity(vx, vy, vz float32) {

	al.Listener3f(al.Velocity, vx, vy, vz)
}

// SetVelocityVec sets the velocity of the listener with a vector.
func (l *Listener) SetVelocityVec(v *math32.Vector3) {

	al.Listener3f(al.Velocity, v.X, v.Y, v.Z)
}

// Velocity returns the velocity of the listener as x, y, z components.
func (l *Listener) Velocity() (float32, float32, float32) {

	return al.GetListener3f(al.Velocity)
}

// VelocityVec returns the velocity of the listener as a vector.
func (l *Listener) VelocityVec() math32.Vector3 {

	vx, vy, vz := al.GetListener3f(al.Velocity)
	return math32.Vector3{vx, vy, vz}
}

// SetGain sets the gain of the listener.
func (l *Listener) SetGain(gain float32) {

	al.Listenerf(al.Gain, gain)
}

// Gain returns the gain of the listener.
func (l *Listener) Gain() float32 {

	return al.GetListenerf(al.Gain)
}

// Render is called by the renderer at each frame.
// Updates the position and orientation of the listener.
func (l *Listener) Render(gl *gls.GLS) {

	// Sets the listener source world position
	var wpos math32.Vector3
	l.WorldPosition(&wpos)
	al.Listener3f(al.Position, wpos.X, wpos.Y, wpos.Z)

	// Get listener current world direction
	var vdir math32.Vector3
	l.WorldDirection(&vdir)

	// Assumes initial UP vector and recalculates current up vector
	vup := math32.Vector3{0, 1, 0}
	var vright math32.Vector3
	vright.CrossVectors(&vdir, &vup)
	vup.CrossVectors(&vright, &vdir)

	// Sets the listener orientation
	orientation := []float32{vdir.X, vdir.Y, vdir.Z, vup.X, vup.Y, vup.Z}
	al.Listenerfv(al.Orientation, orientation)
}
