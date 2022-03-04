package math32

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

// Vector3i is a 3D vector/point with X, Y and Z int32 components.
type Vector3i struct {
	X int32
	Y int32
	Z int32
}

// NewVector3i returns a new Vector3i with the specified x, y and y components.
func NewVector3i(x, y, z int32) Vector3i {
	return Vector3i{X: x, Y: y, Z: z}
}

// NewVector3iScalar returns a new Vector3 with all components set to scalar.
func NewVector3iScalar(s int32) Vector3i {
	return Vector3i{X: s, Y: s, Z: s}
}

// IsNil returns true if all values are 0 (uninitialized).
func (v Vector3i) IsNil() bool {
	if v.X == 0 && v.Y == 0 && v.Z == 0 {
		return true
	}
	return false
}

// Set sets this vector X, Y and Z components.
func (v *Vector3i) Set(x, y, z int32) {
	v.X = x
	v.Y = y
	v.Z = z
}

// SetScalar sets all vector X, Y and Z components to same scalar value.
func (v *Vector3i) SetScalar(s int32) {
	v.X = s
	v.Y = s
	v.Z = s
}

// SetFromVector3 sets from a Vector3 (float32) vector.
func (v *Vector3i) SetFromVec3(vf Vector3) {
	v.X = int32(vf.X)
	v.Y = int32(vf.Y)
	v.Z = int32(vf.Z)
}

// SetDim sets this vector component value by dimension index.
func (v *Vector3i) SetDim(dim Dims, value int32) {
	switch dim {
	case X:
		v.X = value
	case Y:
		v.Y = value
	case Z:
		v.Z = value
	default:
		panic("dim is out of range: ")
	}
}

// Dim returns this vector component
func (v Vector3i) Dim(dim Dims) int32 {
	switch dim {
	case X:
		return v.X
	case Y:
		return v.Y
	case Z:
		return v.Z
	default:
		panic("dim is out of range")
	}
}

// SetByName sets this vector component value by its case insensitive name: "x", "y", or "z".
func (v *Vector3i) SetByName(name string, value int32) {
	switch name {
	case "x", "X":
		v.X = value
	case "y", "Y":
		v.Y = value
	case "z", "Z":
		v.Z = value
	default:
		panic("Invalid Vector3i component name: " + name)
	}
}

// SetZero sets this vector X, Y and Z components to be zero.
func (v *Vector3i) SetZero() {
	v.SetScalar(0)
}

// FromArray sets this vector's components from the specified array and offset
func (v *Vector3i) FromArray(array []int32, offset int) {
	v.X = array[offset]
	v.Y = array[offset+1]
	v.Z = array[offset+2]
}

// ToArray copies this vector's components to array starting at offset.
func (v Vector3i) ToArray(array []int32, offset int) {
	array[offset] = v.X
	array[offset+1] = v.Y
	array[offset+2] = v.Z
}

///////////////////////////////////////////////////////////////////////
//  Basic math operations

// Add adds other vector to this one and returns result in a new vector.
func (v Vector3i) Add(other Vector3i) Vector3i {
	return Vector3i{v.X + other.X, v.Y + other.Y, v.Z + other.Z}
}

// AddScalar adds scalar s to each component of this vector and returns new vector.
func (v Vector3i) AddScalar(s int32) Vector3i {
	return Vector3i{v.X + s, v.Y + s, v.Z + s}
}

// SetAdd sets this to addition with other vector (i.e., += or plus-equals).
func (v *Vector3i) SetAdd(other Vector3i) {
	v.X += other.X
	v.Y += other.Y
	v.Z += other.Z
}

// SetAddScalar sets this to addition with scalar.
func (v *Vector3i) SetAddScalar(s int32) {
	v.X += s
	v.Y += s
	v.Z += s
}

// Sub subtracts other vector from this one and returns result in new vector.
func (v Vector3i) Sub(other Vector3i) Vector3i {
	return Vector3i{v.X - other.X, v.Y - other.Y, v.Z - other.Z}
}

// SubScalar subtracts scalar s from each component of this vector and returns new vector.
func (v Vector3i) SubScalar(s int32) Vector3i {
	return Vector3i{v.X - s, v.Y - s, v.Z - s}
}

// SetSub sets this to subtraction with other vector (i.e., -= or minus-equals).
func (v *Vector3i) SetSub(other Vector3i) {
	v.X -= other.X
	v.Y -= other.Y
	v.Z -= other.Z
}

// SetSubScalar sets this to subtraction of scalar.
func (v *Vector3i) SetSubScalar(s int32) {
	v.X -= s
	v.Y -= s
	v.Z -= s
}

// Mul multiplies each component of this vector by the corresponding one from other
// and returns resulting vector.
func (v Vector3i) Mul(other Vector3i) Vector3i {
	return Vector3i{v.X * other.X, v.Y * other.Y, v.Z * other.Z}
}

// MulScalar multiplies each component of this vector by the scalar s and returns resulting vector.
func (v Vector3i) MulScalar(s int32) Vector3i {
	return Vector3i{v.X * s, v.Y * s, v.Z * s}
}

// SetMul sets this to multiplication with other vector (i.e., *= or times-equals).
func (v *Vector3i) SetMul(other Vector3i) {
	v.X *= other.X
	v.Y *= other.Y
	v.Z *= other.Z
}

// SetMulScalar sets this to multiplication by scalar.
func (v *Vector3i) SetMulScalar(s int32) {
	v.X *= s
	v.Y *= s
	v.Z *= s
}

// Div divides each component of this vector by the corresponding one from other vector
// and returns resulting vector.
func (v Vector3i) Div(other Vector3i) Vector3i {
	return Vector3i{v.X / other.X, v.Y / other.Y, v.Z / other.Z}
}

// DivScalar divides each component of this vector by the scalar s and returns resulting vector.
// If scalar is zero, returns zero.
func (v Vector3i) DivScalar(scalar int32) Vector3i {
	if scalar != 0 {
		return v.MulScalar(1 / scalar)
	} else {
		return Vector3i{}
	}
}

// SetDiv sets this to division by other vector (i.e., /= or divide-equals).
func (v *Vector3i) SetDiv(other Vector3i) {
	v.X /= other.X
	v.Y /= other.Y
	v.Z /= other.Z
}

// SetDivScalar sets this to division by scalar.
func (v *Vector3i) SetDivScalar(s int32) {
	if s != 0 {
		v.SetMulScalar(1 / s)
	} else {
		v.SetZero()
	}
}

// Min returns min of this vector components vs. other vector.
func (v Vector3i) Min(other Vector3i) Vector3i {
	return Vector3i{Min32i(v.X, other.X), Min32i(v.Y, other.Y), Min32i(v.Z, other.Z)}
}

// SetMin sets this vector components to the minimum values of itself and other vector.
func (v *Vector3i) SetMin(other Vector3i) {
	v.X = Min32i(v.X, other.X)
	v.Y = Min32i(v.Y, other.Y)
	v.Z = Min32i(v.Z, other.Z)
}

// Max returns max of this vector components vs. other vector.
func (v Vector3i) Max(other Vector3i) Vector3i {
	return Vector3i{Max32i(v.X, other.X), Max32i(v.Y, other.Y), Max32i(v.Z, other.Z)}
}

// SetMax sets this vector components to the maximum value of itself and other vector.
func (v *Vector3i) SetMax(other Vector3i) {
	v.X = Max32i(v.X, other.X)
	v.Y = Max32i(v.Y, other.Y)
	v.Z = Max32i(v.Z, other.Z)
}

// Clamp sets this vector components to be no less than the corresponding components of min
// and not greater than the corresponding component of max.
// Assumes min < max, if this assumption isn't true it will not operate correctly.
func (v *Vector3i) Clamp(min, max Vector3i) {
	if v.X < min.X {
		v.X = min.X
	} else if v.X > max.X {
		v.X = max.X
	}
	if v.Y < min.Y {
		v.Y = min.Y
	} else if v.Y > max.Y {
		v.Y = max.Y
	}
	if v.Z < min.Z {
		v.Z = min.Z
	} else if v.Z > max.Z {
		v.Z = max.Z
	}
}

// ClampScalar sets this vector components to be no less than minVal and not greater than maxVal.
func (v *Vector3i) ClampScalar(minVal, maxVal int32) {
	v.Clamp(NewVector3iScalar(minVal), NewVector3iScalar(maxVal))
}

// Negate returns vector with each component negated.
func (v Vector3i) Negate() Vector3i {
	return Vector3i{-v.X, -v.Y, -v.Z}
}

// SetNegate negates each of this vector's components.
func (v *Vector3i) SetNegate() {
	v.X = -v.X
	v.Y = -v.Y
	v.Z = -v.Z
}

// IsEqual returns if this vector is equal to other.
func (v Vector3i) IsEqual(other Vector3i) bool {
	return (other.X == v.X) && (other.Y == v.Y) && (other.Z == v.Z)
}
