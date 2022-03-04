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
	"github.com/bhojpur/render/pkg/g3d/core"
	"github.com/bhojpur/render/pkg/g3d/geometry"
	"github.com/bhojpur/render/pkg/g3d/material"
	"github.com/bhojpur/render/pkg/gls"
)

// Lines is a Graphic which is rendered as a collection of independent lines.
type Lines struct {
	Graphic             // Embedded graphic object
	uniMVPm gls.Uniform // Model view projection matrix uniform location cache
}

// NewLines returns a pointer to a new Lines object.
func NewLines(igeom geometry.IGeometry, imat material.IMaterial) *Lines {

	l := new(Lines)
	l.Init(igeom, imat)
	return l
}

// Init initializes the Lines object and adds the specified material.
func (l *Lines) Init(igeom geometry.IGeometry, imat material.IMaterial) {

	l.Graphic.Init(l, igeom, gls.LINES)
	l.AddMaterial(l, imat, 0, 0)
	l.uniMVPm.Init("MVP")
}

// RenderSetup is called by the engine before drawing this geometry.
func (l *Lines) RenderSetup(gs *gls.GLS, rinfo *core.RenderInfo) {

	// Transfer model view projection matrix uniform
	mvpm := l.ModelViewProjectionMatrix()
	location := l.uniMVPm.Location(gs)
	gs.UniformMatrix4fv(location, 1, false, &mvpm[0])
}
