package base

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
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"
	"testing"
)

var (
	flatteningThreshold = 0.5
	testsCubicFloat64   = []float64{
		100, 100, 200, 100, 100, 200, 200, 200,
		100, 100, 300, 200, 200, 200, 300, 100,
		100, 100, 0, 300, 200, 0, 300, 300,
		150, 290, 10, 10, 290, 10, 150, 290,
		10, 290, 10, 10, 290, 10, 290, 290,
		100, 290, 290, 10, 10, 10, 200, 290,
	}
	testsQuadFloat64 = []float64{
		100, 100, 200, 100, 200, 200,
		100, 100, 290, 200, 290, 100,
		100, 100, 0, 290, 200, 290,
		150, 290, 10, 10, 290, 290,
		10, 290, 10, 10, 290, 290,
		100, 290, 290, 10, 120, 290,
	}
)

func init() {
	os.Mkdir("test_results", 0666)
	f, err := os.Create("../output/curve/_test.html")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer f.Close()
	log.Printf("Create html viewer")
	f.Write([]byte("<html><body>"))
	for i := 0; i < len(testsCubicFloat64)/8; i++ {
		f.Write([]byte(fmt.Sprintf("<div><img src='_test%d.png'/></div>\n", i)))
	}
	for i := 0; i < len(testsQuadFloat64); i++ {
		f.Write([]byte(fmt.Sprintf("<div><img src='_testQuad%d.png'/>\n</div>\n", i)))
	}
	f.Write([]byte("</body></html>"))

}

func drawPoints(img draw.Image, c color.Color, s ...float64) image.Image {
	for i := 0; i < len(s); i += 2 {
		x, y := int(s[i]+0.5), int(s[i+1]+0.5)
		img.Set(x, y, c)
		img.Set(x, y+1, c)
		img.Set(x, y-1, c)
		img.Set(x+1, y, c)
		img.Set(x+1, y+1, c)
		img.Set(x+1, y-1, c)
		img.Set(x-1, y, c)
		img.Set(x-1, y+1, c)
		img.Set(x-1, y-1, c)

	}
	return img
}

func TestCubicCurve(t *testing.T) {
	for i := 0; i < len(testsCubicFloat64); i += 8 {
		var p SegmentedPath
		p.MoveTo(testsCubicFloat64[i], testsCubicFloat64[i+1])
		TraceCubic(&p, testsCubicFloat64[i:], flatteningThreshold)
		img := image.NewNRGBA(image.Rect(0, 0, 300, 300))
		PolylineBresenham(img, color.NRGBA{0xff, 0, 0, 0xff}, testsCubicFloat64[i:i+8]...)
		PolylineBresenham(img, image.Black, p.Points...)
		//drawPoints(img, image.NRGBAColor{0, 0, 0, 0xff}, curve[:]...)
		drawPoints(img, color.NRGBA{0, 0, 0, 0xff}, p.Points...)
		SaveToPngFile(fmt.Sprintf("../output/curve/_test%d.png", i/8), img)
		log.Printf("Num of points: %d\n", len(p.Points))
	}
	fmt.Println()
}

func TestQuadCurve(t *testing.T) {
	for i := 0; i < len(testsQuadFloat64); i += 6 {
		var p SegmentedPath
		p.MoveTo(testsQuadFloat64[i], testsQuadFloat64[i+1])
		TraceQuad(&p, testsQuadFloat64[i:], flatteningThreshold)
		img := image.NewNRGBA(image.Rect(0, 0, 300, 300))
		PolylineBresenham(img, color.NRGBA{0xff, 0, 0, 0xff}, testsQuadFloat64[i:i+6]...)
		PolylineBresenham(img, image.Black, p.Points...)
		//drawPoints(img, image.NRGBAColor{0, 0, 0, 0xff}, curve[:]...)
		drawPoints(img, color.NRGBA{0, 0, 0, 0xff}, p.Points...)
		SaveToPngFile(fmt.Sprintf("../output/curve/_testQuad%d.png", i), img)
		log.Printf("Num of points: %d\n", len(p.Points))
	}
	fmt.Println()
}

func TestQuadCurveCombinedPoint(t *testing.T) {
	var p1 SegmentedPath
	TraceQuad(&p1, []float64{0, 0, 0, 0, 0, 0}, flatteningThreshold)
	if len(p1.Points) != 2 {
		t.Error("It must have one point for this curve", len(p1.Points))
	}
	var p2 SegmentedPath
	TraceQuad(&p2, []float64{0, 0, 100, 100, 0, 0}, flatteningThreshold)
	if len(p2.Points) != 2 {
		t.Error("It must have one point for this curve", len(p2.Points))
	}
}

func TestCubicCurveCombinedPoint(t *testing.T) {
	var p1 SegmentedPath
	TraceCubic(&p1, []float64{0, 0, 0, 0, 0, 0, 0, 0}, flatteningThreshold)
	if len(p1.Points) != 2 {
		t.Error("It must have one point for this curve", len(p1.Points))
	}
	var p2 SegmentedPath
	TraceCubic(&p2, []float64{0, 0, 100, 100, 200, 200, 0, 0}, flatteningThreshold)
	if len(p2.Points) != 2 {
		t.Error("It must have one point for this curve", len(p2.Points))
	}
}

func BenchmarkCubicCurve(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for i := 0; i < len(testsCubicFloat64); i += 8 {
			var p SegmentedPath
			p.MoveTo(testsCubicFloat64[i], testsCubicFloat64[i+1])
			TraceCubic(&p, testsCubicFloat64[i:], flatteningThreshold)
		}
	}
}

// SaveToPngFile create and save an image to a file using PNG format
func SaveToPngFile(filePath string, m image.Image) error {
	// Create the file
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	// Create Writer from file
	b := bufio.NewWriter(f)
	// Write the image into the buffer
	err = png.Encode(b, m)
	if err != nil {
		return err
	}
	err = b.Flush()
	if err != nil {
		return err
	}
	return nil
}

func TestOutOfRangeTraceCurve(t *testing.T) {
	c := []float64{
		100, 100, 200, 100, 100, 200,
	}
	var p SegmentedPath
	TraceCubic(&p, c, flatteningThreshold)
}

func TestOutOfRangeTraceQuad(t *testing.T) {
	c := []float64{
		100, 100, 200, 100,
	}
	var p SegmentedPath
	TraceQuad(&p, c, flatteningThreshold)
}
