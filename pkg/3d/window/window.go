package window

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

// It abstracts a platform-specific window.
// Depending on the build tags it can be a GLFW desktop window or a browser WebGlCanvas.

import (
	"fmt"

	"github.com/bhojpur/render/pkg/3d/core"
	"github.com/bhojpur/render/pkg/3d/util/logger"
	"github.com/bhojpur/render/pkg/gls"
)

// Package logger
var log = logger.New("WIN", logger.Default)

// IWindow singleton
var win IWindow

// Get returns the IWindow singleton.
func Get() IWindow {
	// Return singleton if already created
	if win != nil {
		return win
	}
	panic(fmt.Errorf("need to call window.Init() first"))
}

// IWindow is the interface for all windows
type IWindow interface {
	core.IDispatcher
	Gls() *gls.GLS
	GetFramebufferSize() (width int, height int)
	GetSize() (width int, height int)
	GetScale() (x float64, y float64)
	CreateCursor(imgFile string, xhot, yhot int) (Cursor, error)
	SetCursor(cursor Cursor)
	DisposeAllCustomCursors()
	Destroy()
	FullScreen() bool
	SetFullScreen(full bool)
}

// Key corresponds to a keyboard key.
type Key int

// ModifierKey corresponds to a set of modifier keys (bitmask).
type ModifierKey int

// MouseButton corresponds to a mouse button.
type MouseButton int

// InputMode corresponds to an input mode.
type InputMode int

// InputMode corresponds to an input mode.
type CursorMode int

// Cursor corresponds to a Bhojpur Render standard or user-created cursor icon.
type Cursor int

// Standard cursors for Bhojpur Render.
const (
	ArrowCursor = Cursor(iota)
	IBeamCursor
	CrosshairCursor
	HandCursor
	HResizeCursor
	VResizeCursor
	DiagResize1Cursor
	DiagResize2Cursor
	CursorLast = DiagResize2Cursor
)

// Window event names. See availability per platform below ("x" indicates available).
const ( //                               Desktop | Browser |
	OnWindowFocus = "w.OnWindowFocus" //    x    |    x    |
	OnWindowPos   = "w.OnWindowPos"   //    x    |         |
	OnWindowSize  = "w.OnWindowSize"  //    x    |         |
	OnKeyUp       = "w.OnKeyUp"       //    x    |    x    |
	OnKeyDown     = "w.OnKeyDown"     //    x    |    x    |
	OnKeyRepeat   = "w.OnKeyRepeat"   //    x    |         |
	OnChar        = "w.OnChar"        //    x    |    x    |
	OnCursor      = "w.OnCursor"      //    x    |    x    |
	OnMouseUp     = "w.OnMouseUp"     //    x    |    x    |
	OnMouseDown   = "w.OnMouseDown"   //    x    |    x    |
	OnScroll      = "w.OnScroll"      //    x    |    x    |
)

// PosEvent describes a windows position changed event
type PosEvent struct {
	Xpos int
	Ypos int
}

// SizeEvent describers a window size changed event
type SizeEvent struct {
	Width  int
	Height int
}

// KeyEvent describes a window key event
type KeyEvent struct {
	Key  Key
	Mods ModifierKey
}

// CharEvent describes a window char event
type CharEvent struct {
	Char rune
	Mods ModifierKey
}

// MouseEvent describes a mouse event over the window
type MouseEvent struct {
	Xpos   float32
	Ypos   float32
	Button MouseButton
	Mods   ModifierKey
}

// CursorEvent describes a cursor position changed event
type CursorEvent struct {
	Xpos float32
	Ypos float32
	Mods ModifierKey
}

// ScrollEvent describes a scroll event
type ScrollEvent struct {
	Xoffset float32
	Yoffset float32
	Mods    ModifierKey
}

// FocusEvent describes a focus event
type FocusEvent struct {
	Focused bool
}
