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
	"github.com/bhojpur/render/pkg/g3d/text"
	"github.com/bhojpur/render/pkg/math32"
)

// Style contains the styles for all GUI elements
type Style struct {
	Color         ColorStyle
	Font          *text.Font
	FontIcon      *text.Font
	Label         LabelStyle
	Button        ButtonStyles
	CheckRadio    CheckRadioStyles
	Edit          EditStyles
	ScrollBar     ScrollBarStyles
	Slider        SliderStyles
	Splitter      SplitterStyles
	Window        WindowStyles
	ItemScroller  ItemScrollerStyles
	Scroller      ScrollerStyle
	List          ListStyles
	DropDown      DropDownStyles
	Folder        FolderStyles
	Tree          TreeStyles
	ControlFolder ControlFolderStyles
	Menu          MenuStyles
	Table         TableStyles
	ImageButton   ImageButtonStyles
	TabBar        TabBarStyles
}

// ColorStyle defines the main colors used.
type ColorStyle struct {
	BgDark    math32.Color4
	BgMed     math32.Color4
	BgNormal  math32.Color4
	BgOver    math32.Color4
	Highlight math32.Color4
	Select    math32.Color4
	Text      math32.Color4
	TextDis   math32.Color4
}

// States that a GUI element can be in
const (
	StyleOver = iota + 1
	StyleFocus
	StyleDisabled
	StyleNormal
	StyleDef
)
