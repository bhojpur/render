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
	"strconv"

	"github.com/bhojpur/render/pkg/engine/core"
	"github.com/bhojpur/render/pkg/engine/math32"
	"github.com/bhojpur/render/pkg/gls"
)

// MaxBoneInfluencers is the maximum number of bone influencers per vertex.
const MaxBoneInfluencers = 4

// RiggedMesh is a Mesh associated with a skeleton.
type RiggedMesh struct {
	*Mesh    // Embedded mesh
	skeleton *Skeleton
	mBones   gls.Uniform
}

// NewRiggedMesh returns a new rigged mesh.
func NewRiggedMesh(mesh *Mesh) *RiggedMesh {

	rm := new(RiggedMesh)
	rm.Mesh = mesh
	rm.SetIGraphic(rm)
	rm.mBones.Init("mBones")
	rm.ShaderDefines.Set("BONE_INFLUENCERS", strconv.Itoa(MaxBoneInfluencers))
	rm.ShaderDefines.Set("TOTAL_BONES", "0")

	return rm
}

// SetSkeleton sets the skeleton used by the rigged mesh.
func (rm *RiggedMesh) SetSkeleton(sk *Skeleton) {

	rm.skeleton = sk
	rm.ShaderDefines.Set("TOTAL_BONES", strconv.Itoa(len(rm.skeleton.Bones())))
}

// SetSkeleton returns the skeleton used by the rigged mesh.
func (rm *RiggedMesh) Skeleton() *Skeleton {

	return rm.skeleton
}

// RenderSetup is called by the renderer before drawing the geometry.
func (rm *RiggedMesh) RenderSetup(gs *gls.GLS, rinfo *core.RenderInfo) {

	// Call base mesh's RenderSetup
	rm.Mesh.RenderSetup(gs, rinfo)

	// Get inverse matrix world
	var invMat math32.Matrix4
	node := rm.GetNode()
	nMW := node.MatrixWorld()
	err := invMat.GetInverse(&nMW)
	if err != nil {
		log.Error("Skeleton.BoneMatrices: inverting matrix failed!")
	}

	// Transfer bone matrices
	boneMatrices := rm.skeleton.BoneMatrices(&invMat)
	location := rm.mBones.Location(gs)
	gs.UniformMatrix4fv(location, int32(len(boneMatrices)), false, &boneMatrices[0][0])
}
