package constraint

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
	"github.com/bhojpur/render/pkg/experimental/physics/equation"
	"github.com/bhojpur/render/pkg/math32"
)

// ConeTwist constraint.
type ConeTwist struct {
	PointToPoint
	axisA      *math32.Vector3 // Rotation axis, defined locally in bodyA.
	axisB      *math32.Vector3 // Rotation axis, defined locally in bodyB.
	coneEq     *equation.Cone
	twistEq    *equation.Rotational
	angle      float32
	twistAngle float32
}

// NewConeTwist creates and returns a pointer to a new ConeTwist constraint object.
func NewConeTwist(bodyA, bodyB IBody, pivotA, pivotB, axisA, axisB *math32.Vector3, angle, twistAngle, maxForce float32) *ConeTwist {

	ctc := new(ConeTwist)

	// Default of pivots and axes should be vec3(0)

	ctc.initialize(bodyA, bodyB, pivotA, pivotB, maxForce)

	ctc.axisA = axisA
	ctc.axisB = axisB
	ctc.axisA.Normalize()
	ctc.axisB.Normalize()

	ctc.angle = angle
	ctc.twistAngle = twistAngle

	ctc.coneEq = equation.NewCone(bodyA, bodyB, ctc.axisA, ctc.axisB, ctc.angle, maxForce)

	ctc.twistEq = equation.NewRotational(bodyA, bodyB, maxForce)
	ctc.twistEq.SetAxisA(ctc.axisA)
	ctc.twistEq.SetAxisB(ctc.axisB)

	// Make the cone equation push the bodies toward the cone axis, not outward
	ctc.coneEq.SetMaxForce(0)
	ctc.coneEq.SetMinForce(-maxForce)

	// Make the twist equation add torque toward the initial position
	ctc.twistEq.SetMaxForce(0)
	ctc.twistEq.SetMinForce(-maxForce)

	ctc.AddEquation(ctc.coneEq)
	ctc.AddEquation(ctc.twistEq)

	return ctc
}

// Update updates the equations with data.
func (ctc *ConeTwist) Update() {

	ctc.PointToPoint.Update()

	// Update the axes to the cone constraint
	worldAxisA := ctc.bodyA.VectorToWorld(ctc.axisA)
	worldAxisB := ctc.bodyB.VectorToWorld(ctc.axisB)

	ctc.coneEq.SetAxisA(&worldAxisA)
	ctc.coneEq.SetAxisB(&worldAxisB)

	// Update the world axes in the twist constraint
	tA, _ := ctc.axisA.RandomTangents()
	worldTA := ctc.bodyA.VectorToWorld(tA)
	ctc.twistEq.SetAxisA(&worldTA)

	tB, _ := ctc.axisB.RandomTangents()
	worldTB := ctc.bodyB.VectorToWorld(tB)
	ctc.twistEq.SetAxisB(&worldTB)

	ctc.coneEq.SetAngle(ctc.angle)
	ctc.twistEq.SetMaxAngle(ctc.twistAngle)
}
