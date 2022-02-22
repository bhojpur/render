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
	"github.com/bhojpur/render/pkg/engine/math32"
	"github.com/bhojpur/render/pkg/experimental/physics/equation"
)

// PointToPoint is an offset constraint.
// Connects two bodies at the specified offset points.
type PointToPoint struct {
	Constraint
	pivotA *math32.Vector3 // Pivot, defined locally in bodyA.
	pivotB *math32.Vector3 // Pivot, defined locally in bodyB.
	eqX    *equation.Contact
	eqY    *equation.Contact
	eqZ    *equation.Contact
}

// NewPointToPoint creates and returns a pointer to a new PointToPoint constraint object.
func NewPointToPoint(bodyA, bodyB IBody, pivotA, pivotB *math32.Vector3, maxForce float32) *PointToPoint {

	ptpc := new(PointToPoint)
	ptpc.initialize(bodyA, bodyB, pivotA, pivotB, maxForce)

	return ptpc
}

func (ptpc *PointToPoint) initialize(bodyA, bodyB IBody, pivotA, pivotB *math32.Vector3, maxForce float32) {

	ptpc.Constraint.initialize(bodyA, bodyB, true, true)

	ptpc.pivotA = pivotA // default is zero vec3
	ptpc.pivotB = pivotB // default is zero vec3

	ptpc.eqX = equation.NewContact(bodyA, bodyB, -maxForce, maxForce)
	ptpc.eqY = equation.NewContact(bodyA, bodyB, -maxForce, maxForce)
	ptpc.eqZ = equation.NewContact(bodyA, bodyB, -maxForce, maxForce)

	ptpc.eqX.SetNormal(&math32.Vector3{1, 0, 0})
	ptpc.eqY.SetNormal(&math32.Vector3{0, 1, 0})
	ptpc.eqZ.SetNormal(&math32.Vector3{0, 0, 1})

	ptpc.AddEquation(&ptpc.eqX.Equation)
	ptpc.AddEquation(&ptpc.eqY.Equation)
	ptpc.AddEquation(&ptpc.eqZ.Equation)
}

// Update updates the equations with data.
func (ptpc *PointToPoint) Update() {

	// Rotate the pivots to world space
	xRi := ptpc.pivotA.Clone().ApplyQuaternion(ptpc.bodyA.Quaternion())
	xRj := ptpc.pivotA.Clone().ApplyQuaternion(ptpc.bodyA.Quaternion())

	ptpc.eqX.SetRA(xRi)
	ptpc.eqX.SetRB(xRj)
	ptpc.eqY.SetRA(xRi)
	ptpc.eqY.SetRB(xRj)
	ptpc.eqZ.SetRA(xRi)
	ptpc.eqZ.SetRB(xRj)
}
