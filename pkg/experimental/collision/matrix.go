package collision

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

// Matrix is a triangular collision matrix indicating which pairs of bodies are colliding.
type Matrix [][]bool

// NewMatrix creates and returns a pointer to a new collision Matrix.
func NewMatrix() Matrix {

	return make([][]bool, 0)
}

// Set sets whether i and j are colliding.
func (m *Matrix) Set(i, j int, val bool) {

	var s, l int
	if i < j {
		s = i
		l = j
	} else {
		s = j
		l = i
	}
	diff := s + 1 - len(*m)
	if diff > 0 {
		for i := 0; i < diff; i++ {
			*m = append(*m, make([]bool, 0))
		}
	}
	for idx := range *m {
		diff = l + 1 - len((*m)[idx]) - idx
		if diff > 0 {
			for i := 0; i < diff; i++ {
				(*m)[idx] = append((*m)[idx], false)
			}
		}
	}
	(*m)[s][l-s] = val
}

// Get returns whether i and j are colliding.
func (m *Matrix) Get(i, j int) bool {

	var s, l int
	if i < j {
		s = i
		l = j
	} else {
		s = j
		l = i
	}
	return (*m)[s][l-s]
}

// Reset clears all values.
func (m *Matrix) Reset() {

	*m = make([][]bool, 0)
}
