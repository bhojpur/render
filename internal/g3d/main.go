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

import (
	"time"

	engine "github.com/bhojpur/render/pkg/app"
	"github.com/bhojpur/render/pkg/g3d/camera"
	"github.com/bhojpur/render/pkg/g3d/core"
	"github.com/bhojpur/render/pkg/g3d/gui"
	"github.com/bhojpur/render/pkg/g3d/light"
	"github.com/bhojpur/render/pkg/g3d/loader/obj"
	"github.com/bhojpur/render/pkg/g3d/renderer"
	"github.com/bhojpur/render/pkg/g3d/window"
	"github.com/bhojpur/render/pkg/gls"
	"github.com/bhojpur/render/pkg/math32"
)

func main() {

	// Create a 3D Rendering application and scene
	a := engine.BhojpurApp3D()
	scene := core.NewNode()

	// Set the scene to be managed by the gui manager
	gui.Manager().Set(scene)

	// Create perspective camera
	cam := camera.New(1)
	cam.SetPosition(0, 0, 3)
	scene.Add(cam)

	// Set up orbit control for the camera
	camera.NewOrbitControl(cam)

	// Set up callback to update viewport and camera aspect ratio when the window is resized
	onResize := func(evname string, ev interface{}) {
		// Get framebuffer size and update viewport accordingly
		width, height := a.GetSize()
		a.Gls().Viewport(0, 0, int32(width), int32(height))
		// Update the camera's aspect ratio
		cam.SetAspect(float32(width) / float32(height))
	}
	a.Subscribe(window.OnWindowSize, onResize)
	onResize("", nil)

	// Decode model in in OBJ format
	dec, err := obj.Decode("./internal/g3d/gopher.obj", "./internal/g3d/gopher.mtl")
	if err != nil {
		panic(err.Error())
	}

	// Create a new node with all the objects in the decoded file and adds it to the scene
	group, err := dec.NewGroup()
	if err != nil {
		panic(err.Error())
	}
	group.SetScale(0.3, 0.3, 0.3)
	group.SetPosition(0.0, -0.8, -0.2)
	scene.Add(group)

	// Create and add ambient light to scene
	scene.Add(light.NewAmbient(&math32.Color{1.0, 1.0, 1.0}, 0.8))

	// Create and add directional white light to the scene
	dirLight := light.NewDirectional(&math32.Color{1, 1, 1}, 1.0)
	dirLight.SetPosition(1, 0, 0)
	scene.Add(dirLight)

	// Set background color to gray
	a.Gls().ClearColor(0.5, 0.5, 0.5, 1.0)

	// Run the application
	a.Run(func(renderer *renderer.Renderer, deltaTime time.Duration) {
		a.Gls().Clear(gls.DEPTH_BUFFER_BIT | gls.STENCIL_BUFFER_BIT | gls.COLOR_BUFFER_BIT)
		renderer.Render(scene, cam)
	})
}
