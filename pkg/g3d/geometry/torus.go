package geometry

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
	"math"

	"github.com/bhojpur/render/pkg/gls"
	"github.com/bhojpur/render/pkg/math32"
)

// NewTorus creates a torus geometry with the specified revolution radius, tube radius,
// number of radial segments, number of tubular segments, and arc length angle in radians.
// TODO instead of 'arc' have thetaStart and thetaLength for consistency with other generators
// TODO then rename this to NewTorusSector and add a NewTorus constructor
func NewTorus(radius, tubeRadius float64, radialSegments, tubularSegments int, arc float64) *Geometry {

	t := NewGeometry()

	// Create buffers
	positions := math32.NewArrayF32(0, 0)
	normals := math32.NewArrayF32(0, 0)
	uvs := math32.NewArrayF32(0, 0)
	indices := math32.NewArrayU32(0, 0)

	var center math32.Vector3
	for j := 0; j <= radialSegments; j++ {
		for i := 0; i <= tubularSegments; i++ {
			u := float64(i) / float64(tubularSegments) * arc
			v := float64(j) / float64(radialSegments) * math.Pi * 2

			center.X = float32(radius * math.Cos(u))
			center.Y = float32(radius * math.Sin(u))

			var vertex math32.Vector3
			vertex.X = float32((radius + tubeRadius*math.Cos(v)) * math.Cos(u))
			vertex.Y = float32((radius + tubeRadius*math.Cos(v)) * math.Sin(u))
			vertex.Z = float32(tubeRadius * math.Sin(v))
			positions.AppendVector3(&vertex)

			uvs.Append(float32(float64(i)/float64(tubularSegments)), float32(float64(j)/float64(radialSegments)))
			normals.AppendVector3(vertex.Sub(&center).Normalize())
		}
	}

	for j := 1; j <= radialSegments; j++ {
		for i := 1; i <= tubularSegments; i++ {
			a := (tubularSegments+1)*j + i - 1
			b := (tubularSegments+1)*(j-1) + i - 1
			c := (tubularSegments+1)*(j-1) + i
			d := (tubularSegments+1)*j + i
			indices.Append(uint32(a), uint32(b), uint32(d), uint32(b), uint32(c), uint32(d))
		}
	}

	t.SetIndices(indices)
	t.AddVBO(gls.NewVBO(positions).AddAttrib(gls.VertexPosition))
	t.AddVBO(gls.NewVBO(normals).AddAttrib(gls.VertexNormal))
	t.AddVBO(gls.NewVBO(uvs).AddAttrib(gls.VertexTexcoord))

	return t
}
