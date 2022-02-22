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
	"github.com/bhojpur/render/pkg/engine/core"
	"github.com/bhojpur/render/pkg/engine/math32"
	"github.com/bhojpur/render/pkg/gls"
)

// Ambient represents an ambient light
type Ambient struct {
	core.Node              // Embedded node
	color     math32.Color // Light color
	intensity float32      // Light intensity
	uni       gls.Uniform  // Uniform location cache
}

// NewAmbient returns a pointer to a new ambient color with the specified
// color and intensity
func NewAmbient(color *math32.Color, intensity float32) *Ambient {

	la := new(Ambient)
	la.Node.Init(la)
	la.color = *color
	la.intensity = intensity
	la.uni.Init("AmbientLightColor")
	return la
}

// SetColor sets the color of this light
func (la *Ambient) SetColor(color *math32.Color) {

	la.color = *color
}

// Color returns the current color of this light
func (la *Ambient) Color() math32.Color {

	return la.color
}

// SetIntensity sets the intensity of this light
func (la *Ambient) SetIntensity(intensity float32) {

	la.intensity = intensity
}

// Intensity returns the current intensity of this light
func (la *Ambient) Intensity() float32 {

	return la.intensity
}

// RenderSetup is called by the engine before rendering the scene
func (la *Ambient) RenderSetup(gs *gls.GLS, rinfo *core.RenderInfo, idx int) {

	color := la.color
	color.MultiplyScalar(la.intensity)
	location := la.uni.LocationIdx(gs, int32(idx))
	gs.Uniform3f(location, color.R, color.G, color.B)
}
