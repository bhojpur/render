package solver

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

// It implements a basic physics engine.

import (
	"github.com/bhojpur/render/pkg/engine/math32"
)

// GaussSeidel equation solver.
// See https://en.wikipedia.org/wiki/Gauss-Seidel_method.
// The number of solver iterations determines the quality of the solution.
// More iterations yield a better solution but require more computation.
type GaussSeidel struct {
	Solver
	maxIter   int     // Number of solver iterations.
	tolerance float32 // When the error is less than the tolerance, the system is assumed to be converged.

	solveInvCs  []float32
	solveBs     []float32
	solveLambda []float32
}

// NewGaussSeidel creates and returns a pointer to a new GaussSeidel constraint equation solver.
func NewGaussSeidel() *GaussSeidel {

	gs := new(GaussSeidel)
	gs.maxIter = 20
	gs.tolerance = 1e-7

	gs.VelocityDeltas = make([]math32.Vector3, 0)
	gs.AngularVelocityDeltas = make([]math32.Vector3, 0)

	gs.solveInvCs = make([]float32, 0)
	gs.solveBs = make([]float32, 0)
	gs.solveLambda = make([]float32, 0)

	return gs
}

func (gs *GaussSeidel) Reset(numBodies int) {

	// Reset solution
	gs.VelocityDeltas = make([]math32.Vector3, numBodies)
	gs.AngularVelocityDeltas = make([]math32.Vector3, numBodies)
	gs.Iterations = 0

	// Reset internal arrays
	gs.solveInvCs = gs.solveInvCs[0:0]
	gs.solveBs = gs.solveBs[0:0]
	gs.solveLambda = gs.solveLambda[0:0]
}

// Solve
func (gs *GaussSeidel) Solve(dt float32, nBodies int) *Solution {

	gs.Reset(nBodies)

	iter := 0
	nEquations := len(gs.equations)
	h := dt

	// Things that do not change during iteration can be computed once
	for i := 0; i < nEquations; i++ {
		eq := gs.equations[i]
		gs.solveInvCs = append(gs.solveInvCs, 1.0/eq.ComputeC())
		gs.solveBs = append(gs.solveBs, eq.ComputeB(h))
		gs.solveLambda = append(gs.solveLambda, 0.0)
	}

	if nEquations > 0 {
		tolSquared := gs.tolerance * gs.tolerance

		// Iterate over equations
		for iter = 0; iter < gs.maxIter; iter++ {

			// Accumulate the total error for each iteration.
			deltaLambdaTot := float32(0)

			for j := 0; j < nEquations; j++ {
				eq := gs.equations[j]

				// Compute iteration
				lambdaJ := gs.solveLambda[j]

				idxBodyA := eq.BodyA().Index()
				idxBodyB := eq.BodyB().Index()

				vA := gs.VelocityDeltas[idxBodyA]
				vB := gs.VelocityDeltas[idxBodyB]
				wA := gs.AngularVelocityDeltas[idxBodyA]
				wB := gs.AngularVelocityDeltas[idxBodyB]
				jeA := eq.JeA()
				jeB := eq.JeB()
				GWlambda := jeA.MultiplyVectors(&vA, &wA) + jeB.MultiplyVectors(&vB, &wB)

				deltaLambda := gs.solveInvCs[j] * (gs.solveBs[j] - GWlambda - eq.Eps()*lambdaJ)

				// Clamp if we are outside the min/max interval
				if (lambdaJ + deltaLambda) < eq.MinForce() {
					deltaLambda = eq.MinForce() - lambdaJ
				} else if (lambdaJ + deltaLambda) > eq.MaxForce() {
					deltaLambda = eq.MaxForce() - lambdaJ
				}
				gs.solveLambda[j] += deltaLambda
				deltaLambdaTot += math32.Abs(deltaLambda)

				// Add to velocity deltas
				spatA := jeA.Spatial()
				spatB := jeB.Spatial()
				gs.VelocityDeltas[idxBodyA].Add(spatA.MultiplyScalar(eq.BodyA().InvMassEff() * deltaLambda))
				gs.VelocityDeltas[idxBodyB].Add(spatB.MultiplyScalar(eq.BodyB().InvMassEff() * deltaLambda))

				// Add to angular velocity deltas
				rotA := jeA.Rotational()
				rotB := jeB.Rotational()
				gs.AngularVelocityDeltas[idxBodyA].Add(rotA.ApplyMatrix3(eq.BodyA().InvRotInertiaWorldEff()).MultiplyScalar(deltaLambda))
				gs.AngularVelocityDeltas[idxBodyB].Add(rotB.ApplyMatrix3(eq.BodyB().InvRotInertiaWorldEff()).MultiplyScalar(deltaLambda))
			}

			// If the total error is small enough - stop iterating
			if deltaLambdaTot*deltaLambdaTot < tolSquared {
				break
			}
		}

		// Set the .multiplier property of each equation
		for i := range gs.equations {
			gs.equations[i].SetMultiplier(gs.solveLambda[i] / h)
		}
		iter += 1
	}

	gs.Iterations = iter

	return &gs.Solution
}
