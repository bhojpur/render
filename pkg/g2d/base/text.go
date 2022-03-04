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

import d2d "github.com/bhojpur/render/pkg/g2d/draw"

// GlyphCache manage a cache of glyphs
type GlyphCache interface {
	// Fetch fetches a glyph from the cache, storing with Render first if it doesn't already exist
	Fetch(gc d2d.GraphicContext, fontName string, chr rune) *Glyph
}

// GlyphCacheImp manage a map of glyphs without sync mecanism, not thread safe
type GlyphCacheImp struct {
	glyphs map[string]map[rune]*Glyph
}

// NewGlyphCache initializes a GlyphCache
func NewGlyphCache() *GlyphCacheImp {
	glyphs := make(map[string]map[rune]*Glyph)
	return &GlyphCacheImp{
		glyphs: glyphs,
	}
}

// Fetch fetches a glyph from the cache, calling renderGlyph first if it doesn't already exist
func (glyphCache *GlyphCacheImp) Fetch(gc d2d.GraphicContext, fontName string, chr rune) *Glyph {
	if glyphCache.glyphs[fontName] == nil {
		glyphCache.glyphs[fontName] = make(map[rune]*Glyph, 60)
	}
	if glyphCache.glyphs[fontName][chr] == nil {
		glyphCache.glyphs[fontName][chr] = renderGlyph(gc, fontName, chr)
	}
	return glyphCache.glyphs[fontName][chr].Copy()
}

// renderGlyph renders a glyph then caches and returns it
func renderGlyph(gc d2d.GraphicContext, fontName string, chr rune) *Glyph {
	gc.Save()
	defer gc.Restore()
	gc.BeginPath()
	width := gc.CreateStringPath(string(chr), 0, 0)
	path := gc.GetPath()
	return &Glyph{
		Path:  &path,
		Width: width,
	}
}

// Glyph represents a rune which has been converted to a Path and width
type Glyph struct {
	// path represents a glyph, it is always at (0, 0)
	Path *d2d.Path
	// Width of the glyph
	Width float64
}

// Copy Returns a copy of a Glyph
func (g *Glyph) Copy() *Glyph {
	return &Glyph{
		Path:  g.Path.Copy(),
		Width: g.Width,
	}
}

// Fill copies a glyph from the cache, and fills it
func (g *Glyph) Fill(gc d2d.GraphicContext, x, y float64) float64 {
	gc.Save()
	gc.BeginPath()
	gc.Translate(x, y)
	gc.Fill(g.Path)
	gc.Restore()
	return g.Width
}

// Stroke fetches a glyph from the cache, and strokes it
func (g *Glyph) Stroke(gc d2d.GraphicContext, x, y float64) float64 {
	gc.Save()
	gc.BeginPath()
	gc.Translate(x, y)
	gc.Stroke(g.Path)
	gc.Restore()
	return g.Width
}
