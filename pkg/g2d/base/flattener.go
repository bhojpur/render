package base

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
	d2d "github.com/bhojpur/render/pkg/g2d/draw"
)

// Liner receive segment definition
type Liner interface {
	// LineTo Draw a line from the current position to the point (x, y)
	LineTo(x, y float64)
}

// Flattener receive segment definition
type Flattener interface {
	// MoveTo Start a New line from the point (x, y)
	MoveTo(x, y float64)
	// LineTo Draw a line from the current position to the point (x, y)
	LineTo(x, y float64)
	// LineJoin use Round, Bevel or miter to join points
	LineJoin()
	// Close add the most recent starting point to close the path to create a polygon
	Close()
	// End mark the current line as finished so we can draw caps
	End()
}

// Flatten convert curves into straight segments keeping join segments info
func Flatten(path *d2d.Path, flattener Flattener, scale float64) {
	// First Point
	var startX, startY float64 = 0, 0
	// Current Point
	var x, y float64 = 0, 0
	i := 0
	for _, cmp := range path.Components {
		switch cmp {
		case d2d.MoveToCmp:
			x, y = path.Points[i], path.Points[i+1]
			startX, startY = x, y
			if i != 0 {
				flattener.End()
			}
			flattener.MoveTo(x, y)
			i += 2
		case d2d.LineToCmp:
			x, y = path.Points[i], path.Points[i+1]
			flattener.LineTo(x, y)
			flattener.LineJoin()
			i += 2
		case d2d.QuadCurveToCmp:
			TraceQuad(flattener, path.Points[i-2:], 0.5)
			x, y = path.Points[i+2], path.Points[i+3]
			flattener.LineTo(x, y)
			i += 4
		case d2d.CubicCurveToCmp:
			TraceCubic(flattener, path.Points[i-2:], 0.5)
			x, y = path.Points[i+4], path.Points[i+5]
			flattener.LineTo(x, y)
			i += 6
		case d2d.ArcToCmp:
			x, y = TraceArc(flattener, path.Points[i], path.Points[i+1], path.Points[i+2], path.Points[i+3], path.Points[i+4], path.Points[i+5], scale)
			flattener.LineTo(x, y)
			i += 6
		case d2d.CloseCmp:
			flattener.LineTo(startX, startY)
			flattener.Close()
		}
	}
	flattener.End()
}

// Transformer apply the Matrix transformation tr
type Transformer struct {
	Tr        d2d.Matrix
	Flattener Flattener
}

func (t Transformer) MoveTo(x, y float64) {
	u := x*t.Tr[0] + y*t.Tr[2] + t.Tr[4]
	v := x*t.Tr[1] + y*t.Tr[3] + t.Tr[5]
	t.Flattener.MoveTo(u, v)
}

func (t Transformer) LineTo(x, y float64) {
	u := x*t.Tr[0] + y*t.Tr[2] + t.Tr[4]
	v := x*t.Tr[1] + y*t.Tr[3] + t.Tr[5]
	t.Flattener.LineTo(u, v)
}

func (t Transformer) LineJoin() {
	t.Flattener.LineJoin()
}

func (t Transformer) Close() {
	t.Flattener.Close()
}

func (t Transformer) End() {
	t.Flattener.End()
}

type SegmentedPath struct {
	Points []float64
}

func (p *SegmentedPath) MoveTo(x, y float64) {
	p.Points = append(p.Points, x, y)
	// TODO need to mark this point as moveto
}

func (p *SegmentedPath) LineTo(x, y float64) {
	p.Points = append(p.Points, x, y)
}

func (p *SegmentedPath) LineJoin() {
	// TODO need to mark the current point as linejoin
}

func (p *SegmentedPath) Close() {
	// TODO Close
}

func (p *SegmentedPath) End() {
	// Nothing to do
}
