package collada

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
	"encoding/xml"
	"fmt"
	"io"
)

//
// LibraryMaterials
//
type LibraryMaterials struct {
	Id       string
	Name     string
	Asset    *Asset
	Material []*Material
}

// Dump prints out information about the LibraryMaterials
func (lm *LibraryMaterials) Dump(out io.Writer, indent int) {

	if lm == nil {
		return
	}
	fmt.Fprintf(out, "%sLibraryMaterials id:%s name:%s\n", sIndent(indent), lm.Id, lm.Name)
	for _, mat := range lm.Material {
		mat.Dump(out, indent+step)
	}
}

//
// Material
//
type Material struct {
	Id             string
	Name           string
	Asset          *Asset
	InstanceEffect InstanceEffect
}

// Dump prints out information about the Material
func (mat *Material) Dump(out io.Writer, indent int) {

	fmt.Fprintf(out, "%sMaterial id:%s name:%s\n", sIndent(indent), mat.Id, mat.Name)
	ind := indent + step
	mat.InstanceEffect.Dump(out, ind)
}

//
// InstanceEffect
//
type InstanceEffect struct {
	Sid  string
	Name string
	Url  string
}

// Dump prints out information about the InstanceEffect
func (ie *InstanceEffect) Dump(out io.Writer, indent int) {

	fmt.Fprintf(out, "%sInstanceEffect id:%s name:%s url:%s\n",
		sIndent(indent), ie.Sid, ie.Name, ie.Url)
}

func (d *Decoder) decLibraryMaterials(start xml.StartElement, dom *Collada) error {

	lm := new(LibraryMaterials)
	dom.LibraryMaterials = lm
	lm.Id = findAttrib(start, "id").Value
	lm.Name = findAttrib(start, "name").Value

	for {
		// Get next child element
		child, _, err := d.decNextChild(start)
		if err != nil || child.Name.Local == "" {
			return err
		}
		// Decodes <material>
		if child.Name.Local == "material" {
			err := d.decMaterial(child, lm)
			if err != nil {
				return err
			}
			continue
		}
	}
}

func (d *Decoder) decMaterial(start xml.StartElement, lm *LibraryMaterials) error {

	mat := new(Material)
	mat.Id = findAttrib(start, "id").Value
	mat.Name = findAttrib(start, "name").Value
	lm.Material = append(lm.Material, mat)

	for {
		child, _, err := d.decNextChild(start)
		if err != nil || child.Name.Local == "" {
			return err
		}
		if child.Name.Local == "instance_effect" {
			err := d.decInstanceEffect(child, &mat.InstanceEffect)
			if err != nil {
				return err
			}
			continue
		}
	}
}

func (d *Decoder) decInstanceEffect(start xml.StartElement, ie *InstanceEffect) error {

	ie.Sid = findAttrib(start, "sid").Value
	ie.Name = findAttrib(start, "name").Value
	ie.Url = findAttrib(start, "url").Value

	for {
		child, _, err := d.decNextChild(start)
		if err != nil || child.Name.Local == "" {
			return err
		}
		if child.Name.Local == "technique_hint setparam" {
			log.Warn("<technique_hint> not implemented")
			continue
		}
		if child.Name.Local == "setparam" {
			log.Warn("<setparam> not implemented")
			continue
		}
	}
}
