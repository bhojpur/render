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

// DockLayout is the layout for docking panels to the internal edges of their parent.
type DockLayout struct {
}

// DockLayoutParams specifies the edge to dock to.
type DockLayoutParams struct {
	Edge int
}

// The different types of docking.
const (
	DockTop = iota + 1
	DockRight
	DockBottom
	DockLeft
	DockCenter
)

// NewDockLayout returns a pointer to a new DockLayout.
func NewDockLayout() *DockLayout {

	return new(DockLayout)
}

// Recalc (which satisfies the ILayout interface) recalculates the positions and sizes of the children panels.
func (dl *DockLayout) Recalc(ipan IPanel) {

	pan := ipan.GetPanel()
	width := pan.Width()
	topY := float32(0)
	bottomY := pan.Height()
	leftX := float32(0)
	rightX := width

	// Top and bottom first
	for _, iobj := range pan.Children() {
		child := iobj.(IPanel).GetPanel()
		if child.layoutParams == nil {
			continue
		}
		params := child.layoutParams.(*DockLayoutParams)
		if params.Edge == DockTop {
			child.SetPosition(0, topY)
			topY += child.Height()
			child.SetWidth(width)
			continue
		}
		if params.Edge == DockBottom {
			child.SetPosition(0, bottomY-child.Height())
			bottomY -= child.Height()
			child.SetWidth(width)
			continue
		}
	}
	// Left and right
	for _, iobj := range pan.Children() {
		child := iobj.(IPanel).GetPanel()
		if child.layoutParams == nil {
			continue
		}
		params := child.layoutParams.(*DockLayoutParams)
		if params.Edge == DockLeft {
			child.SetPosition(leftX, topY)
			leftX += child.Width()
			child.SetHeight(bottomY - topY)
			continue
		}
		if params.Edge == DockRight {
			child.SetPosition(rightX-child.Width(), topY)
			rightX -= child.Width()
			child.SetHeight(bottomY - topY)
			continue
		}
	}
	// Center (only the first found)
	for _, iobj := range pan.Children() {
		child := iobj.(IPanel).GetPanel()
		if child.layoutParams == nil {
			continue
		}
		params := child.layoutParams.(*DockLayoutParams)
		if params.Edge == DockCenter {
			child.SetPosition(leftX, topY)
			child.SetSize(rightX-leftX, bottomY-topY)
			break
		}
	}
}
