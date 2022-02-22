package vorbis

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

// It implements the Go bindings of a subset (only one function) of the functions of the libvorbis library
// See API reference at: https://xiph.org/vorbis/doc/libvorbis/reference.html

// #cgo darwin,amd64  CFLAGS:  -DGO_DARWIN  -I/usr/include/vorbis -I/usr/local/include/vorbis
// #cgo darwin,arm64  CFLAGS:  -DGO_DARWIN  -I/opt/homebrew/include -I/opt/homebrew/include/vorbis
// #cgo freebsd       CFLAGS:  -DGO_FREEBSD -I/usr/local/include/vorbis
// #cgo linux         CFLAGS:  -DGO_LINUX   -I/usr/include/vorbis
// #cgo windows       CFLAGS:  -DGO_WINDOWS -I${SRCDIR}/../windows/libvorbis-1.3.5/include/vorbis -I${SRCDIR}/../windows/libogg-1.3.3/include
// #cgo darwin,amd64  LDFLAGS: -L/usr/lib -L/usr/local/lib -lvorbis
// #cgo darwin,arm64  LDFLAGS: -L/opt/homebrew/lib -lvorbis
// #cgo freebsd       LDFLAGS: -L/usr/local/lib -lvorbis
// #cgo linux         LDFLAGS: -lvorbis
// #cgo windows       LDFLAGS: -L${SRCDIR}/../windows/bin -llibvorbis
// #include "codec.h"
import "C"

// VersionString returns a string giving version information for libvorbis
func VersionString() string {

	cstr := C.vorbis_version_string()
	return C.GoString(cstr)
}
