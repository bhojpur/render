package physics

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

type Material struct {
	name        string
	friction    float32
	restitution float32
}

type ContactMaterial struct {
	mat1                       *Material
	mat2                       *Material
	friction                   float32
	restitution                float32
	contactEquationStiffness   float32
	contactEquationRelaxation  float32
	frictionEquationStiffness  float32
	frictionEquationRelaxation float32
}

func NewContactMaterial() *ContactMaterial {

	cm := new(ContactMaterial)
	cm.friction = 0.3
	cm.restitution = 0.3
	cm.contactEquationStiffness = 1e7
	cm.contactEquationRelaxation = 3
	cm.frictionEquationStiffness = 1e7
	cm.frictionEquationRelaxation = 3
	return cm
}

//type intPair struct {
//	i int
//	j int
//}

//type ContactMaterialTable map[intPair]*ContactMaterial
