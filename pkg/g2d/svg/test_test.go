package svg_test

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

// It gives test coverage with the command:
// go test -cover ./... | grep -v "no test"

import (
	"testing"

	d2d "github.com/bhojpur/render/pkg/g2d/draw"
	svgkit "github.com/bhojpur/render/pkg/g2d/svg"
)

type sample func(gc d2d.GraphicContext, ext string) (string, error)

func test(t *testing.T, draw sample) {
	// Initialize the graphic context on a .pdf document
	dest := svgkit.NewSvg()
	gc := svgkit.NewGraphicContext(dest)
	// Draw sample
	output, err := draw(gc, "svg")
	if err != nil {
		t.Errorf("Drawing %q failed: %v", output, err)
		return
	}
	err = svgkit.SaveToSvgFile(output, dest)
	if err != nil {
		t.Errorf("Saving %q failed: %v", output, err)
	}
}
