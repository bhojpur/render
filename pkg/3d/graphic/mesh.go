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
	"github.com/bhojpur/render/pkg/3d/core"
	"github.com/bhojpur/render/pkg/3d/geometry"
	"github.com/bhojpur/render/pkg/3d/material"
	"github.com/bhojpur/render/pkg/3d/math32"
	"github.com/bhojpur/render/pkg/gls"
)

// Mesh is a Graphic with uniforms for the model, view, projection, and normal matrices.
type Mesh struct {
	Graphic             // Embedded graphic
	uniMm   gls.Uniform // Model matrix uniform location cache
	uniMVm  gls.Uniform // Model view matrix uniform location cache
	uniMVPm gls.Uniform // Model view projection matrix uniform cache
	uniNm   gls.Uniform // Normal matrix uniform cache
}

// NewMesh creates and returns a pointer to a mesh with the specified geometry and material.
// If the mesh has multi materials, the material specified here must be nil and the
// individual materials must be add using "AddMaterial" or AddGroupMaterial".
func NewMesh(igeom geometry.IGeometry, imat material.IMaterial) *Mesh {

	m := new(Mesh)
	m.Init(igeom, imat)
	return m
}

// Init initializes the Mesh and its uniforms.
func (m *Mesh) Init(igeom geometry.IGeometry, imat material.IMaterial) {

	m.Graphic.Init(m, igeom, gls.TRIANGLES)

	// Initialize uniforms
	m.uniMm.Init("ModelMatrix")
	m.uniMVm.Init("ModelViewMatrix")
	m.uniMVPm.Init("MVP")
	m.uniNm.Init("NormalMatrix")

	// Adds single material if not nil
	if imat != nil {
		m.AddMaterial(imat, 0, 0)
	}
}

// SetMaterial clears all materials and adds the specified material for all vertices.
func (m *Mesh) SetMaterial(imat material.IMaterial) {

	m.Graphic.ClearMaterials()
	m.Graphic.AddMaterial(m, imat, 0, 0)
}

// AddMaterial adds a material for the specified subset of vertices.
func (m *Mesh) AddMaterial(imat material.IMaterial, start, count int) {

	m.Graphic.AddMaterial(m, imat, start, count)
}

// AddGroupMaterial adds a material for the specified geometry group.
func (m *Mesh) AddGroupMaterial(imat material.IMaterial, gindex int) {

	m.Graphic.AddGroupMaterial(m, imat, gindex)
}

// Clone clones the mesh and satisfies the INode interface.
func (m *Mesh) Clone() core.INode {

	clone := new(Mesh)
	clone.Graphic = *m.Graphic.Clone().(*Graphic)
	clone.SetIGraphic(clone)

	// Initialize uniforms
	clone.uniMm.Init("ModelMatrix")
	clone.uniMVm.Init("ModelViewMatrix")
	clone.uniMVPm.Init("MVP")
	clone.uniNm.Init("NormalMatrix")

	return clone
}

// RenderSetup is called by the engine before drawing the mesh geometry
// It is responsible to updating the current shader uniforms with
// the model matrices.
func (m *Mesh) RenderSetup(gs *gls.GLS, rinfo *core.RenderInfo) {

	// Transfer uniform for model matrix
	mm := m.ModelMatrix()
	location := m.uniMm.Location(gs)
	gs.UniformMatrix4fv(location, 1, false, &mm[0])

	// Transfer uniform for model view matrix
	mvm := m.ModelViewMatrix()
	location = m.uniMVm.Location(gs)
	gs.UniformMatrix4fv(location, 1, false, &mvm[0])

	// Transfer uniform for model view projection matrix
	mvpm := m.ModelViewProjectionMatrix()
	location = m.uniMVPm.Location(gs)
	gs.UniformMatrix4fv(location, 1, false, &mvpm[0])

	// Calculates normal matrix and transfer uniform
	var nm math32.Matrix3
	nm.GetNormalMatrix(mvm)
	location = m.uniNm.Location(gs)
	gs.UniformMatrix3fv(location, 1, false, &nm[0])
}
