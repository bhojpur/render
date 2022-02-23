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

// FillLayout is the simple layout where the assigned panel "fills" its parent in the specified dimension(s)
type FillLayout struct {
	width  bool
	height bool
}

// NewFillLayout creates and returns a pointer of a new fill layout
func NewFillLayout(width, height bool) *FillLayout {

	f := new(FillLayout)
	f.width = width
	f.height = height
	return f
}

// Recalc is called by the panel which has this layout
func (f *FillLayout) Recalc(ipan IPanel) {

	parent := ipan.GetPanel()
	children := parent.Children()
	if len(children) == 0 {
		return
	}
	child := children[0].(IPanel).GetPanel()

	if f.width {
		child.SetWidth(parent.ContentWidth())
	}
	if f.height {
		child.SetHeight(parent.ContentHeight())
	}
}
