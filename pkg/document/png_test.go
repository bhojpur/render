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
	"os"
	"testing"
)

func BenchmarkParsePNG_rgb(b *testing.B) {
	raw, err := os.ReadFile("image/golang-gopher.png")
	if err != nil {
		b.Fatal(err)
	}

	pdf := New("P", "mm", "A4", "")
	pdf.AddPage()

	const readDPI = true
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = pdf.parsepng(bytes.NewReader(raw), readDPI)
	}
}

func BenchmarkParsePNG_gray(b *testing.B) {
	raw, err := os.ReadFile("image/logo-gray.png")
	if err != nil {
		b.Fatal(err)
	}

	pdf := New("P", "mm", "A4", "")
	pdf.AddPage()

	const readDPI = true
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = pdf.parsepng(bytes.NewReader(raw), readDPI)
	}
}

func BenchmarkParsePNG_small(b *testing.B) {
	raw, err := os.ReadFile("image/logo.png")
	if err != nil {
		b.Fatal(err)
	}

	pdf := New("P", "mm", "A4", "")
	pdf.AddPage()

	const readDPI = true
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = pdf.parsepng(bytes.NewReader(raw), readDPI)
	}
}

func BenchmarkParseJPG(b *testing.B) {
	raw, err := os.ReadFile("image/logo_bdf.jpg")
	if err != nil {
		b.Fatal(err)
	}

	pdf := New("P", "mm", "A4", "")
	pdf.AddPage()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = pdf.parsejpg(bytes.NewReader(raw))
	}
}

func BenchmarkParseGIF(b *testing.B) {
	raw, err := os.ReadFile("image/logo.gif")
	if err != nil {
		b.Fatal(err)
	}

	pdf := New("P", "mm", "A4", "")
	pdf.AddPage()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = pdf.parsegif(bytes.NewReader(raw))
	}
}
