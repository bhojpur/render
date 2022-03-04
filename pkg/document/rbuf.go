package document

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
	"bytes"
	"encoding/binary"
	"io"
)

type rbuffer struct {
	p []byte
	c int
}

// newRBuffer returns a new buffer populated with the contents of the specified Reader
func newRBuffer(r io.Reader) (b *rbuffer, err error) {
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(r)
	b = &rbuffer{p: buf.Bytes()}
	return
}

func (r *rbuffer) Read(p []byte) (int, error) {
	if r.c >= len(r.p) {
		return 0, io.EOF
	}
	n := copy(p, r.p[r.c:])
	r.c += n
	return n, nil
}

func (r *rbuffer) ReadByte() (byte, error) {
	if r.c >= len(r.p) {
		return 0, io.EOF
	}
	v := r.p[r.c]
	r.c++
	return v, nil
}

func (r *rbuffer) u8() uint8 {
	if r.c >= len(r.p) {
		panic(io.ErrShortBuffer)
	}
	v := r.p[r.c]
	r.c++
	return v
}

func (r *rbuffer) u32() uint32 {
	const n = 4
	if r.c+n >= len(r.p) {
		panic(io.ErrShortBuffer)
	}
	beg := r.c
	r.c += n
	v := binary.BigEndian.Uint32(r.p[beg:])
	return v
}

func (r *rbuffer) i32() int32 {
	return int32(r.u32())
}

func (r *rbuffer) Next(n int) []byte {
	c := r.c
	r.c += n
	return r.p[c:r.c]
}
