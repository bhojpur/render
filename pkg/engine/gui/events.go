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
	"github.com/bhojpur/render/pkg/engine/window"
)

// Core events sent by the GUI manager.
// The target panel is the panel immediately under the mouse cursor.
const (
	// Events sent to target panel's lowest subscribed ancestor
	OnMouseDown = window.OnMouseDown // Any mouse button is pressed
	OnMouseUp   = window.OnMouseUp   // Any mouse button is released
	OnScroll    = window.OnScroll    // Scrolling mouse wheel

	// Events sent to all panels except the ancestors of the target panel
	OnMouseDownOut = "gui.OnMouseDownOut" // Any mouse button is pressed
	OnMouseUpOut   = "gui.OnMouseUpOut"   // Any mouse button is released

	// Event sent to new target panel and all of its ancestors up to (not including) the common ancestor of the new and old targets
	OnCursorEnter = "gui.OnCursorEnter" // Cursor entered the panel or a descendant
	// Event sent to old target panel and all of its ancestors up to (not including) the common ancestor of the new and old targets
	OnCursorLeave = "gui.OnCursorLeave" // Cursor left the panel or a descendant
	// Event sent to the cursor-focused IDispatcher if any, else sent to target panel's lowest subscribed ancestor
	OnCursor = window.OnCursor // Cursor is over the panel

	// Event sent to the new key-focused IDispatcher, specified on a call to gui.Manager().SetKeyFocus(core.IDispatcher)
	OnFocus = "gui.OnFocus" // All keyboard events will be exclusively sent to the receiving IDispatcher
	// Event sent to the previous key-focused IDispatcher when another panel is key-focused
	OnFocusLost = "gui.OnFocusLost" // Keyboard events will stop being sent to the receiving IDispatcher

	// Events sent to the key-focused IDispatcher
	OnKeyDown   = window.OnKeyDown   // A key is pressed
	OnKeyUp     = window.OnKeyUp     // A key is released
	OnKeyRepeat = window.OnKeyRepeat // A key was pressed and is now automatically repeating
	OnChar      = window.OnChar      // A unicode key is pressed
)

const (
	OnResize     = "gui.OnResize"     // Panel size changed (no parameters)
	OnEnable     = "gui.OnEnable"     // Panel enabled/disabled (no parameters)
	OnClick      = "gui.OnClick"      // Widget clicked by mouse left button or via key press
	OnChange     = "gui.OnChange"     // Value was changed. Emitted by List, DropDownList, CheckBox and Edit
	OnRadioGroup = "gui.OnRadioGroup" // Radio button within a group changed state
)
