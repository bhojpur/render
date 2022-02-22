package graphic

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
	"github.com/bhojpur/render/pkg/engine/core"
	"github.com/bhojpur/render/pkg/engine/geometry"
	"github.com/bhojpur/render/pkg/engine/material"
	"github.com/bhojpur/render/pkg/engine/math32"
	"github.com/bhojpur/render/pkg/gls"
)

// Sprite is a potentially animated image positioned in space that always faces the camera.
type Sprite struct {
	Graphic             // Embedded graphic
	uniMVPM gls.Uniform // Model view projection matrix uniform location cache
}

// NewSprite creates and returns a pointer to a sprite with the specified dimensions and material
func NewSprite(width, height float32, imat material.IMaterial) *Sprite {

	s := new(Sprite)

	// Creates geometry
	geom := geometry.NewGeometry()
	w := width / 2
	h := height / 2

	// Builds array with vertex positions and texture coordinates
	positions := math32.NewArrayF32(0, 12)
	positions.Append(
		-w, -h, 0, 0, 0,
		w, -h, 0, 1, 0,
		w, h, 0, 1, 1,
		-w, h, 0, 0, 1,
	)
	// Builds array of indices
	indices := math32.NewArrayU32(0, 6)
	indices.Append(0, 1, 2, 0, 2, 3)

	// Set geometry buffers
	geom.SetIndices(indices)
	geom.AddVBO(
		gls.NewVBO(positions).
			AddAttrib(gls.VertexPosition).
			AddAttrib(gls.VertexTexcoord),
	)

	s.Graphic.Init(s, geom, gls.TRIANGLES)
	s.AddMaterial(s, imat, 0, 0)

	s.uniMVPM.Init("MVP")
	return s
}

// RenderSetup sets up the rendering of the sprite.
func (s *Sprite) RenderSetup(gs *gls.GLS, rinfo *core.RenderInfo) {

	// Calculates model view matrix
	mw := s.MatrixWorld()
	var mvm math32.Matrix4
	mvm.MultiplyMatrices(&rinfo.ViewMatrix, &mw)

	// Decomposes model view matrix
	var position math32.Vector3
	var quaternion math32.Quaternion
	var scale math32.Vector3
	mvm.Decompose(&position, &quaternion, &scale)

	// Removes any rotation in X and Y axes and compose new model view matrix
	rotation := s.Rotation()
	actualScale := s.Scale()
	if actualScale.X >= 0 {
		rotation.Y = 0
	} else {
		rotation.Y = math32.Pi
	}
	if actualScale.Y >= 0 {
		rotation.X = 0
	} else {
		rotation.X = math32.Pi
	}
	quaternion.SetFromEuler(&rotation)
	var mvmNew math32.Matrix4
	mvmNew.Compose(&position, &quaternion, &scale)

	// Calculates final MVP and updates uniform
	var mvpm math32.Matrix4
	mvpm.MultiplyMatrices(&rinfo.ProjMatrix, &mvmNew)
	location := s.uniMVPM.Location(gs)
	gs.UniformMatrix4fv(location, 1, false, &mvpm[0])
}
