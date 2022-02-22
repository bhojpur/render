package graphic

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
	"github.com/bhojpur/render/pkg/engine/geometry"
	"github.com/bhojpur/render/pkg/engine/material"
	"github.com/bhojpur/render/pkg/engine/math32"
	"github.com/bhojpur/render/pkg/engine/texture"
	"github.com/bhojpur/render/pkg/gls"
)

// Skybox is the Graphic that represents a skybox.
type Skybox struct {
	Graphic             // embedded graphic object
	uniMVm  gls.Uniform // model view matrix uniform location cache
	uniMVPm gls.Uniform // model view projection matrix uniform cache
	uniNm   gls.Uniform // normal matrix uniform cache
}

// SkyboxData contains the data necessary to locate the textures for a Skybox in a concise manner.
type SkyboxData struct {
	DirAndPrefix string
	Extension    string
	Suffixes     [6]string
}

// NewSkybox creates and returns a pointer to a Skybox with the specified textures.
func NewSkybox(data SkyboxData) (*Skybox, error) {

	skybox := new(Skybox)

	geom := geometry.NewCube(1)
	skybox.Graphic.Init(skybox, geom, gls.TRIANGLES)
	skybox.Graphic.SetCullable(false)

	for i := 0; i < 6; i++ {
		tex, err := texture.NewTexture2DFromImage(data.DirAndPrefix + data.Suffixes[i] + "." + data.Extension)
		if err != nil {
			return nil, err
		}
		matFace := material.NewStandard(math32.NewColor("white"))
		matFace.AddTexture(tex)
		matFace.SetSide(material.SideBack)
		matFace.SetUseLights(material.UseLightNone)

		// Disable writes to the depth buffer (call glDepthMask(GL_FALSE)).
		// This will cause every other object to draw over the skybox, making it always appear behind everything else.
		// It doesn't matter how small/big the skybox is as long as it's visible by the camera (within near/far planes).
		matFace.SetDepthMask(false)

		skybox.AddGroupMaterial(skybox, matFace, i)
	}

	// Creates uniforms
	skybox.uniMVm.Init("ModelViewMatrix")
	skybox.uniMVPm.Init("MVP")
	skybox.uniNm.Init("NormalMatrix")

	// The skybox should always be rendered last among the opaque objects
	skybox.SetRenderOrder(100)

	return skybox, nil
}

// RenderSetup is called by the engine before drawing the skybox geometry
// It is responsible to updating the current shader uniforms with
// the model matrices.
func (skybox *Skybox) RenderSetup(gs *gls.GLS, rinfo *core.RenderInfo) {

	mvm := *skybox.ModelViewMatrix()

	// Clear translation
	mvm[12] = 0
	mvm[13] = 0
	mvm[14] = 0
	// mvm.ExtractRotation(&rinfo.ViewMatrix) // TODO <- ExtractRotation does not work as expected?

	// Transfer mvp uniform
	location := skybox.uniMVm.Location(gs)
	gs.UniformMatrix4fv(location, 1, false, &mvm[0])

	// Calculates model view projection matrix and updates uniform
	var mvpm math32.Matrix4
	mvpm.MultiplyMatrices(&rinfo.ProjMatrix, &mvm)
	location = skybox.uniMVPm.Location(gs)
	gs.UniformMatrix4fv(location, 1, false, &mvpm[0])

	// Calculates normal matrix and updates uniform
	var nm math32.Matrix3
	nm.GetNormalMatrix(&mvm)
	location = skybox.uniNm.Location(gs)
	gs.UniformMatrix3fv(location, 1, false, &nm[0])
}
