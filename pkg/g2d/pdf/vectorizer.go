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

// Vectorizer defines the minimal interface for document.Bdf
// to be passed to a PathConvertor.
// It is also implemented by for example VertexMatrixTransform
type Vectorizer interface {
	// MoveTo creates a new subpath that start at the specified point
	MoveTo(x, y float64)
	// LineTo adds a line to the current subpath
	LineTo(x, y float64)
	// CurveTo adds a quadratic bezier curve to the current subpath
	CurveTo(cx, cy, x, y float64)
	// CurveTo adds a cubic bezier curve to the current subpath
	CurveBezierCubicTo(cx1, cy1, cx2, cy2, x, y float64)
	// ArcTo adds an arc to the current subpath
	ArcTo(x, y, rx, ry, degRotate, degStart, degEnd float64)
	// ClosePath closes the subpath
	ClosePath()
}
