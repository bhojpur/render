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

// It implements physics constraints.

import (
	"github.com/bhojpur/render/pkg/3d/math32"
	"github.com/bhojpur/render/pkg/experimental/physics/equation"
)

type IBody interface {
	equation.IBody
	WakeUp()
	VectorToWorld(*math32.Vector3) math32.Vector3
	PointToLocal(*math32.Vector3) math32.Vector3
	VectorToLocal(*math32.Vector3) math32.Vector3
	Quaternion() *math32.Quaternion
}

type IConstraint interface {
	Update() // Update all the equations with data.
	Equations() []equation.IEquation
	CollideConnected() bool
	BodyA() IBody
	BodyB() IBody
}

// Constraint base struct.
type Constraint struct {
	equations []equation.IEquation // Equations to be solved in this constraint
	bodyA     IBody
	bodyB     IBody
	colConn   bool // Set to true if you want the bodies to collide when they are connected.
}

// NewConstraint creates and returns a pointer to a new Constraint object.
//func NewConstraint(bodyA, bodyB IBody, colConn, wakeUpBodies bool) *Constraint {
//
//	c := new(Constraint)
//	c.initialize(bodyA, bodyB, colConn, wakeUpBodies)
//	return c
//}

func (c *Constraint) initialize(bodyA, bodyB IBody, colConn, wakeUpBodies bool) {

	c.bodyA = bodyA
	c.bodyB = bodyB
	c.colConn = colConn // true

	if wakeUpBodies { // true
		if bodyA != nil {
			bodyA.WakeUp()
		}
		if bodyB != nil {
			bodyB.WakeUp()
		}
	}
}

// AddEquation adds an equation to the constraint.
func (c *Constraint) AddEquation(eq equation.IEquation) {

	c.equations = append(c.equations, eq)
}

// Equations returns the constraint's equations.
func (c *Constraint) Equations() []equation.IEquation {

	return c.equations
}

func (c *Constraint) CollideConnected() bool {

	return c.colConn
}

func (c *Constraint) BodyA() IBody {

	return c.bodyA
}

func (c *Constraint) BodyB() IBody {

	return c.bodyB
}

// SetEnable sets the enabled flag of the constraint equations.
func (c *Constraint) SetEnabled(state bool) {

	for i := range c.equations {
		c.equations[i].SetEnabled(state)
	}
}
