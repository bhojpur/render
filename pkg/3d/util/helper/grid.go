package helper

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
	"github.com/bhojpur/render/pkg/3d/geometry"
	"github.com/bhojpur/render/pkg/3d/graphic"
	"github.com/bhojpur/render/pkg/3d/material"
	"github.com/bhojpur/render/pkg/3d/math32"
	"github.com/bhojpur/render/pkg/gls"
)

// Grid is a visual representation of a grid.
type Grid struct {
	graphic.Lines
}

// NewGrid creates and returns a pointer to a new grid helper with the specified size and step.
func NewGrid(size, step float32, color *math32.Color) *Grid {

	grid := new(Grid)

	half := size / 2
	positions := math32.NewArrayF32(0, 0)
	for i := -half; i <= half; i += step {
		positions.Append(
			-half, 0, i, color.R, color.G, color.B,
			half, 0, i, color.R, color.G, color.B,
			i, 0, -half, color.R, color.G, color.B,
			i, 0, half, color.R, color.G, color.B,
		)
	}

	// Create geometry
	geom := geometry.NewGeometry()
	geom.AddVBO(
		gls.NewVBO(positions).
			AddAttrib(gls.VertexPosition).
			AddAttrib(gls.VertexColor),
	)

	// Create material
	mat := material.NewBasic()

	// Initialize lines with the specified geometry and material
	grid.Lines.Init(geom, mat)
	return grid
}
