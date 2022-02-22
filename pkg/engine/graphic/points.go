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
	"github.com/bhojpur/render/pkg/gls"
)

// Points represents a geometry containing only points
type Points struct {
	Graphic             // Embedded graphic
	uniMVPm gls.Uniform // Model view projection matrix uniform location cache
	uniMVm  gls.Uniform // Model view matrix uniform location cache
}

// NewPoints creates and returns a graphic points object with the specified
// geometry and material.
func NewPoints(igeom geometry.IGeometry, imat material.IMaterial) *Points {

	p := new(Points)
	p.Graphic.Init(p, igeom, gls.POINTS)
	if imat != nil {
		p.AddMaterial(p, imat, 0, 0)
	}
	p.uniMVPm.Init("MVP")
	p.uniMVm.Init("MV")
	return p
}

// RenderSetup is called by the engine before rendering this graphic.
func (p *Points) RenderSetup(gs *gls.GLS, rinfo *core.RenderInfo) {

	// Transfer model view projection matrix uniform
	mvpm := p.ModelViewProjectionMatrix()
	location := p.uniMVPm.Location(gs)
	gs.UniformMatrix4fv(location, 1, false, &mvpm[0])

	// Transfer model view matrix uniform
	mvm := p.ModelViewMatrix()
	location = p.uniMVm.Location(gs)
	gs.UniformMatrix4fv(location, 1, false, &mvm[0])
}
