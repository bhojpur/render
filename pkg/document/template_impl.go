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
	"crypto/sha1"
	"encoding/gob"
	"errors"
	"fmt"
)

// newTpl creates a template, copying graphics settings from a template if one is given
func newTpl(corner PointType, size SizeType, orientationStr, unitStr, fontDirStr string, fn func(*Tpl), copyFrom *Bdf) Template {
	sizeStr := ""

	bdf := bdfNew(orientationStr, unitStr, sizeStr, fontDirStr, size)
	tpl := Tpl{*bdf}
	if copyFrom != nil {
		tpl.loadParamsFromBdf(copyFrom)
	}
	tpl.Bdf.AddPage()
	fn(&tpl)

	bytes := make([][]byte, len(tpl.Bdf.pages))
	// skip the first page as it will always be empty
	for x := 1; x < len(bytes); x++ {
		bytes[x] = tpl.Bdf.pages[x].Bytes()
	}

	templates := make([]Template, 0, len(tpl.Bdf.templates))
	for _, key := range templateKeyList(tpl.Bdf.templates, true) {
		templates = append(templates, tpl.Bdf.templates[key])
	}
	images := tpl.Bdf.images

	template := BdfTpl{corner, size, bytes, images, templates, tpl.Bdf.page}
	return &template
}

// BdfTpl is a concrete implementation of the Template interface.
type BdfTpl struct {
	corner    PointType
	size      SizeType
	bytes     [][]byte
	images    map[string]*ImageInfoType
	templates []Template
	page      int
}

// ID returns the global template identifier
func (t *BdfTpl) ID() string {
	return fmt.Sprintf("%x", sha1.Sum(t.Bytes()))
}

// Size gives the bounding dimensions of this template
func (t *BdfTpl) Size() (corner PointType, size SizeType) {
	return t.corner, t.size
}

// Bytes returns the actual template data, not including resources
func (t *BdfTpl) Bytes() []byte {
	return t.bytes[t.page]
}

// FromPage creates a new template from a specific Page
func (t *BdfTpl) FromPage(page int) (Template, error) {
	// pages start at 1
	if page == 0 {
		return nil, errors.New("bdf: pages start at 1 No template will have a page 0")
	}

	if page > t.NumPages() {
		return nil, fmt.Errorf("bdf: the template does not have a page %d", page)
	}
	// if it is already pointing to the correct page
	// there is no need to create a new template
	if t.page == page {
		return t, nil
	}

	t2 := *t
	t2.page = page
	return &t2, nil
}

// FromPages creates a template slice with all the pages within a template.
func (t *BdfTpl) FromPages() []Template {
	p := make([]Template, t.NumPages())
	for x := 1; x <= t.NumPages(); x++ {
		// the only error is when accessing a
		// non existing template... that can't happen
		// here
		p[x-1], _ = t.FromPage(x)
	}

	return p
}

// Images returns a list of the images used in this template
func (t *BdfTpl) Images() map[string]*ImageInfoType {
	return t.images
}

// Templates returns a list of templates used in this template
func (t *BdfTpl) Templates() []Template {
	return t.templates
}

// NumPages returns the number of available pages within the template. Look at FromPage and FromPages on access to that content.
func (t *BdfTpl) NumPages() int {
	// the first page is empty to
	// make the pages begin at one
	return len(t.bytes) - 1
}

// Serialize turns a template into a byte string for later deserialization
func (t *BdfTpl) Serialize() ([]byte, error) {
	b := new(bytes.Buffer)
	enc := gob.NewEncoder(b)
	err := enc.Encode(t)

	return b.Bytes(), err
}

// DeserializeTemplate creaties a template from a previously serialized
// template
func DeserializeTemplate(b []byte) (Template, error) {
	tpl := new(BdfTpl)
	dec := gob.NewDecoder(bytes.NewBuffer(b))
	err := dec.Decode(tpl)
	return tpl, err
}

// childrenImages returns the next layer of children images, it doesn't dig into
// children of children. Applies template namespace to keys to ensure
// no collisions. See UseTemplateScaled
func (t *BdfTpl) childrenImages() map[string]*ImageInfoType {
	childrenImgs := make(map[string]*ImageInfoType)

	for x := 0; x < len(t.templates); x++ {
		imgs := t.templates[x].Images()
		for key, val := range imgs {
			name := sprintf("t%s-%s", t.templates[x].ID(), key)
			childrenImgs[name] = val
		}
	}

	return childrenImgs
}

// childrensTemplates returns the next layer of children templates, it doesn't dig into
// children of children.
func (t *BdfTpl) childrensTemplates() []Template {
	childrenTmpls := make([]Template, 0)

	for x := 0; x < len(t.templates); x++ {
		tmpls := t.templates[x].Templates()
		childrenTmpls = append(childrenTmpls, tmpls...)
	}

	return childrenTmpls
}

// GobEncode encodes the receiving template into a byte buffer. Use GobDecode
// to decode the byte buffer back to a template.
func (t *BdfTpl) GobEncode() ([]byte, error) {
	w := new(bytes.Buffer)
	encoder := gob.NewEncoder(w)

	childrensTemplates := t.childrensTemplates()
	firstClassTemplates := make([]Template, 0)

found_continue:
	for x := 0; x < len(t.templates); x++ {
		for y := 0; y < len(childrensTemplates); y++ {
			if childrensTemplates[y].ID() == t.templates[x].ID() {
				continue found_continue
			}
		}

		firstClassTemplates = append(firstClassTemplates, t.templates[x])
	}
	err := encoder.Encode(firstClassTemplates)

	childrenImgs := t.childrenImages()
	firstClassImgs := make(map[string]*ImageInfoType)

	for key, img := range t.images {
		if _, ok := childrenImgs[key]; !ok {
			firstClassImgs[key] = img
		}
	}

	if err == nil {
		err = encoder.Encode(firstClassImgs)
	}
	if err == nil {
		err = encoder.Encode(t.corner)
	}
	if err == nil {
		err = encoder.Encode(t.size)
	}
	if err == nil {
		err = encoder.Encode(t.bytes)
	}
	if err == nil {
		err = encoder.Encode(t.page)
	}

	return w.Bytes(), err
}

// GobDecode decodes the specified byte buffer into the receiving template.
func (t *BdfTpl) GobDecode(buf []byte) error {
	r := bytes.NewBuffer(buf)
	decoder := gob.NewDecoder(r)

	firstClassTemplates := make([]*BdfTpl, 0)
	err := decoder.Decode(&firstClassTemplates)
	t.templates = make([]Template, len(firstClassTemplates))

	for x := 0; x < len(t.templates); x++ {
		t.templates[x] = Template(firstClassTemplates[x])
	}

	firstClassImages := t.childrenImages()

	t.templates = append(t.childrensTemplates(), t.templates...)

	t.images = make(map[string]*ImageInfoType)
	if err == nil {
		err = decoder.Decode(&t.images)
	}

	for k, v := range firstClassImages {
		t.images[k] = v
	}

	if err == nil {
		err = decoder.Decode(&t.corner)
	}
	if err == nil {
		err = decoder.Decode(&t.size)
	}
	if err == nil {
		err = decoder.Decode(&t.bytes)
	}
	if err == nil {
		err = decoder.Decode(&t.page)
	}

	return err
}

// Tpl is an Bdf used for writing a template. It has most of the facilities of
// an Bdf, but cannot add more pages. Tpl is used directly only during the
// limited time a template is writable.
type Tpl struct {
	Bdf
}

func (t *Tpl) loadParamsFromBdf(f *Bdf) {
	t.Bdf.compress = false

	t.Bdf.k = f.k
	t.Bdf.x = f.x
	t.Bdf.y = f.y
	t.Bdf.lineWidth = f.lineWidth
	t.Bdf.capStyle = f.capStyle
	t.Bdf.joinStyle = f.joinStyle

	t.Bdf.color.draw = f.color.draw
	t.Bdf.color.fill = f.color.fill
	t.Bdf.color.text = f.color.text

	t.Bdf.fonts = f.fonts
	t.Bdf.currentFont = f.currentFont
	t.Bdf.fontFamily = f.fontFamily
	t.Bdf.fontSize = f.fontSize
	t.Bdf.fontSizePt = f.fontSizePt
	t.Bdf.fontStyle = f.fontStyle
	t.Bdf.ws = f.ws

	for key, value := range f.images {
		t.Bdf.images[key] = value
	}
}
