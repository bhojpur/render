package light

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
	"github.com/bhojpur/render/pkg/3d/core"
	"github.com/bhojpur/render/pkg/3d/math32"
	"github.com/bhojpur/render/pkg/gls"
)

// Directional represents a directional, positionless light
type Directional struct {
	core.Node              // Embedded node
	color     math32.Color // Light color
	intensity float32      // Light intensity
	uni       gls.Uniform  // Uniform location cache
	udata     struct {     // Combined uniform data in 2 vec3:
		color    math32.Color   // Light color
		position math32.Vector3 // Light position
	}
}

// NewDirectional creates and returns a pointer of a new directional light
// the specified color and intensity.
func NewDirectional(color *math32.Color, intensity float32) *Directional {

	ld := new(Directional)
	ld.Node.Init(ld)

	ld.color = *color
	ld.intensity = intensity
	ld.uni.Init("DirLight")
	ld.SetColor(color)
	return ld
}

// SetColor sets the color of this light
func (ld *Directional) SetColor(color *math32.Color) {

	ld.color = *color
	ld.udata.color = ld.color
	ld.udata.color.MultiplyScalar(ld.intensity)
}

// Color returns the current color of this light
func (ld *Directional) Color() math32.Color {

	return ld.color
}

// SetIntensity sets the intensity of this light
func (ld *Directional) SetIntensity(intensity float32) {

	ld.intensity = intensity
	ld.udata.color = ld.color
	ld.udata.color.MultiplyScalar(ld.intensity)
}

// Intensity returns the current intensity of this light
func (ld *Directional) Intensity() float32 {

	return ld.intensity
}

// RenderSetup is called by the engine before rendering the scene
func (ld *Directional) RenderSetup(gs *gls.GLS, rinfo *core.RenderInfo, idx int) {

	// Calculates light position in camera coordinates and updates uniform
	var pos math32.Vector3
	ld.WorldPosition(&pos)
	pos4 := math32.Vector4{pos.X, pos.Y, pos.Z, 0.0}
	pos4.ApplyMatrix4(&rinfo.ViewMatrix)
	ld.udata.position.X = pos4.X
	ld.udata.position.Y = pos4.Y
	ld.udata.position.Z = pos4.Z

	// Transfer uniform data
	const vec3count = 2
	location := ld.uni.LocationIdx(gs, vec3count*int32(idx))
	gs.Uniform3fv(location, vec3count, &ld.udata.color.R)
}
