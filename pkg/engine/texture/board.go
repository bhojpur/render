package texture

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
	"github.com/bhojpur/render/pkg/engine/math32"
	"github.com/bhojpur/render/pkg/gls"
)

// NewBoard creates and returns a pointer to a new checker board 2D texture.
// A checker board texture contains up to 4 different colors arranged in
// the following order:
//  +------+------+
//  |      |      |
//  |  c3  |  c4  |
//  |      |      |
//  +------+------+
//  |      |      |
//  |  c1  |  c2  | height (pixels)
//  |      |      |
//  +------+------+
//    width
//  (pixels)
//
func NewBoard(width, height int, c1, c2, c3, c4 *math32.Color, alpha float32) *Texture2D {

	// Generates texture data
	data := make([]float32, width*height*4*4)
	colorData := func(sx, sy int, c *math32.Color) {
		for y := sy; y < sy+height; y++ {
			for x := sx; x < sx+width; x++ {
				pos := (x + y*2*width) * 4
				data[pos] = c.R
				data[pos+1] = c.G
				data[pos+2] = c.B
				data[pos+3] = alpha
			}
		}
	}
	colorData(0, 0, c1)
	colorData(width, 0, c2)
	colorData(0, height, c3)
	colorData(width, height, c4)

	// Creates, initializes and returns board texture object
	return NewTexture2DFromData(width*2, height*2, gls.RGBA, gls.FLOAT, gls.RGBA8, data)
}
