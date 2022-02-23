//go:generate rendershaders -in=. -out=sources.go -pkg=shaders -v

package shaders

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

// It contains the several shaders used by the engine. It generates shaders
// sources from this directory and include directory *.glsl files

// ProgramInfo contains information for a registered shader program
type ProgramInfo struct {
	Vertex   string // Vertex shader name
	Fragment string // Fragment shader name
	Geometry string // Geometry shader name (optional)
}

// AddInclude adds a chunk of shader code to the default shaders registry
// which can be included in a shader using the "#include <name>" directive
func AddInclude(name string, source string) {

	if len(name) == 0 || len(source) == 0 {
		panic("Invalid include name and/or source")
	}
	includeMap[name] = source
}

// AddShader add a shader to default shaders registry.
// The specified name can be used when adding programs to the registry
func AddShader(name string, source string) {

	if len(name) == 0 || len(source) == 0 {
		panic("Invalid shader name and/or source")
	}
	shaderMap[name] = source
}

// AddProgram adds a shader program to the default registry of programs.
// Currently up to 3 shaders: vertex, fragment and geometry (optional) can be specified.
func AddProgram(name string, vertex string, frag string, others ...string) {

	if len(name) == 0 || len(vertex) == 0 || len(frag) == 0 {
		panic("Program and/or shader name empty")
	}
	if shaderMap[vertex] == "" {
		panic("Invalid vertex shader name")
	}
	if shaderMap[frag] == "" {
		panic("Invalid vertex shader name")
	}
	var geom = ""
	if len(others) > 0 {
		geom = others[0]
		if shaderMap[geom] == "" {
			panic("Invalid geometry shader name")
		}
	}
	programMap[name] = ProgramInfo{
		Vertex:   vertex,
		Fragment: frag,
		Geometry: geom,
	}
}

// Includes returns list with the names of all include chunks currently in the default shaders registry.
func Includes() []string {

	list := make([]string, 0)
	for name := range includeMap {
		list = append(list, name)
	}
	return list
}

// IncludeSource returns the source code of the specified shader include chunk.
// If the name is not found an empty string is returned.
func IncludeSource(name string) string {

	return includeMap[name]
}

// Shaders returns list with the names of all shaders currently in the default shaders registry.
func Shaders() []string {

	list := make([]string, 0)
	for name := range shaderMap {
		list = append(list, name)
	}
	return list
}

// ShaderSource returns the source code of the specified shader in the default shaders registry.
// If the name is not found an empty string is returned
func ShaderSource(name string) string {

	return shaderMap[name]
}

// Programs returns list with the names of all programs currently in the default shaders registry.
func Programs() []string {

	list := make([]string, 0)
	for name := range programMap {
		list = append(list, name)
	}
	return list
}

// GetProgramInfo returns ProgramInfo struct for the specified program name
// in the default shaders registry
func GetProgramInfo(name string) ProgramInfo {

	return programMap[name]
}
