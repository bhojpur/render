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
	"github.com/bhojpur/render/pkg/3d/math32"
	"github.com/bhojpur/render/pkg/gls"
)

// NewPlane creates a plane geometry with the specified width and height.
// The plane is generated centered in the XY plane with Z=0.
func NewPlane(width, height float32) *Geometry {
	return NewSegmentedPlane(width, height, 1, 1)
}

// NewSegmentedPlane creates a segmented plane geometry with the specified width, height, and number of
// segments in each dimension (minimum 1 in each). The plane is generated centered in the XY plane with Z=0.
func NewSegmentedPlane(width, height float32, widthSegments, heightSegments int) *Geometry {

	plane := NewGeometry()

	widthHalf := width / 2
	heightHalf := height / 2
	gridX := widthSegments
	gridY := heightSegments
	gridX1 := gridX + 1
	gridY1 := gridY + 1
	segmentWidth := width / float32(gridX)
	segmentHeight := height / float32(gridY)

	// Create buffers
	positions := math32.NewArrayF32(0, 16)
	normals := math32.NewArrayF32(0, 16)
	uvs := math32.NewArrayF32(0, 16)
	indices := math32.NewArrayU32(0, 16)

	// Generate plane vertices, vertices normals and vertices texture mappings.
	for iy := 0; iy < gridY1; iy++ {
		y := float32(iy)*segmentHeight - heightHalf
		for ix := 0; ix < gridX1; ix++ {
			x := float32(ix)*segmentWidth - widthHalf
			positions.Append(float32(x), float32(-y), 0)
			normals.Append(0, 0, 1)
			uvs.Append(float32(float64(ix)/float64(gridX)), float32(float64(1)-(float64(iy)/float64(gridY))))
		}
	}

	// Generate plane vertices indices for the faces
	for iy := 0; iy < gridY; iy++ {
		for ix := 0; ix < gridX; ix++ {
			a := ix + gridX1*iy
			b := ix + gridX1*(iy+1)
			c := (ix + 1) + gridX1*(iy+1)
			d := (ix + 1) + gridX1*iy
			indices.Append(uint32(a), uint32(b), uint32(d))
			indices.Append(uint32(b), uint32(c), uint32(d))
		}
	}

	plane.SetIndices(indices)
	plane.AddVBO(gls.NewVBO(positions).AddAttrib(gls.VertexPosition))
	plane.AddVBO(gls.NewVBO(normals).AddAttrib(gls.VertexNormal))
	plane.AddVBO(gls.NewVBO(uvs).AddAttrib(gls.VertexTexcoord))

	// Update area
	plane.area = width * height
	plane.areaValid = true

	// Update volume
	plane.volume = 0
	plane.volumeValid = true

	return plane
}
