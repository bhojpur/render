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
	"fmt"
	"image"
	"image/color"

	d2d "github.com/bhojpur/render/pkg/g2d/draw"
	"github.com/golang/freetype/truetype"
)

var DefaultFontData = d2d.FontData{Name: "luxi", Family: d2d.FontFamilySans, Style: d2d.FontStyleNormal}

type StackGraphicContext struct {
	Current *ContextStack
}

type ContextStack struct {
	Tr          d2d.Matrix
	Path        *d2d.Path
	LineWidth   float64
	Dash        []float64
	DashOffset  float64
	StrokeColor color.Color
	FillColor   color.Color
	FillRule    d2d.FillRule
	Cap         d2d.LineCap
	Join        d2d.LineJoin
	FontSize    float64
	FontData    d2d.FontData

	Font *truetype.Font
	// fontSize and dpi are used to calculate scale. scale is the number of
	// 26.6 fixed point units in 1 em.
	Scale float64

	Previous *ContextStack
}

// GetFontName gets the current FontData with fontSize as a string
func (cs *ContextStack) GetFontName() string {
	fontData := cs.FontData
	return fmt.Sprintf("%s:%d:%d:%9.2f", fontData.Name, fontData.Family, fontData.Style, cs.FontSize)
}

/**
 * Create a new Graphic context from an image
 */
func NewStackGraphicContext() *StackGraphicContext {
	gc := &StackGraphicContext{}
	gc.Current = new(ContextStack)
	gc.Current.Tr = d2d.NewIdentityMatrix()
	gc.Current.Path = new(d2d.Path)
	gc.Current.LineWidth = 1.0
	gc.Current.StrokeColor = image.Black
	gc.Current.FillColor = image.White
	gc.Current.Cap = d2d.RoundCap
	gc.Current.FillRule = d2d.FillRuleEvenOdd
	gc.Current.Join = d2d.RoundJoin
	gc.Current.FontSize = 10
	gc.Current.FontData = DefaultFontData
	return gc
}

func (gc *StackGraphicContext) GetMatrixTransform() d2d.Matrix {
	return gc.Current.Tr
}

func (gc *StackGraphicContext) SetMatrixTransform(Tr d2d.Matrix) {
	gc.Current.Tr = Tr
}

func (gc *StackGraphicContext) ComposeMatrixTransform(Tr d2d.Matrix) {
	gc.Current.Tr.Compose(Tr)
}

func (gc *StackGraphicContext) Rotate(angle float64) {
	gc.Current.Tr.Rotate(angle)
}

func (gc *StackGraphicContext) Translate(tx, ty float64) {
	gc.Current.Tr.Translate(tx, ty)
}

func (gc *StackGraphicContext) Scale(sx, sy float64) {
	gc.Current.Tr.Scale(sx, sy)
}

func (gc *StackGraphicContext) SetStrokeColor(c color.Color) {
	gc.Current.StrokeColor = c
}

func (gc *StackGraphicContext) SetFillColor(c color.Color) {
	gc.Current.FillColor = c
}

func (gc *StackGraphicContext) SetFillRule(f d2d.FillRule) {
	gc.Current.FillRule = f
}

func (gc *StackGraphicContext) SetLineWidth(lineWidth float64) {
	gc.Current.LineWidth = lineWidth
}

func (gc *StackGraphicContext) SetLineCap(cap d2d.LineCap) {
	gc.Current.Cap = cap
}

func (gc *StackGraphicContext) SetLineJoin(join d2d.LineJoin) {
	gc.Current.Join = join
}

func (gc *StackGraphicContext) SetLineDash(dash []float64, dashOffset float64) {
	gc.Current.Dash = dash
	gc.Current.DashOffset = dashOffset
}

func (gc *StackGraphicContext) SetFontSize(fontSize float64) {
	gc.Current.FontSize = fontSize
}

func (gc *StackGraphicContext) GetFontSize() float64 {
	return gc.Current.FontSize
}

func (gc *StackGraphicContext) SetFontData(fontData d2d.FontData) {
	gc.Current.FontData = fontData
}

func (gc *StackGraphicContext) GetFontData() d2d.FontData {
	return gc.Current.FontData
}

func (gc *StackGraphicContext) BeginPath() {
	gc.Current.Path.Clear()
}

func (gc *StackGraphicContext) GetPath() d2d.Path {
	return *gc.Current.Path.Copy()
}

func (gc *StackGraphicContext) IsEmpty() bool {
	return gc.Current.Path.IsEmpty()
}

func (gc *StackGraphicContext) LastPoint() (float64, float64) {
	return gc.Current.Path.LastPoint()
}

func (gc *StackGraphicContext) MoveTo(x, y float64) {
	gc.Current.Path.MoveTo(x, y)
}

func (gc *StackGraphicContext) LineTo(x, y float64) {
	gc.Current.Path.LineTo(x, y)
}

func (gc *StackGraphicContext) QuadCurveTo(cx, cy, x, y float64) {
	gc.Current.Path.QuadCurveTo(cx, cy, x, y)
}

func (gc *StackGraphicContext) CubicCurveTo(cx1, cy1, cx2, cy2, x, y float64) {
	gc.Current.Path.CubicCurveTo(cx1, cy1, cx2, cy2, x, y)
}

func (gc *StackGraphicContext) ArcTo(cx, cy, rx, ry, startAngle, angle float64) {
	gc.Current.Path.ArcTo(cx, cy, rx, ry, startAngle, angle)
}

func (gc *StackGraphicContext) Close() {
	gc.Current.Path.Close()
}

func (gc *StackGraphicContext) Save() {
	context := new(ContextStack)
	context.FontSize = gc.Current.FontSize
	context.FontData = gc.Current.FontData
	context.LineWidth = gc.Current.LineWidth
	context.StrokeColor = gc.Current.StrokeColor
	context.FillColor = gc.Current.FillColor
	context.FillRule = gc.Current.FillRule
	context.Dash = gc.Current.Dash
	context.DashOffset = gc.Current.DashOffset
	context.Cap = gc.Current.Cap
	context.Join = gc.Current.Join
	context.Path = gc.Current.Path.Copy()
	context.Font = gc.Current.Font
	context.Scale = gc.Current.Scale
	copy(context.Tr[:], gc.Current.Tr[:])
	context.Previous = gc.Current
	gc.Current = context
}

func (gc *StackGraphicContext) Restore() {
	if gc.Current.Previous != nil {
		oldContext := gc.Current
		gc.Current = gc.Current.Previous
		oldContext.Previous = nil
	}
}

func (gc *StackGraphicContext) GetFontName() string {
	return gc.Current.GetFontName()
}
