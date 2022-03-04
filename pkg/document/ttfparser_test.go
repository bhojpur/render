package document_test

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
	"fmt"

	"github.com/bhojpur/render/internal/example"
	bdf "github.com/bhojpur/render/pkg/document"
)

func ExampleTtfParse() {
	ttf, err := bdf.TtfParse(example.FontDir() + "/calligra.ttf")
	if err == nil {
		fmt.Printf("Postscript name:  %s\n", ttf.PostScriptName)
		fmt.Printf("unitsPerEm:       %8d\n", ttf.UnitsPerEm)
		fmt.Printf("Xmin:             %8d\n", ttf.Xmin)
		fmt.Printf("Ymin:             %8d\n", ttf.Ymin)
		fmt.Printf("Xmax:             %8d\n", ttf.Xmax)
		fmt.Printf("Ymax:             %8d\n", ttf.Ymax)
	} else {
		fmt.Printf("%s\n", err)
	}
	// Output:
	// Postscript name:  CalligrapherRegular
	// unitsPerEm:           1000
	// Xmin:                 -173
	// Ymin:                 -234
	// Xmax:                 1328
	// Ymax:                  899
}

func hexStr(s string) string {
	var b bytes.Buffer
	b.WriteString("\"")
	for _, ch := range []byte(s) {
		b.WriteString(fmt.Sprintf("\\x%02x", ch))
	}
	b.WriteString("\":")
	return b.String()
}

func ExampleBdf_GetStringWidth() {
	pdf := bdf.New("", "", "", example.FontDir())
	pdf.SetFont("Helvetica", "", 12)
	pdf.AddPage()
	for _, s := range []string{"Hello", "世界", "\xe7a va?"} {
		fmt.Printf("%-32s width %5.2f, bytes %2d, runes %2d\n",
			hexStr(s), pdf.GetStringWidth(s), len(s), len([]rune(s)))
		if pdf.Err() {
			fmt.Println(pdf.Error())
		}
	}
	pdf.Close()
	// Output:
	// "\x48\x65\x6c\x6c\x6f":          width  9.64, bytes  5, runes  5
	// "\xe4\xb8\x96\xe7\x95\x8c":      width 13.95, bytes  6, runes  2
	// "\xe7\x61\x20\x76\x61\x3f":      width 12.47, bytes  6, runes  6
}
