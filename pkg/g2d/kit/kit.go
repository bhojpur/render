package kit

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

// It provides helpers to draw common figures using a Path

import (
	"math"

	d2d "github.com/bhojpur/render/pkg/g2d/draw"
)

// Rectangle draws a rectangle using a path between (x1,y1) and (x2,y2)
func Rectangle(path d2d.PathBuilder, x1, y1, x2, y2 float64) {
	path.MoveTo(x1, y1)
	path.LineTo(x2, y1)
	path.LineTo(x2, y2)
	path.LineTo(x1, y2)
	path.Close()
}

// RoundedRectangle draws a rectangle using a path between (x1,y1) and (x2,y2)
func RoundedRectangle(path d2d.PathBuilder, x1, y1, x2, y2, arcWidth, arcHeight float64) {
	arcWidth = arcWidth / 2
	arcHeight = arcHeight / 2
	path.MoveTo(x1, y1+arcHeight)
	path.QuadCurveTo(x1, y1, x1+arcWidth, y1)
	path.LineTo(x2-arcWidth, y1)
	path.QuadCurveTo(x2, y1, x2, y1+arcHeight)
	path.LineTo(x2, y2-arcHeight)
	path.QuadCurveTo(x2, y2, x2-arcWidth, y2)
	path.LineTo(x1+arcWidth, y2)
	path.QuadCurveTo(x1, y2, x1, y2-arcHeight)
	path.Close()
}

// Ellipse draws an ellipse using a path with center (cx,cy) and radius (rx,ry)
func Ellipse(path d2d.PathBuilder, cx, cy, rx, ry float64) {
	path.ArcTo(cx, cy, rx, ry, 0, -math.Pi*2)
	path.Close()
}

// Circle draws a circle using a path with center (cx,cy) and radius
func Circle(path d2d.PathBuilder, cx, cy, radius float64) {
	path.ArcTo(cx, cy, radius, radius, 0, -math.Pi*2)
	path.Close()
}
