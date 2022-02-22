package renderer

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

// It implements the scene renderer.

import (
	"github.com/bhojpur/render/pkg/gls"
)

type Postprocessor struct {
	Width    int32
	Height   int32
	Fbo      uint32
	Tex      uint32
	Vao      uint32
	Prg      *gls.Program
	screen   []float32
	Renderer *Renderer
}

func (r *Renderer) CreatePostprocessor(width, height int32, vertexShaderSource, fragmentShaderSource string) *Postprocessor {
	pp := &Postprocessor{
		Width:    width,
		Height:   height,
		Renderer: r,
		screen: []float32{
			// xyz		color		texture coords
			-1, 1, 0, 1, 1, 1, 0, 1,
			-1, -1, 0, 1, 1, 1, 0, 0,
			1, -1, 0, 1, 1, 1, 1, 0,
			1, 1, 0, 1, 1, 1, 1, 1,
			-1, 1, 0, 1, 1, 1, 0, 1,
			1, -1, 0, 1, 1, 1, 1, 0,
		},
	}

	pp.Fbo = r.gs.GenFramebuffer()
	r.gs.BindFramebuffer(pp.Fbo)

	// set up a texture to render into
	pp.Tex = r.gs.GenTexture()
	r.gs.BindTexture(gls.TEXTURE_2D, pp.Tex)
	r.gs.TexImage2D(gls.TEXTURE_2D, 0, gls.RGB, width, height, gls.RGB, gls.UNSIGNED_BYTE, nil)
	r.gs.TexParameteri(gls.TEXTURE_2D, gls.TEXTURE_WRAP_S, gls.CLAMP_TO_EDGE)
	r.gs.TexParameteri(gls.TEXTURE_2D, gls.TEXTURE_WRAP_T, gls.CLAMP_TO_EDGE)
	r.gs.TexParameteri(gls.TEXTURE_2D, gls.TEXTURE_MIN_FILTER, gls.NEAREST)
	r.gs.TexParameteri(gls.TEXTURE_2D, gls.TEXTURE_MAG_FILTER, gls.NEAREST)
	r.gs.BindTexture(gls.TEXTURE_2D, 0)
	r.gs.FramebufferTexture2D(gls.COLOR_ATTACHMENT0, gls.TEXTURE_2D, pp.Tex)

	// attach depth and stencil buffers
	rbo := r.gs.GenRenderbuffer()
	r.gs.BindRenderbuffer(rbo)
	r.gs.RenderbufferStorage(gls.DEPTH24_STENCIL8, int(width), int(height))
	r.gs.BindRenderbuffer(0)
	r.gs.FramebufferRenderbuffer(gls.DEPTH_STENCIL_ATTACHMENT, rbo)

	// check the framebuffer status
	if r.gs.CheckFramebufferStatus() != gls.FRAMEBUFFER_COMPLETE {
		log.Fatal("Can't create frame buffer")
	}
	r.gs.BindFramebuffer(0)

	// create the "screen" quad
	vbo := r.gs.GenBuffer()
	r.gs.BindBuffer(gls.ARRAY_BUFFER, vbo)
	r.gs.BufferData(gls.ARRAY_BUFFER, 4*len(pp.screen), pp.screen, gls.STATIC_DRAW)

	pp.Vao = r.gs.GenVertexArray()
	r.gs.BindVertexArray(pp.Vao)
	r.gs.BindBuffer(gls.ARRAY_BUFFER, vbo)
	var offset uint32

	// position attribute
	r.gs.VertexAttribPointer(0, 3, gls.FLOAT, false, 8*4, offset)
	r.gs.EnableVertexAttribArray(0)
	offset += 3 * 4

	// color attribute
	r.gs.VertexAttribPointer(1, 3, gls.FLOAT, false, 8*4, offset)
	r.gs.EnableVertexAttribArray(1)
	offset += 3 * 4

	// texture coord attribute
	r.gs.VertexAttribPointer(2, 2, gls.FLOAT, false, 8*4, offset)
	r.gs.EnableVertexAttribArray(2)
	offset += 2 * 4

	// the screen shaders
	pp.Prg = r.gs.NewProgram()
	pp.Prg.AddShader(gls.VERTEX_SHADER, vertexShaderSource)
	pp.Prg.AddShader(gls.FRAGMENT_SHADER, fragmentShaderSource)
	err := pp.Prg.Build()
	if err != nil {
		log.Fatal("can't create shader: %e", err)
	}

	return pp
}

func (pp *Postprocessor) Render(fbwidth, fbheight int, render func()) {
	// render into the low-res texture
	gs := pp.Renderer.gs
	gs.Viewport(0, 0, pp.Width, pp.Height)
	gs.BindFramebuffer(pp.Fbo)
	gs.Enable(gls.DEPTH_TEST)
	render()

	// show texture on screen
	gs.Viewport(0, 0, int32(fbwidth), int32(fbheight))
	gs.BindFramebuffer(0)
	gs.ClearColor(1, 1, 1, 1)
	gs.Clear(gls.COLOR_BUFFER_BIT)
	gs.UseProgram(pp.Prg)
	gs.Disable(gls.DEPTH_TEST)
	gs.BindTexture(gls.TEXTURE_2D, pp.Tex)
	gs.BindVertexArray(pp.Vao)
	gs.DrawArrays(gls.TRIANGLES, 0, int32(len(pp.screen)/8))
}
