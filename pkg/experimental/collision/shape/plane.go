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

import "github.com/bhojpur/render/pkg/math32"

// Plane is an analytical collision Plane.
// A plane, facing in the +Z direction. The plane has its surface at z=0 and everything below z=0 is assumed to be solid.
type Plane struct {
	normal math32.Vector3
}

// NewPlane creates and returns a pointer to a new analytical collision plane.
func NewPlane() *Plane {

	p := new(Plane)
	p.normal = *math32.NewVector3(0, 0, 1) //*normal
	return p
}

// SetRadius sets the radius of the analytical collision sphere.
//func (p *Plane) SetNormal(normal *math32.Vector3) {
//
//	p.normal = *normal
//}

// Normal returns the normal of the analytical collision plane.
func (p *Plane) Normal() math32.Vector3 {

	return p.normal
}

// IShape =============================================================

// BoundingBox computes and returns the bounding box of the analytical collision plane.
func (p *Plane) BoundingBox() math32.Box3 {

	//return math32.Box3{math32.Vector3{math32.Inf(-1), math32.Inf(-1), math32.Inf(-1)}, math32.Vector3{math32.Inf(1), 0, math32.Inf(1)}}
	return math32.Box3{math32.Vector3{-1000, -1000, -1000}, math32.Vector3{1000, 1000, 0}}
}

// BoundingSphere computes and returns the bounding sphere of the analytical collision plane.
func (p *Plane) BoundingSphere() math32.Sphere {

	return *math32.NewSphere(math32.NewVec3(), math32.Inf(1))
}

// Area returns the surface area of the analytical collision plane.
func (p *Plane) Area() float32 {

	return math32.Inf(1)
}

// Volume returns the volume of the analytical collision sphere.
func (p *Plane) Volume() float32 {

	return math32.Inf(1)
}

// RotationalInertia computes and returns the rotational inertia of the analytical collision plane.
func (p *Plane) RotationalInertia(mass float32) math32.Matrix3 {

	return *math32.NewMatrix3().Zero()
}

// ProjectOntoAxis returns the minimum and maximum distances of the analytical collision plane projected onto the specified local axis.
func (p *Plane) ProjectOntoAxis(localAxis *math32.Vector3) (float32, float32) {

	if localAxis.Equals(&p.normal) {
		return math32.Inf(-1), 0
	} else {
		return math32.Inf(-1), math32.Inf(1)
	}
}
