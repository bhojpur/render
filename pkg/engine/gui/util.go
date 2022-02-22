package gui

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

// RectBounds specifies the size of the boundaries of a rectangle.
// It can represent the thickness of the borders, the margins, or the padding of a rectangle.
type RectBounds struct {
	Top    float32
	Right  float32
	Bottom float32
	Left   float32
}

// Set sets the values of the RectBounds.
func (bs *RectBounds) Set(top, right, bottom, left float32) {

	if top >= 0 {
		bs.Top = top
	}
	if right >= 0 {
		bs.Right = right
	}
	if bottom >= 0 {
		bs.Bottom = bottom
	}
	if left >= 0 {
		bs.Left = left
	}
}

// Rect represents a rectangle.
type Rect struct {
	X      float32
	Y      float32
	Width  float32
	Height float32
}

// Contains determines whether a 2D point is inside the Rect.
func (r *Rect) Contains(x, y float32) bool {

	if x < r.X || x > r.X+r.Width {
		return false
	}
	if y < r.Y || y > r.Y+r.Height {
		return false
	}
	return true
}
