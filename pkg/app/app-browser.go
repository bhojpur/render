//go:build wasm
// +build wasm

package app

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
	"github.com/bhojpur/render/pkg/engine/renderer"
	"github.com/bhojpur/render/pkg/engine/window"
	"syscall/js"
	"time"
)

// Default canvas Id
const canvasId = "bhojpur-canvas"

// Application
type Application struct {
	window.IWindow                    // Embedded WebGLCanvas
	keyState       *window.KeyState   // Keep track of keyboard state
	renderer       *renderer.Renderer // Renderer object
	startTime      time.Time          // Application start time
	frameStart     time.Time          // Frame start time
	frameDelta     time.Duration      // Duration of last frame
	exit           bool
	cbid           js.Value
}

// App returns the Application singleton, creating it the first time.
func App() *Application {

	// Return singleton if already created
	if a != nil {
		return a
	}
	a = new(Application)
	// Initialize window
	err := window.Init(canvasId)
	if err != nil {
		panic(err)
	}
	a.IWindow = window.Get()
	// TODO audio setup here
	a.keyState = window.NewKeyState(a) // Create KeyState
	// Create renderer and add default shaders
	a.renderer = renderer.NewRenderer(a.Gls())
	err = a.renderer.AddDefaultShaders()
	if err != nil {
		panic(fmt.Errorf("AddDefaultShaders:%v", err))
	}
	return a
}

// Run starts the update loop.
// It calls the user-provided update function every frame.
func (a *Application) Run(update func(rend *renderer.Renderer, deltaTime time.Duration)) {

	// Create channel so later we can prevent application from finishing while we wait for callbacks
	done := make(chan bool)

	// Initialize start and frame time
	a.startTime = time.Now()
	a.frameStart = time.Now()

	// Set up recurring calls to user's update function
	var tick js.Func
	tick = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		// Update frame start and frame delta
		now := time.Now()
		a.frameDelta = now.Sub(a.frameStart)
		a.frameStart = now
		// Call user's update function
		update(a.renderer, a.frameDelta)
		// Set up new callback if not exiting
		if !a.exit {
			a.cbid = js.Global().Call("requestAnimationFrame", tick)
		} else {
			a.Dispatch(OnExit, nil)
			done <- true // Write to done channel to exit the app
		}
		return nil
	})
	defer tick.Release()

	a.cbid = js.Global().Call("requestAnimationFrame", tick)

	// Read from done channel
	// This channel will be empty (except when we want to exit the app)
	// It keeps the app from finishing while we wait for the next call to tick()
	<-done

	// Destroy the window
	a.IWindow.Destroy()
}

// Exit exits the app.
func (a *Application) Exit() {

	a.exit = true
}

// Renderer returns the application's renderer.
func (a *Application) Renderer() *renderer.Renderer {

	return a.renderer
}

// KeyState returns the application's KeyState.
func (a *Application) KeyState() *window.KeyState {

	return a.keyState
}

// RunTime returns the elapsed duration since the call to Run().
func (a *Application) RunTime() time.Duration {

	return time.Now().Sub(a.startTime)
}
