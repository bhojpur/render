package pdf

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

	d2d "github.com/bhojpur/render/pkg/g2d/draw"
)

const deg = 180 / math.Pi

// ConvertPath converts a paths to the pdf api
func ConvertPath(path *d2d.Path, pdf Vectorizer) {
	var startX, startY float64 = 0, 0
	i := 0
	for _, cmp := range path.Components {
		switch cmp {
		case d2d.MoveToCmp:
			startX, startY = path.Points[i], path.Points[i+1]
			pdf.MoveTo(startX, startY)
			i += 2
		case d2d.LineToCmp:
			pdf.LineTo(path.Points[i], path.Points[i+1])
			i += 2
		case d2d.QuadCurveToCmp:
			pdf.CurveTo(path.Points[i], path.Points[i+1], path.Points[i+2], path.Points[i+3])
			i += 4
		case d2d.CubicCurveToCmp:
			pdf.CurveBezierCubicTo(path.Points[i], path.Points[i+1], path.Points[i+2], path.Points[i+3], path.Points[i+4], path.Points[i+5])
			i += 6
		case d2d.ArcToCmp:
			pdf.ArcTo(path.Points[i], path.Points[i+1], path.Points[i+2], path.Points[i+3],
				0,                    // degRotate
				path.Points[i+4]*deg, // degStart = startAngle
				(path.Points[i+4]-path.Points[i+5])*deg) // degEnd = startAngle-angle
			i += 6
		case d2d.CloseCmp:
			pdf.LineTo(startX, startY)
			pdf.ClosePath()
		}
	}
}
