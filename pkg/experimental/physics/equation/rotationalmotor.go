package equation

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
	"github.com/bhojpur/render/pkg/3d/math32"
)

// RotationalMotor is a rotational motor constraint equation.
// Tries to keep the relative angular velocity of the bodies to a given value.
type RotationalMotor struct {
	Equation                    // TODO maybe this should embed Rotational instead ?
	axisA       *math32.Vector3 // World oriented rotational axis
	axisB       *math32.Vector3 // World oriented rotational axis
	targetSpeed float32         // Target speed
}

// NewRotationalMotor creates and returns a pointer to a new RotationalMotor equation object.
func NewRotationalMotor(bodyA, bodyB IBody, maxForce float32) *RotationalMotor {

	re := new(RotationalMotor)
	re.Equation.initialize(bodyA, bodyB, -maxForce, maxForce)

	return re
}

// SetAxisA sets the axis of body A.
func (ce *RotationalMotor) SetAxisA(axisA *math32.Vector3) {

	ce.axisA = axisA
}

// AxisA returns the axis of body A.
func (ce *RotationalMotor) AxisA() math32.Vector3 {

	return *ce.axisA
}

// SetAxisB sets the axis of body B.
func (ce *RotationalMotor) SetAxisB(axisB *math32.Vector3) {

	ce.axisB = axisB
}

// AxisB returns the axis of body B.
func (ce *RotationalMotor) AxisB() math32.Vector3 {

	return *ce.axisB
}

// SetTargetSpeed sets the target speed.
func (ce *RotationalMotor) SetTargetSpeed(speed float32) {

	ce.targetSpeed = speed
}

// TargetSpeed returns the target speed.
func (ce *RotationalMotor) TargetSpeed() float32 {

	return ce.targetSpeed
}

// ComputeB
func (re *RotationalMotor) ComputeB(h float32) float32 {

	// g = 0
	// gdot = axisA * wi - axisB * wj
	// gdot = G * W = G * [vi wi vj wj]
	// =>
	// G = [0 axisA 0 -axisB]
	re.jeA.SetRotational(re.axisA.Clone())
	re.jeB.SetRotational(re.axisB.Clone().Negate())

	GW := re.ComputeGW() - re.targetSpeed
	GiMf := re.ComputeGiMf()

	return -GW*re.b - h*GiMf
}
