package material

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
)

// Point material is normally used for single point sprites
type Point struct {
	Standard // Embedded standard material
}

// NewPoint creates and returns a pointer to a new point material
func NewPoint(color *math32.Color) *Point {

	pm := new(Point)
	pm.Standard.Init("point", color)

	// Sets uniform's initial values
	pm.udata.emissive = *color
	pm.udata.psize = 1.0
	pm.udata.protationZ = 0
	return pm
}

// SetEmissiveColor sets the material emissive color
// The default is {0,0,0}
func (pm *Point) SetEmissiveColor(color *math32.Color) {

	pm.udata.emissive = *color
}

// SetSize sets the point size
func (pm *Point) SetSize(size float32) {

	pm.udata.psize = size
}

// SetRotationZ sets the point rotation around the Z axis.
func (pm *Point) SetRotationZ(rot float32) {

	pm.udata.protationZ = rot
}
