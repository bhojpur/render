package shape

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

import "github.com/bhojpur/render/pkg/engine/math32"

// IShape is the interface for all collision shapes.
// Shapes in this package satisfy this interface and also geometry.Geometry.
type IShape interface {
	BoundingBox() math32.Box3
	BoundingSphere() math32.Sphere
	Area() float32
	Volume() float32
	RotationalInertia(mass float32) math32.Matrix3
	ProjectOntoAxis(localAxis *math32.Vector3) (float32, float32)
}

// Shape is a collision shape.
// It can be an analytical geometry such as a sphere, plane, etc.. or it can be defined by a polygonal Geometry.
type Shape struct {

	// TODO
	//material

	// Collision filtering
	colFilterGroup int
	colFilterMask  int
}

func (s *Shape) initialize() {

	// Collision filtering
	s.colFilterGroup = 1
	s.colFilterMask = -1
}
