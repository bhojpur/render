package graphic

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
	"github.com/bhojpur/render/pkg/3d/core"
	"github.com/bhojpur/render/pkg/3d/math32"
)

// Skeleton contains armature information.
type Skeleton struct {
	inverseBindMatrices []math32.Matrix4
	boneMatrices        []math32.Matrix4
	bones               []*core.Node
}

// NewSkeleton creates and returns a pointer to a new Skeleton.
func NewSkeleton() *Skeleton {

	sk := new(Skeleton)
	sk.boneMatrices = make([]math32.Matrix4, 0)
	sk.bones = make([]*core.Node, 0)
	return sk
}

// AddBone adds a bone to the skeleton along with an optional inverseBindMatrix.
func (sk *Skeleton) AddBone(node *core.Node, inverseBindMatrix *math32.Matrix4) {

	// Useful for debugging:
	//node.Add(NewAxisHelper(0.2))

	sk.bones = append(sk.bones, node)
	sk.boneMatrices = append(sk.boneMatrices, *math32.NewMatrix4())
	if inverseBindMatrix == nil {
		inverseBindMatrix = math32.NewMatrix4() // Identity matrix
	}

	sk.inverseBindMatrices = append(sk.inverseBindMatrices, *inverseBindMatrix)
}

// Bones returns the list of bones in the skeleton.
func (sk *Skeleton) Bones() []*core.Node {

	return sk.bones
}

// BoneMatrices calculates and returns the bone world matrices to be sent to the shader.
func (sk *Skeleton) BoneMatrices(invMat *math32.Matrix4) []math32.Matrix4 {

	// Update bone matrices based on inverseBindMatrices and the provided invMat
	for i := range sk.bones {
		bMat := sk.bones[i].MatrixWorld()
		bMat.MultiplyMatrices(&bMat, &sk.inverseBindMatrices[i])
		sk.boneMatrices[i].MultiplyMatrices(invMat, &bMat)
	}

	return sk.boneMatrices
}
