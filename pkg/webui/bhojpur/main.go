package main

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

// It uses of the Bhojpur Render - Wasm Canvas library,  and avoids context calls for drawing.

import (
	"image/color"
	"syscall/js"
	"time"

	"github.com/bhojpur/render/pkg/g2d/canvas"
	"github.com/bhojpur/render/pkg/g2d/img"
	"github.com/bhojpur/render/pkg/g2d/kit"
)

// Default canvas Id
const defaultCanvasId = "canvas"

type gameState struct{ laserX, laserY, directionX, directionY, laserSize float64 }

var done chan struct{}

var cvs *canvas.WebGlCanvas2D
var width float64
var height float64

var gs = gameState{laserSize: 35, directionX: 13.7, directionY: -13.7, laserX: 40, laserY: 40}

// This specifies how long a delay between calls to 'render'.     To get Frame Rate,   1s / renderDelay
var renderDelay time.Duration = 20 * time.Millisecond

func main() {

	FrameRate := time.Second / renderDelay
	println("Hello, web browser frames per second:", FrameRate)
	//cvs, _ = canvas.NewCanvas2D(defaultCanvasId, true)

	cvs, _ = canvas.NewCanvas2D(false)
	cvs.Create(int(js.Global().Get("innerWidth").Float()*0.9), int(js.Global().Get("innerHeight").Float()*0.9)) // Make Canvas 90% of window size.  For testing rendering canvas smaller than full windows

	height = float64(cvs.Height())
	width = float64(cvs.Width())

	cvs.Start(60, Render)

	//go doEvery(renderDelay, Render) // Kick off the Render function as go routine as it never returns
	<-done
}

// Helper function which calls the required func (in this case 'render') every time.Duration,  Call as a go-routine to prevent blocking, as this never returns
func doEvery(d time.Duration, f func(time.Time)) {
	for x := range time.Tick(d) {
		f(x)
	}
}

// Called from the 'requestAnnimationFrame' function.   It may also be called seperatly from a 'doEvery' function, if the user prefers drawing to be seperate from the annimationFrame callback
func Render(gc *img.GraphicContext) bool {

	if gs.laserX+gs.directionX > width-gs.laserSize || gs.laserX+gs.directionX < gs.laserSize {
		gs.directionX = -gs.directionX
	}
	if gs.laserY+gs.directionY > height-gs.laserSize || gs.laserY+gs.directionY < gs.laserSize {
		gs.directionY = -gs.directionY
	}

	gc.SetFillColor(color.RGBA{0xff, 0xff, 0xff, 0xff})
	gc.Clear()
	// move red laser
	gs.laserX += gs.directionX
	gs.laserY += gs.directionY

	// draws red 🔴 laser
	gc.SetFillColor(color.RGBA{0xff, 0x00, 0xff, 0xff})
	gc.SetStrokeColor(color.RGBA{0xff, 0x00, 0xff, 0xff})

	gc.BeginPath()
	//gc.ArcTo(gs.laserX, gs.laserY, gs.laserSize, gs.laserSize, 0, math.Pi*2)
	kit.Circle(gc, gs.laserX, gs.laserY, gs.laserSize)
	gc.FillStroke()
	gc.Close()

	return true
}
