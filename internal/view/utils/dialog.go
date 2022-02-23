package utils

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

	"github.com/bhojpur/render/pkg/3d/gui"
	"github.com/bhojpur/render/pkg/3d/math32"
	"github.com/bhojpur/render/pkg/app"
)

type ErrorDialog struct {
	gui.Panel
	msg *gui.ImageLabel
	bok *gui.Button
}

func NewErrorDialog(width, height float32) *ErrorDialog {

	e := new(ErrorDialog)
	e.Initialize(e, width, height)
	e.SetBorders(2, 2, 2, 2)
	e.SetPaddings(4, 4, 4, 4)
	e.SetColor(math32.NewColor("White"))
	e.SetVisible(false)
	e.SetBounded(false)

	// Set vertical box layout for the whole panel
	l := gui.NewVBoxLayout()
	l.SetSpacing(4)
	e.SetLayout(l)

	// Creates error message label
	e.msg = gui.NewImageLabel("")
	e.msg.SetColor(math32.NewColor("black"))
	e.msg.SetLayoutParams(&gui.VBoxLayoutParams{Expand: 2, AlignH: gui.AlignWidth})
	e.Add(e.msg)

	// Creates button
	e.bok = gui.NewButton("OK")
	e.bok.SetLayoutParams(&gui.VBoxLayoutParams{Expand: 1, AlignH: gui.AlignCenter})
	e.bok.Subscribe(gui.OnClick, func(evname string, ev interface{}) {
		e.SetVisible(false)
	})
	e.Add(e.bok)

	return e
}

func (e *ErrorDialog) Show(msg string) {

	e.msg.SetText(msg)
	fmt.Println(msg)
	e.SetVisible(true)
	width, height := app.App().GetSize()
	px := (float32(width) - e.Width()) / 2
	py := (float32(height) - e.Height()) / 2
	e.SetPosition(px, py)
}
