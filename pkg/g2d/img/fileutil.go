package img

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
	"image"
	"image/png"
	"os"
)

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

// LoadFromPngFile Open a png file
func LoadFromPngFile(filePath string) (image.Image, error) {
	// Open file
	f, err := os.OpenFile(filePath, 0, 0)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	b := bufio.NewReader(f)
	img, err := png.Decode(b)
	if err != nil {
		return nil, err
	}
	return img, nil
}
