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
	"github.com/bhojpur/render/pkg/experimental/physics/equation"
)

// ISolver is the interface type for all constraint solvers.
type ISolver interface {
	Solve(dt float32, nBodies int) *Solution // Solve should return the number of iterations performed
	AddEquation(eq equation.IEquation)
	RemoveEquation(eq equation.IEquation) bool
	ClearEquations()
}

// Solution represents a solver solution
type Solution struct {
	VelocityDeltas        []math32.Vector3
	AngularVelocityDeltas []math32.Vector3
	Iterations            int
}

// Constraint equation solver base class.
type Solver struct {
	Solution
	equations []equation.IEquation // All equations to be solved
}

// AddEquation adds an equation to the solver.
func (s *Solver) AddEquation(eq equation.IEquation) {

	s.equations = append(s.equations, eq)
}

// RemoveEquation removes the specified equation from the solver.
// Returns true if found, false otherwise.
func (s *Solver) RemoveEquation(eq equation.IEquation) bool {

	for pos, current := range s.equations {
		if current == eq {
			copy(s.equations[pos:], s.equations[pos+1:])
			s.equations[len(s.equations)-1] = nil
			s.equations = s.equations[:len(s.equations)-1]
			return true
		}
	}
	return false
}

// ClearEquations removes all equations from the solver.
func (s *Solver) ClearEquations() {

	s.equations = s.equations[0:0]
}
