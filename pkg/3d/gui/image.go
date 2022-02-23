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

import (
	"image"

	"github.com/bhojpur/render/pkg/3d/texture"
)

// Image is a Panel which contains a single Image
type Image struct {
	Panel                    // Embedded panel
	tex   *texture.Texture2D // pointer to image texture
}

// NewImage creates and returns an image panel with the image
// from the specified image used as a texture.
// Initially the size of the panel content area is the exact size of the image.
func NewImage(imgfile string) (image *Image, err error) {

	tex, err := texture.NewTexture2DFromImage(imgfile)
	if err != nil {
		return nil, err
	}
	return NewImageFromTex(tex), nil
}

// NewImageFromRGBA creates and returns an image panel from the
// specified image
func NewImageFromRGBA(rgba *image.RGBA) *Image {

	tex := texture.NewTexture2DFromRGBA(rgba)
	return NewImageFromTex(tex)
}

// NewImageFromTex creates and returns an image panel from the specified texture2D
func NewImageFromTex(tex *texture.Texture2D) *Image {

	i := new(Image)
	i.Panel.Initialize(i, 0, 0)
	i.tex = tex
	i.Panel.SetContentSize(float32(i.tex.Width()), float32(i.tex.Height()))
	i.Material().AddTexture(i.tex)
	return i
}

// SetTexture changes the image texture to the specified texture2D.
// It returns a pointer to the previous texture.
func (i *Image) SetTexture(tex *texture.Texture2D) *texture.Texture2D {

	prevtex := i.tex
	i.Material().RemoveTexture(prevtex)
	i.tex = tex
	i.Panel.SetContentSize(float32(i.tex.Width()), float32(i.tex.Height()))
	i.Material().AddTexture(i.tex)
	return prevtex
}

// SetImage sets the image from the specified image file
func (i *Image) SetImage(imgfile string) error {

	tex, err := texture.NewTexture2DFromImage(imgfile)
	if err != nil {
		return err
	}
	i.SetTexture(tex)
	return nil
}
