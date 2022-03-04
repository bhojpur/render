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

// NewSphere creates a sphere geometry with the specified radius and number of radial segments in each dimension.
func NewSphere(radius float64, widthSegments, heightSegments int) *Geometry {
	return NewSphereSector(radius, widthSegments, heightSegments, 0, math.Pi*2, 0, math.Pi)
}

// NewSphereSector creates a sphere sector geometry with the specified radius, number of radial segments in each dimension, elevation
// start angle in radians, elevation size angle in radians, sector start angle in radians, and sector size angle in radians.
func NewSphereSector(radius float64, widthSegments, heightSegments int, phiStart, phiLength, thetaStart, thetaLength float64) *Geometry {

	s := NewGeometry()

	thetaEnd := thetaStart + thetaLength
	vertexCount := (widthSegments + 1) * (heightSegments + 1)

	// Create buffers
	positions := math32.NewArrayF32(vertexCount*3, vertexCount*3)
	normals := math32.NewArrayF32(vertexCount*3, vertexCount*3)
	uvs := math32.NewArrayF32(vertexCount*2, vertexCount*2)
	indices := math32.NewArrayU32(0, vertexCount)

	index := 0
	vertices := make([][]uint32, 0)
	var normal math32.Vector3

	for y := 0; y <= heightSegments; y++ {
		verticesRow := make([]uint32, 0)
		v := float64(y) / float64(heightSegments)
		for x := 0; x <= widthSegments; x++ {
			u := float64(x) / float64(widthSegments)
			px := -radius * math.Cos(phiStart+u*phiLength) * math.Sin(thetaStart+v*thetaLength)
			py := radius * math.Cos(thetaStart+v*thetaLength)
			pz := radius * math.Sin(phiStart+u*phiLength) * math.Sin(thetaStart+v*thetaLength)
			normal.Set(float32(px), float32(py), float32(pz)).Normalize()

			positions.Set(index*3, float32(px), float32(py), float32(pz))
			normals.SetVector3(index*3, &normal)
			uvs.Set(index*2, float32(u), float32(v))
			verticesRow = append(verticesRow, uint32(index))
			index++
		}
		vertices = append(vertices, verticesRow)
	}

	for y := 0; y < heightSegments; y++ {
		for x := 0; x < widthSegments; x++ {
			v1 := vertices[y][x+1]
			v2 := vertices[y][x]
			v3 := vertices[y+1][x]
			v4 := vertices[y+1][x+1]
			if y != 0 || thetaStart > 0 {
				indices.Append(v1, v2, v4)
			}
			if y != heightSegments-1 || thetaEnd < math.Pi {
				indices.Append(v2, v3, v4)
			}
		}
	}

	s.SetIndices(indices)
	s.AddVBO(gls.NewVBO(positions).AddAttrib(gls.VertexPosition))
	s.AddVBO(gls.NewVBO(normals).AddAttrib(gls.VertexNormal))
	s.AddVBO(gls.NewVBO(uvs).AddAttrib(gls.VertexTexcoord))

	r := float32(radius)

	// Update bounding sphere
	s.boundingSphere.Radius = 3
	s.boundingSphereValid = true

	// Update bounding box
	s.boundingBox = math32.Box3{math32.Vector3{-r, -r, -r}, math32.Vector3{r, r, r}}
	s.boundingBoxValid = true

	return s
}
