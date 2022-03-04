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
	"bytes"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"
	"strconv"

	"github.com/golang/freetype/truetype"

	"github.com/bhojpur/render/pkg/document"
	base "github.com/bhojpur/render/pkg/g2d/base"
	d2d "github.com/bhojpur/render/pkg/g2d/draw"
	kit "github.com/bhojpur/render/pkg/g2d/kit"
)

const (
	// DPI of a pdf document is fixed at 72.
	DPI  = 72
	c255 = 255.0 / 65535.0
)

var (
	caps = map[d2d.LineCap]string{
		d2d.RoundCap:  "round",
		d2d.ButtCap:   "butt",
		d2d.SquareCap: "square"}
	joins = map[d2d.LineJoin]string{
		d2d.RoundJoin: "round",
		d2d.BevelJoin: "bevel",
		d2d.MiterJoin: "miter",
	}
	imageCount uint32
	white      color.Color = color.RGBA{255, 255, 255, 255}
)

// NewPdf creates a new pdf document with the d2d fontfolder, adds
// a page and set fill color to white.
func NewPdf(orientationStr, unitStr, sizeStr string) *document.Bdf {
	pdf := document.New(orientationStr, unitStr, sizeStr, d2d.GetFontFolder())
	// to be compatible with d2d
	pdf.SetMargins(0, 0, 0)
	pdf.SetDrawColor(0, 0, 0)
	pdf.SetFillColor(255, 255, 255)
	pdf.SetLineCapStyle("round")
	pdf.SetLineJoinStyle("round")
	pdf.SetLineWidth(1)
	pdf.AddPage()
	return pdf
}

// rgb converts a color (used by d2d) into 3 int (used by Bdf)
func rgb(c color.Color) (int, int, int) {
	r, g, b, _ := c.RGBA()
	return int(float64(r) * c255), int(float64(g) * c255), int(float64(b) * c255)
}

// clearRect draws a white rectangle
func clearRect(gc *GraphicContext, x1, y1, x2, y2 float64) {
	// save state
	f := gc.Current.FillColor
	x, y := gc.pdf.GetXY()
	// cover page with white rectangle
	gc.SetFillColor(white)
	kit.Rectangle(gc, x1, y1, x2, y2)
	gc.Fill()
	// restore state
	gc.SetFillColor(f)
	gc.pdf.MoveTo(x, y)
}

// GraphicContext implements the d2d.GraphicContext interface
// It provides d2d with a pdf backend (based on Bdf)
type GraphicContext struct {
	*base.StackGraphicContext
	pdf *document.Bdf
	DPI int
}

// NewGraphicContext creates a new pdf GraphicContext
func NewGraphicContext(pdf *document.Bdf) *GraphicContext {
	gc := &GraphicContext{base.NewStackGraphicContext(), pdf, DPI}
	gc.SetDPI(DPI)
	return gc
}

// DrawImage draws an image as PNG
// TODO: add type (tp) as parameter to argument list?
func (gc *GraphicContext) DrawImage(image image.Image) {
	name := strconv.Itoa(int(imageCount))
	imageCount++
	tp := "PNG" // "JPG", "JPEG", "PNG" and "GIF"
	b := &bytes.Buffer{}
	png.Encode(b, image)
	gc.pdf.RegisterImageReader(name, tp, b)
	bounds := image.Bounds()
	x0, y0 := float64(bounds.Min.X), float64(bounds.Min.Y)
	w, h := float64(bounds.Dx()), float64(bounds.Dy())
	gc.pdf.Image(name, x0, y0, w, h, false, tp, 0, "")
}

// Clear draws a white rectangle over the whole page
func (gc *GraphicContext) Clear() {
	width, height := gc.pdf.GetPageSize()
	clearRect(gc, 0, 0, width, height)
}

// ClearRect draws a white rectangle over the specified area.
// Samples: line.
func (gc *GraphicContext) ClearRect(x1, y1, x2, y2 int) {
	clearRect(gc, float64(x1), float64(y1), float64(x2), float64(y2))
}

// recalc recalculates scale and bounds values from the font size, screen
// resolution and font metrics, and invalidates the glyph cache.
func (gc *GraphicContext) recalc() {
	// TODO: resolve properly the font size for pdf and bitmap
	gc.Current.Scale = 3 * float64(gc.DPI) / 72
}

// SetDPI sets the DPI which influences the font size.
func (gc *GraphicContext) SetDPI(dpi int) {
	gc.DPI = dpi
	gc.recalc()
}

// GetDPI returns the DPI which influences the font size.
// (Note that Bdf uses a fixed dpi of 72:
func (gc *GraphicContext) GetDPI() int {
	return gc.DPI
}

// GetStringBounds returns the approximate pixel bounds of the string s at x, y.
// The left edge of the em square of the first character of s
// and the baseline intersect at 0, 0 in the returned coordinates.
// Therefore the top and left coordinates may well be negative.
func (gc *GraphicContext) GetStringBounds(s string) (left, top, right, bottom float64) {
	_, h := gc.pdf.GetFontSize()
	d := gc.pdf.GetFontDesc("", "")
	if d.Ascent == 0 {
		// not defined (standard font?), use average of 81%
		top = 0.81 * h
	} else {
		top = -float64(d.Ascent) * h / float64(d.Ascent-d.Descent)
	}
	return 0, top, gc.pdf.GetStringWidth(s), top + h
}

// CreateStringPath creates a path from the string s at x, y, and returns the string width.
func (gc *GraphicContext) CreateStringPath(text string, x, y float64) (cursor float64) {
	//Bdf uses the top left corner
	left, top, right, bottom := gc.GetStringBounds(text)
	w := right - left
	h := bottom - top
	// gc.pdf.SetXY(x, y-h) do not use this as y-h might be negative
	margin := gc.pdf.GetCellMargin()
	gc.pdf.MoveTo(x-left-margin, y+top)
	gc.pdf.CellFormat(w, h, text, "", 0, "BL", false, 0, "")
	return w
}

// FillString draws a string at 0, 0
func (gc *GraphicContext) FillString(text string) (cursor float64) {
	return gc.FillStringAt(text, 0, 0)
}

// FillStringAt draws a string at x, y
func (gc *GraphicContext) FillStringAt(text string, x, y float64) (cursor float64) {
	return gc.CreateStringPath(text, x, y)
}

// StrokeString draws a string at 0, 0 (stroking is unsupported,
// string will be filled)
func (gc *GraphicContext) StrokeString(text string) (cursor float64) {
	return gc.StrokeStringAt(text, 0, 0)
}

// StrokeStringAt draws a string at x, y (stroking is unsupported,
// string will be filled)
func (gc *GraphicContext) StrokeStringAt(text string, x, y float64) (cursor float64) {
	return gc.CreateStringPath(text, x, y)
}

// Stroke strokes the paths with the color specified by SetStrokeColor
func (gc *GraphicContext) Stroke(paths ...*d2d.Path) {
	_, _, _, alphaS := gc.Current.StrokeColor.RGBA()
	gc.draw("D", alphaS, paths...)
	gc.Current.Path.Clear()
}

// Fill fills the paths with the color specified by SetFillColor
func (gc *GraphicContext) Fill(paths ...*d2d.Path) {
	style := "F"
	if gc.Current.FillRule != d2d.FillRuleWinding {
		style += "*"
	}
	_, _, _, alphaF := gc.Current.FillColor.RGBA()
	gc.draw(style, alphaF, paths...)
	gc.Current.Path.Clear()
}

// FillStroke first fills the paths and than strokes them
func (gc *GraphicContext) FillStroke(paths ...*d2d.Path) {
	var rule string
	if gc.Current.FillRule != d2d.FillRuleWinding {
		rule = "*"
	}
	_, _, _, alphaS := gc.Current.StrokeColor.RGBA()
	_, _, _, alphaF := gc.Current.FillColor.RGBA()
	if alphaS == alphaF {
		gc.draw("FD"+rule, alphaF, paths...)
	} else {
		gc.draw("F"+rule, alphaF, paths...)
		gc.draw("S", alphaS, paths...)
	}
	gc.Current.Path.Clear()
}

var logger = log.New(os.Stdout, "", log.Lshortfile)

const alphaMax = float64(0xFFFF)

// draw fills and/or strokes paths
func (gc *GraphicContext) draw(style string, alpha uint32, paths ...*d2d.Path) {
	paths = append(paths, gc.Current.Path)
	for _, p := range paths {
		ConvertPath(p, gc.pdf)
	}
	a := float64(alpha) / alphaMax
	current, blendMode := gc.pdf.GetAlpha()
	if a != current {
		gc.pdf.SetAlpha(a, blendMode)
	}
	gc.pdf.DrawPath(style)
}

// overwrite StackGraphicContext methods

// SetStrokeColor sets the stroke color
func (gc *GraphicContext) SetStrokeColor(c color.Color) {
	gc.StackGraphicContext.SetStrokeColor(c)
	gc.pdf.SetDrawColor(rgb(c))
}

// SetFillColor sets the fill and text color
func (gc *GraphicContext) SetFillColor(c color.Color) {
	gc.StackGraphicContext.SetFillColor(c)
	gc.pdf.SetFillColor(rgb(c))
	gc.pdf.SetTextColor(rgb(c))
}

// SetFont is unsupported by the pdf graphic context, use SetFontData
// instead.
func (gc *GraphicContext) SetFont(font *truetype.Font) {
	// TODO: what to do with this api conflict between d2d and Bdf?!
}

// SetFontData sets the current font used to draw text. Always use
// this method, as SetFont is unsupported by the pdf graphic context.
// It is mandatory to call this method at least once before printing
// text or the resulting document will not be valid.
// It is necessary to generate a font definition file first with the
// makefont utility. It is not necessary to call this function for the
// core PDF fonts (courier, helvetica, times, zapfdingbats).
func (gc *GraphicContext) SetFontData(fontData d2d.FontData) {
	// TODO: call Makefont embed if json file does not exist yet
	gc.StackGraphicContext.SetFontData(fontData)
	var style string
	if fontData.Style&d2d.FontStyleBold != 0 {
		style += "B"
	}
	if fontData.Style&d2d.FontStyleItalic != 0 {
		style += "I"
	}
	fn := d2d.FontFileName(fontData)
	fn = fn[:len(fn)-4]
	size, _ := gc.pdf.GetFontSize()
	gc.pdf.AddFont(fontData.Name, style, fn+".json")
	gc.pdf.SetFont(fontData.Name, style, size)
}

// SetFontSize sets the font size in points (as in ``a 12 point font'').
// TODO: resolve this with ImgGraphicContext (now done with gc.Current.Scale)
func (gc *GraphicContext) SetFontSize(fontSize float64) {
	gc.StackGraphicContext.SetFontSize(fontSize)
	gc.recalc()
	gc.pdf.SetFontSize(fontSize * gc.Current.Scale)
}

// SetLineDash sets the line dash pattern
func (gc *GraphicContext) SetLineDash(Dash []float64, DashOffset float64) {
	gc.StackGraphicContext.SetLineDash(Dash, DashOffset)
	gc.pdf.SetDashPattern(Dash, DashOffset)
}

// SetLineWidth sets the line width
func (gc *GraphicContext) SetLineWidth(LineWidth float64) {
	gc.StackGraphicContext.SetLineWidth(LineWidth)
	gc.pdf.SetLineWidth(LineWidth)
}

// SetLineCap sets the line cap (round, but or square)
func (gc *GraphicContext) SetLineCap(Cap d2d.LineCap) {
	gc.StackGraphicContext.SetLineCap(Cap)
	gc.pdf.SetLineCapStyle(caps[Cap])
}

// SetLineJoin sets the line cap (round, bevel or miter)
func (gc *GraphicContext) SetLineJoin(Join d2d.LineJoin) {
	gc.StackGraphicContext.SetLineJoin(Join)
	gc.pdf.SetLineJoinStyle(joins[Join])
}

// Transformations

// Scale generally scales the following text, drawings and images.
// sx and sy are the scaling factors for width and height.
// This must be placed between gc.Save() and gc.Restore(), otherwise
// the pdf is invalid.
func (gc *GraphicContext) Scale(sx, sy float64) {
	gc.StackGraphicContext.Scale(sx, sy)
	gc.pdf.TransformScale(sx*100, sy*100, 0, 0)
}

// Rotate rotates the following text, drawings and images.
// Angle is specified in radians and measured clockwise from the
// 3 o'clock position.
// This must be placed between gc.Save() and gc.Restore(), otherwise
// the pdf is invalid.
func (gc *GraphicContext) Rotate(angle float64) {
	gc.StackGraphicContext.Rotate(angle)
	gc.pdf.TransformRotate(-angle*180/math.Pi, 0, 0)
}

// Translate moves the following text, drawings and images
// horizontally and vertically by the amounts specified by tx and ty.
// This must be placed between gc.Save() and gc.Restore(), otherwise
// the pdf is invalid.
func (gc *GraphicContext) Translate(tx, ty float64) {
	gc.StackGraphicContext.Translate(tx, ty)
	gc.pdf.TransformTranslate(tx, ty)
}

// Save saves the current context stack
// (transformation, font, color,...).
func (gc *GraphicContext) Save() {
	gc.StackGraphicContext.Save()
	gc.pdf.TransformBegin()
}

// Restore restores the current context stack
// (transformation, color,...). Restoring the font is not supported.
func (gc *GraphicContext) Restore() {
	gc.pdf.TransformEnd()
	gc.StackGraphicContext.Restore()
	c := gc.Current
	gc.SetFontSize(c.FontSize)
	// gc.SetFontData(c.FontData) unsupported, causes bug (do not enable)
	gc.SetLineWidth(c.LineWidth)
	gc.SetStrokeColor(c.StrokeColor)
	gc.SetFillColor(c.FillColor)
	gc.SetFillRule(c.FillRule)
	// gc.SetLineDash(c.Dash, c.DashOffset) // TODO
	gc.SetLineCap(c.Cap)
	gc.SetLineJoin(c.Join)
	// c.Path unsupported
	// c.Font unsupported
}
