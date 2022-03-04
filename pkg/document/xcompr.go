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
	"compress/zlib"
	"fmt"
	"sync"
)

var xmem = xmempool{
	Pool: sync.Pool{
		New: func() interface{} {
			var m membuffer
			return &m
		},
	},
}

type xmempool struct{ sync.Pool }

func (pool *xmempool) compress(data []byte) *membuffer {
	mem := pool.Get().(*membuffer)
	buf := &mem.buf
	buf.Grow(len(data))

	zw, err := zlib.NewWriterLevel(buf, zlib.BestSpeed)
	if err != nil {
		panic(fmt.Errorf("could not create zlib writer: %w", err))
	}
	_, err = zw.Write(data)
	if err != nil {
		panic(fmt.Errorf("could not zlib-compress slice: %w", err))
	}

	err = zw.Close()
	if err != nil {
		panic(fmt.Errorf("could not close zlib writer: %w", err))
	}
	return mem
}

func (pool *xmempool) uncompress(data []byte) (*membuffer, error) {
	zr, err := zlib.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	mem := pool.Get().(*membuffer)
	mem.buf.Reset()

	_, err = mem.buf.ReadFrom(zr)
	if err != nil {
		mem.release()
		return nil, err
	}

	return mem, nil
}

type membuffer struct {
	buf bytes.Buffer
}

func (mem *membuffer) bytes() []byte { return mem.buf.Bytes() }
func (mem *membuffer) release() {
	mem.buf.Reset()
	xmem.Put(mem)
}

func (mem *membuffer) copy() []byte {
	src := mem.bytes()
	dst := make([]byte, len(src))
	copy(dst, src)
	return dst
}
