package draw

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
	"image"
	"image/color"
)

// GraphicContext describes the interface for the various backends (images, pdf, opengl, ...)
type GraphicContext interface {
	// PathBuilder describes the interface for path drawing
	PathBuilder
	// BeginPath creates a new path
	BeginPath()
	// GetPath copies the current path, then returns it
	GetPath() Path
	// GetMatrixTransform returns the current transformation matrix
	GetMatrixTransform() Matrix
	// SetMatrixTransform sets the current transformation matrix
	SetMatrixTransform(tr Matrix)
	// ComposeMatrixTransform composes the current transformation matrix with tr
	ComposeMatrixTransform(tr Matrix)
	// Rotate applies a rotation to the current transformation matrix. angle is in radian.
	Rotate(angle float64)
	// Translate applies a translation to the current transformation matrix.
	Translate(tx, ty float64)
	// Scale applies a scale to the current transformation matrix.
	Scale(sx, sy float64)
	// SetStrokeColor sets the current stroke color
	SetStrokeColor(c color.Color)
	// SetFillColor sets the current fill color
	SetFillColor(c color.Color)
	// SetFillRule sets the current fill rule
	SetFillRule(f FillRule)
	// SetLineWidth sets the current line width
	SetLineWidth(lineWidth float64)
	// SetLineCap sets the current line cap
	SetLineCap(cap LineCap)
	// SetLineJoin sets the current line join
	SetLineJoin(join LineJoin)
	// SetLineDash sets the current dash
	SetLineDash(dash []float64, dashOffset float64)
	// SetFontSize sets the current font size
	SetFontSize(fontSize float64)
	// GetFontSize gets the current font size
	GetFontSize() float64
	// SetFontData sets the current FontData
	SetFontData(fontData FontData)
	// GetFontData gets the current FontData
	GetFontData() FontData
	// GetFontName gets the current FontData as a string
	GetFontName() string
	// DrawImage draws the raster image in the current canvas
	DrawImage(image image.Image)
	// Save the context and push it to the context stack
	Save()
	// Restore remove the current context and restore the last one
	Restore()
	// Clear fills the current canvas with a default transparent color
	Clear()
	// ClearRect fills the specified rectangle with a default transparent color
	ClearRect(x1, y1, x2, y2 int)
	// SetDPI sets the current DPI
	SetDPI(dpi int)
	// GetDPI gets the current DPI
	GetDPI() int
	// GetStringBounds gets pixel bounds(dimensions) of given string
	GetStringBounds(s string) (left, top, right, bottom float64)
	// CreateStringPath creates a path from the string s at x, y
	CreateStringPath(text string, x, y float64) (cursor float64)
	// FillString draws the text at point (0, 0)
	FillString(text string) (cursor float64)
	// FillStringAt draws the text at the specified point (x, y)
	FillStringAt(text string, x, y float64) (cursor float64)
	// StrokeString draws the contour of the text at point (0, 0)
	StrokeString(text string) (cursor float64)
	// StrokeStringAt draws the contour of the text at point (x, y)
	StrokeStringAt(text string, x, y float64) (cursor float64)
	// Stroke strokes the paths with the color specified by SetStrokeColor
	Stroke(paths ...*Path)
	// Fill fills the paths with the color specified by SetFillColor
	Fill(paths ...*Path)
	// FillStroke first fills the paths and than strokes them
	FillStroke(paths ...*Path)
}
