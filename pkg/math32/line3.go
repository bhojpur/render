package math32

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

// Line3 represents a 3D line segment defined by a start and an end point.
type Line3 struct {
	start Vector3
	end   Vector3
}

// NewLine3 creates and returns a pointer to a new Line3 with the
// specified start and end points.
func NewLine3(start, end *Vector3) *Line3 {

	l := new(Line3)
	l.Set(start, end)
	return l
}

// Set sets this line segment start and end points.
// Returns pointer to this updated line segment.
func (l *Line3) Set(start, end *Vector3) *Line3 {

	if start != nil {
		l.start = *start
	}
	if end != nil {
		l.end = *end
	}
	return l
}

// Copy copy other line segment to this one.
// Returns pointer to this updated line segment.
func (l *Line3) Copy(other *Line3) *Line3 {

	*l = *other
	return l
}

// Center calculates this line segment center point.
// Store its pointer into optionalTarget, if not nil, and also returns it.
func (l *Line3) Center(optionalTarget *Vector3) *Vector3 {

	var result *Vector3
	if optionalTarget == nil {
		result = NewVector3(0, 0, 0)
	} else {
		result = optionalTarget
	}
	return result.AddVectors(&l.start, &l.end).MultiplyScalar(0.5)
}

// Delta calculates the vector from the start to end point of this line segment.
// Store its pointer in optionalTarget, if not nil, and also returns it.
func (l *Line3) Delta(optionalTarget *Vector3) *Vector3 {

	var result *Vector3
	if optionalTarget == nil {
		result = NewVector3(0, 0, 0)
	} else {
		result = optionalTarget
	}
	return result.SubVectors(&l.end, &l.start)
}

// DistanceSq returns the square of the distance from the start point to the end point.
func (l *Line3) DistanceSq() float32 {

	return l.start.DistanceToSquared(&l.end)
}

// Distance returns the distance from the start point to the end point.
func (l *Line3) Distance() float32 {

	return l.start.DistanceTo(&l.end)
}

// ApplyMatrix4 applies the specified matrix to this line segment start and end points.
// Returns pointer to this updated line segment.
func (l *Line3) ApplyMatrix4(matrix *Matrix4) *Line3 {

	l.start.ApplyMatrix4(matrix)
	l.end.ApplyMatrix4(matrix)
	return l
}

// Equals returns if this line segement is equal to other.
func (l *Line3) Equals(other *Line3) bool {

	return other.start.Equals(&l.start) && other.end.Equals(&l.end)
}

// Clone creates and returns a pointer to a copy of this line segment.
func (l *Line3) Clone() *Line3 {

	return NewLine3(&l.start, &l.end)
}
