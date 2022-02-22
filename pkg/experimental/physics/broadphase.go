package physics

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

// It implements collision related algorithms and data structures.

import (
	"github.com/bhojpur/render/pkg/experimental/physics/object"
)

// CollisionPair is a pair of bodies that may be colliding.
type CollisionPair struct {
	BodyA *object.Body
	BodyB *object.Body
}

// Broadphase is the base class for broadphase implementations.
type Broadphase struct{}

// NewBroadphase creates and returns a pointer to a new Broadphase.
func NewBroadphase() *Broadphase {

	b := new(Broadphase)
	return b
}

// FindCollisionPairs (naive implementation)
func (b *Broadphase) FindCollisionPairs(objects []*object.Body) []CollisionPair {

	pairs := make([]CollisionPair, 0)

	for iA, bodyA := range objects {
		for _, bodyB := range objects[iA+1:] {
			if b.NeedTest(bodyA, bodyB) {
				BBa := bodyA.BoundingBox()
				BBb := bodyB.BoundingBox()
				if BBa.IsIntersectionBox(&BBb) {
					pairs = append(pairs, CollisionPair{bodyA, bodyB})
				}
			}
		}
	}

	return pairs
}

func (b *Broadphase) NeedTest(bodyA, bodyB *object.Body) bool {

	if !bodyA.CollidableWith(bodyB) || (bodyA.Sleeping() && bodyB.Sleeping()) {
		return false
	}

	return true
}
