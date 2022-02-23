package physics

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
	"github.com/bhojpur/render/pkg/3d/graphic"
	"github.com/bhojpur/render/pkg/3d/material"
	"github.com/bhojpur/render/pkg/3d/math32"
	"github.com/bhojpur/render/pkg/experimental/collision"
	"github.com/bhojpur/render/pkg/gls"
)

// This file contains helpful infrastructure for debugging physics
type DebugHelper struct {
}

func ShowWorldFace(scene *core.Node, face []math32.Vector3, color *math32.Color) {

	if len(face) == 0 {
		return
	}

	vertices := math32.NewArrayF32(0, 16)
	for i := range face {
		vertices.AppendVector3(&face[i])
	}
	vertices.AppendVector3(&face[0])

	geom := geometry.NewGeometry()
	geom.AddVBO(gls.NewVBO(vertices).AddAttrib(gls.VertexPosition))

	mat := material.NewStandard(color)
	faceGraphic := graphic.NewLineStrip(geom, mat)
	scene.Add(faceGraphic)
}

func ShowPenAxis(scene *core.Node, axis *math32.Vector3) { //}, min, max float32) {

	vertices := math32.NewArrayF32(0, 16)

	size := float32(100)
	minPoint := axis.Clone().MultiplyScalar(size)
	maxPoint := axis.Clone().MultiplyScalar(-size)
	//vertices.AppendVector3(minPoint.Clone().SetX(minPoint.X - size))
	//vertices.AppendVector3(minPoint.Clone().SetX(minPoint.X + size))
	//vertices.AppendVector3(minPoint.Clone().SetY(minPoint.Y - size))
	//vertices.AppendVector3(minPoint.Clone().SetY(minPoint.Y + size))
	//vertices.AppendVector3(minPoint.Clone().SetZ(minPoint.Z - size))
	//vertices.AppendVector3(minPoint.Clone().SetZ(minPoint.Z + size))
	vertices.AppendVector3(minPoint)
	//vertices.AppendVector3(maxPoint.Clone().SetX(maxPoint.X - size))
	//vertices.AppendVector3(maxPoint.Clone().SetX(maxPoint.X + size))
	//vertices.AppendVector3(maxPoint.Clone().SetY(maxPoint.Y - size))
	//vertices.AppendVector3(maxPoint.Clone().SetY(maxPoint.Y + size))
	//vertices.AppendVector3(maxPoint.Clone().SetZ(maxPoint.Z - size))
	//vertices.AppendVector3(maxPoint.Clone().SetZ(maxPoint.Z + size))
	vertices.AppendVector3(maxPoint)

	geom := geometry.NewGeometry()
	geom.AddVBO(gls.NewVBO(vertices).AddAttrib(gls.VertexPosition))

	mat := material.NewStandard(&math32.Color{1, 1, 1})
	faceGraphic := graphic.NewLines(geom, mat)
	scene.Add(faceGraphic)
}

func ShowContact(scene *core.Node, contact *collision.Contact) {

	vertices := math32.NewArrayF32(0, 16)

	size := float32(0.0005)
	otherPoint := contact.Point.Clone().Add(contact.Normal.Clone().MultiplyScalar(-contact.Depth))
	vertices.AppendVector3(contact.Point.Clone().SetX(contact.Point.X - size))
	vertices.AppendVector3(contact.Point.Clone().SetX(contact.Point.X + size))
	vertices.AppendVector3(contact.Point.Clone().SetY(contact.Point.Y - size))
	vertices.AppendVector3(contact.Point.Clone().SetY(contact.Point.Y + size))
	vertices.AppendVector3(contact.Point.Clone().SetZ(contact.Point.Z - size))
	vertices.AppendVector3(contact.Point.Clone().SetZ(contact.Point.Z + size))
	vertices.AppendVector3(contact.Point.Clone())
	vertices.AppendVector3(otherPoint)

	geom := geometry.NewGeometry()
	geom.AddVBO(gls.NewVBO(vertices).AddAttrib(gls.VertexPosition))

	mat := material.NewStandard(&math32.Color{0, 0, 1})
	faceGraphic := graphic.NewLines(geom, mat)
	scene.Add(faceGraphic)
}
