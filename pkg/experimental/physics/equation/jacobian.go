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

import "github.com/bhojpur/render/pkg/math32"

// JacobianElement contains 6 entries, 3 spatial and 3 rotational degrees of freedom.
type JacobianElement struct {
	spatial    math32.Vector3
	rotational math32.Vector3
}

// SetSpatial sets the spatial component of the JacobianElement.
func (je *JacobianElement) SetSpatial(spatial *math32.Vector3) {

	je.spatial = *spatial
}

// Spatial returns the spatial component of the JacobianElement.
func (je *JacobianElement) Spatial() math32.Vector3 {

	return je.spatial
}

// Rotational sets the rotational component of the JacobianElement.
func (je *JacobianElement) SetRotational(rotational *math32.Vector3) {

	je.rotational = *rotational
}

// Rotational returns the rotational component of the JacobianElement.
func (je *JacobianElement) Rotational() math32.Vector3 {

	return je.rotational
}

// MultiplyElement multiplies the JacobianElement with another JacobianElement.
// None of the elements are changed.
func (je *JacobianElement) MultiplyElement(je2 *JacobianElement) float32 {

	return je.spatial.Dot(&je2.spatial) + je.rotational.Dot(&je2.rotational)
}

// MultiplyElement multiplies the JacobianElement with two vectors.
// None of the elements are changed.
func (je *JacobianElement) MultiplyVectors(spatial *math32.Vector3, rotational *math32.Vector3) float32 {

	return je.spatial.Dot(spatial) + je.rotational.Dot(rotational)
}
