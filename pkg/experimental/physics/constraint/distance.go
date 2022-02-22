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
)

// Distance is a distance constraint.
// Constrains two bodies to be at a constant distance from each others center of mass.
type Distance struct {
	Constraint
	distance float32 // Distance
	equation *equation.Contact
}

// NewDistance creates and returns a pointer to a new Distance constraint object.
func NewDistance(bodyA, bodyB IBody, distance, maxForce float32) *Distance {

	dc := new(Distance)
	dc.initialize(bodyA, bodyB, true, true)

	// Default distance should be: bodyA.position.distanceTo(bodyB.position)
	// Default maxForce should be: 1e6

	dc.distance = distance

	dc.equation = equation.NewContact(bodyA, bodyB, -maxForce, maxForce) // Make it bidirectional
	dc.AddEquation(dc.equation)

	return dc
}

// Update updates the equation with data.
func (dc *Distance) Update() {

	halfDist := dc.distance * 0.5

	posA := dc.bodyA.Position()
	posB := dc.bodyB.Position()

	normal := posB.Sub(&posA)
	normal.Normalize()

	dc.equation.SetRA(normal.Clone().MultiplyScalar(halfDist))
	dc.equation.SetRB(normal.Clone().MultiplyScalar(-halfDist))
}
