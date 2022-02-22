package gls

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

// ShaderDefines is a store of shader defines ("#define <key> <value>").
type ShaderDefines map[string]string

// NewShaderDefines creates and returns a pointer to a ShaderDefines object.
func NewShaderDefines() *ShaderDefines {

	sd := ShaderDefines(make(map[string]string))
	return &sd
}

// Set sets a shader define with the specified value.
func (sd *ShaderDefines) Set(name, value string) {

	(*sd)[name] = value
}

// Unset removes the specified name from the shader defines.
func (sd *ShaderDefines) Unset(name string) {

	delete(*sd, name)
}

// Add adds to this ShaderDefines all the key-value pairs in the specified ShaderDefines.
func (sd *ShaderDefines) Add(other *ShaderDefines) {

	for k, v := range map[string]string(*other) {
		(*sd)[k] = v
	}
}

// Equals compares two ShaderDefines and return true if they contain the same key-value pairs.
func (sd *ShaderDefines) Equals(other *ShaderDefines) bool {

	if sd == nil && other == nil {
		return true
	}
	if sd != nil && other != nil {
		if len(*sd) != len(*other) {
			return false
		}
		for k := range map[string]string(*sd) {
			v1, ok1 := (*sd)[k]
			v2, ok2 := (*other)[k]
			if v1 != v2 || ok1 != ok2 {
				return false
			}
		}
		return true
	}
	// One is nil and the other is not nil
	return false
}
