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

import "testing"

// Test simple matrix operations
func Test(t *testing.T) {

	m := NewMatrix()

	// m.Get(1, 1) // panics with "runtime error: index out of range" as expected

	m.Set(2, 4, true)
	m.Set(3, 2, true)
	if m.Get(1, 1) != false {
		t.Error("Get failed")
	}
	if m.Get(2, 4) != true {
		t.Error("Get failed")
	}
	if m.Get(3, 2) != true {
		t.Error("Get failed")
	}

	m.Set(2, 4, false)
	m.Set(3, 2, false)
	if m.Get(2, 4) != false {
		t.Error("Get failed")
	}
	if m.Get(3, 2) != false {
		t.Error("Get failed")
	}

	// m.Get(100, 100) // panics with "runtime error: index out of range" as expected
}
