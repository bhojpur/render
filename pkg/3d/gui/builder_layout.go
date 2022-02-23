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

// BuilderLayoutHBox is builder for HBox layout
type BuilderLayoutHBox struct{}

// BuildLayout builds and returns an HBoxLayout with the specified attributes
func (bl *BuilderLayoutHBox) BuildLayout(b *Builder, am map[string]interface{}) (ILayout, error) {

	// Creates layout and sets optional spacing
	l := NewHBoxLayout()
	var spacing float32
	if sp := am[AttribSpacing]; sp != nil {
		spacing = sp.(float32)
	}
	l.SetSpacing(spacing)

	// Sets optional horizontal alignment
	if ah := am[AttribAlignh]; ah != nil {
		l.SetAlignH(ah.(Align))
	}

	// Sets optional minheight flag
	if mh := am[AttribAutoHeight]; mh != nil {
		l.SetAutoHeight(mh.(bool))
	}

	// Sets optional minwidth flag
	if mw := am[AttribAutoWidth]; mw != nil {
		l.SetAutoWidth(mw.(bool))
	}
	return l, nil
}

// BuildParams builds and returns a pointer to HBoxLayoutParams with the specified attributes
func (bl *BuilderLayoutHBox) BuildParams(b *Builder, am map[string]interface{}) (interface{}, error) {

	// Creates layout parameters with default values
	params := HBoxLayoutParams{Expand: 0, AlignV: AlignNone}

	// Sets optional expand parameter
	if expand := am[AttribExpand]; expand != nil {
		params.Expand = expand.(float32)
	}

	// Sets optional align parameter
	if alignv := am[AttribAlignv]; alignv != nil {
		params.AlignV = alignv.(Align)
	}
	return &params, nil
}

//
// BuilderLayoutVBox is builder for VBox layout
//
type BuilderLayoutVBox struct{}

// BuildLayout builds and returns an VBoxLayout with the specified attributes
func (bl *BuilderLayoutVBox) BuildLayout(b *Builder, am map[string]interface{}) (ILayout, error) {

	// Creates layout and sets optional spacing
	l := NewVBoxLayout()
	var spacing float32
	if sp := am[AttribSpacing]; sp != nil {
		spacing = sp.(float32)
	}
	l.SetSpacing(spacing)

	// Sets optional vertical alignment
	if av := am[AttribAlignh]; av != nil {
		l.SetAlignV(av.(Align))
	}

	// Sets optional minheight flag
	if mh := am[AttribAutoHeight]; mh != nil {
		l.SetAutoHeight(mh.(bool))
	}

	// Sets optional minwidth flag
	if mw := am[AttribAutoWidth]; mw != nil {
		l.SetAutoWidth(mw.(bool))
	}
	return l, nil
}

// BuildParams builds and returns a pointer to VBoxLayoutParams with the specified attributes
func (bl *BuilderLayoutVBox) BuildParams(b *Builder, am map[string]interface{}) (interface{}, error) {

	// Creates layout parameters with default values
	params := VBoxLayoutParams{Expand: 0, AlignH: AlignNone}

	// Sets optional expand parameter
	if expand := am[AttribExpand]; expand != nil {
		params.Expand = expand.(float32)
	}

	// Sets optional align parameter
	if alignh := am[AttribAlignh]; alignh != nil {
		params.AlignH = alignh.(Align)
	}
	return &params, nil
}

//
// BuilderLayoutGrid is builder for Grid layout
//
type BuilderLayoutGrid struct{}

// BuildLayout builds and returns a GridLayout with the specified attributes
func (bl *BuilderLayoutGrid) BuildLayout(b *Builder, am map[string]interface{}) (ILayout, error) {

	// Get number of columns
	v := am[AttribCols]
	if v == nil {
		return nil, b.err(am, AttribCols, "Number of columns must be supplied")
	}
	cols := v.(int)
	if cols <= 0 {
		return nil, b.err(am, AttribCols, "Invalid number of columns")
	}

	// Creates layout
	l := NewGridLayout(cols)

	// Sets optional horizontal alignment
	if ah := am[AttribAlignh]; ah != nil {
		l.SetAlignH(ah.(Align))
	}

	// Sets optional vertical alignment
	if av := am[AttribAlignv]; av != nil {
		l.SetAlignV(av.(Align))
	}

	// Sets optional horizontal expand flag
	if eh := am[AttribExpandh]; eh != nil {
		l.SetExpandH(eh.(bool))
	}

	// Sets optional vertical expand flag
	if ev := am[AttribExpandv]; ev != nil {
		l.SetExpandV(ev.(bool))
	}

	return l, nil
}

// BuildParams builds and returns a pointer to GridLayoutParams with the specified attributes
func (bl *BuilderLayoutGrid) BuildParams(b *Builder, am map[string]interface{}) (interface{}, error) {

	// Creates layout parameters with default values
	params := GridLayoutParams{
		ColSpan: 0,
		AlignH:  AlignNone,
		AlignV:  AlignNone,
	}

	// Sets optional colspan
	if cs := am[AttribColSpan]; cs != nil {
		params.ColSpan = cs.(int)
	}

	// Sets optional alignh parameter
	if align := am[AttribAlignh]; align != nil {
		params.AlignH = align.(Align)
	}

	// Sets optional alignv parameter
	if align := am[AttribAlignv]; align != nil {
		params.AlignV = align.(Align)
	}
	return &params, nil
}

//
// BuilderLayoutDock is builder for Dock layout
//
type BuilderLayoutDock struct{}

// BuildLayout builds and returns a DockLayout with the specified attributes
func (bl *BuilderLayoutDock) BuildLayout(b *Builder, am map[string]interface{}) (ILayout, error) {

	return NewDockLayout(), nil
}

// BuildParams builds and returns a pointer to DockLayoutParams with the specified attributes
func (bl *BuilderLayoutDock) BuildParams(b *Builder, am map[string]interface{}) (interface{}, error) {

	edge := am[AttribEdge]
	if edge == nil {
		return nil, b.err(am, AttribEdge, "Edge name not found")
	}
	params := DockLayoutParams{Edge: edge.(int)}
	return &params, nil
}