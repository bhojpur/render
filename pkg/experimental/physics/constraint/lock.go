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

// Lock constraint.
// Removes all degrees of freedom between the bodies.
type Lock struct {
	PointToPoint
	rotEq1 *equation.Rotational
	rotEq2 *equation.Rotational
	rotEq3 *equation.Rotational
	xA     *math32.Vector3
	xB     *math32.Vector3
	yA     *math32.Vector3
	yB     *math32.Vector3
	zA     *math32.Vector3
	zB     *math32.Vector3
}

// NewLock creates and returns a pointer to a new Lock constraint object.
func NewLock(bodyA, bodyB IBody, maxForce float32) *Lock {

	lc := new(Lock)

	// Set pivot point in between
	posA := bodyA.Position()
	posB := bodyB.Position()

	halfWay := math32.NewVec3().AddVectors(&posA, &posB)
	halfWay.MultiplyScalar(0.5)

	pivotB := bodyB.PointToLocal(halfWay)
	pivotA := bodyA.PointToLocal(halfWay)

	// The point-to-point constraint will keep a point shared between the bodies
	lc.initialize(bodyA, bodyB, &pivotA, &pivotB, maxForce)

	// Store initial rotation of the bodies as unit vectors in the local body spaces
	UnitX := math32.NewVector3(1, 0, 0)

	localA := bodyA.VectorToLocal(UnitX)
	localB := bodyB.VectorToLocal(UnitX)

	lc.xA = &localA
	lc.xB = &localB
	lc.yA = &localA
	lc.yB = &localB
	lc.zA = &localA
	lc.zB = &localB

	// ...and the following rotational equations will keep all rotational DOF's in place

	lc.rotEq1 = equation.NewRotational(bodyA, bodyB, maxForce)
	lc.rotEq2 = equation.NewRotational(bodyA, bodyB, maxForce)
	lc.rotEq3 = equation.NewRotational(bodyA, bodyB, maxForce)

	lc.AddEquation(lc.rotEq1)
	lc.AddEquation(lc.rotEq2)
	lc.AddEquation(lc.rotEq3)

	return lc
}

// Update updates the equations with data.
func (lc *Lock) Update() {

	lc.PointToPoint.Update()

	// These vector pairs must be orthogonal

	xAw := lc.bodyA.VectorToWorld(lc.xA)
	yBw := lc.bodyA.VectorToWorld(lc.yB)

	yAw := lc.bodyA.VectorToWorld(lc.yA)
	zBw := lc.bodyB.VectorToWorld(lc.zB)

	zAw := lc.bodyA.VectorToWorld(lc.zA)
	xBw := lc.bodyB.VectorToWorld(lc.xB)

	lc.rotEq1.SetAxisA(&xAw)
	lc.rotEq1.SetAxisB(&yBw)

	lc.rotEq2.SetAxisA(&yAw)
	lc.rotEq2.SetAxisB(&zBw)

	lc.rotEq3.SetAxisA(&zAw)
	lc.rotEq3.SetAxisB(&xBw)
}
