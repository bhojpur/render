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

// Vector2i is a 2D vector/point with X and Y integer32 components.
type Vector2i struct {
	X int32
	Y int32
}

// NewVector2i returns a new Vector2i with the specified x and y components.
func NewVector2i(x, y int32) Vector2i {
	return Vector2i{X: x, Y: y}
}

// NewVector2iScalar returns a new Vector2i with all components set to scalar.
func NewVector2iScalar(s int32) Vector2i {
	return Vector2i{X: s, Y: s}
}

// IsNil returns true if all values are 0 (uninitialized).
func (v Vector2i) IsNil() bool {
	if v.X == 0 && v.Y == 0 {
		return true
	}
	return false
}

// Set sets this vector X and Y components.
func (v *Vector2i) Set(x, y int32) {
	v.X = x
	v.Y = y
}

// SetScalar sets all vector components to same scalar value.
func (v *Vector2i) SetScalar(s int32) {
	v.X = s
	v.Y = s
}

// SetFromVector2 sets from a Vector2 (float32) vector.
func (v *Vector2i) SetFromVector2(vf Vector2) {
	v.X = int32(vf.X)
	v.Y = int32(vf.Y)
}

// SetDim sets this vector component value by its dimension index
func (v *Vector2i) SetDim(dim Dims, value int32) {
	switch dim {
	case X:
		v.X = value
	case Y:
		v.Y = value
	default:
		panic("dim is out of range")
	}
}

// Dim returns this vector component
func (v Vector2i) Dim(dim Dims) int32 {
	switch dim {
	case X:
		return v.X
	case Y:
		return v.Y
	default:
		panic("dim is out of range")
	}
}

// SetByName sets this vector component value by its case insensitive name: "x" or "y".
func (v *Vector2i) SetByName(name string, value int32) {
	switch name {
	case "x", "X":
		v.X = value
	case "y", "Y":
		v.Y = value
	default:
		panic("Invalid Vector2i component name: " + name)
	}
}

// SetZero sets this vector X and Y components to be zero.
func (v *Vector2i) SetZero() {
	v.SetScalar(0)
}

// FromArray sets this vector's components from the specified array and offset.
func (v *Vector2i) FromArray(array []int32, offset int) {
	v.X = array[offset]
	v.Y = array[offset+1]
}

// ToArray copies this vector's components to array starting at offset.
func (v Vector2i) ToArray(array []int32, offset int) {
	array[offset] = v.X
	array[offset+1] = v.Y
}

///////////////////////////////////////////////////////////////////////
//  Basic math operations

// Add adds other vector to this one and returns result in a new vector.
func (v Vector2i) Add(other Vector2i) Vector2i {
	return Vector2i{v.X + other.X, v.Y + other.Y}
}

// AddScalar adds scalar s to each component of this vector and returns new vector.
func (v Vector2i) AddScalar(s int32) Vector2i {
	return Vector2i{v.X + s, v.Y + s}
}

// SetAdd sets this to addition with other vector (i.e., += or plus-equals).
func (v *Vector2i) SetAdd(other Vector2i) {
	v.X += other.X
	v.Y += other.Y
}

// SetAddScalar sets this to addition with scalar.
func (v *Vector2i) SetAddScalar(s int32) {
	v.X += s
	v.Y += s
}

// Sub subtracts other vector from this one and returns result in new vector.
func (v Vector2i) Sub(other Vector2i) Vector2i {
	return Vector2i{v.X - other.X, v.Y - other.Y}
}

// SubScalar subtracts scalar s from each component of this vector and returns new vector.
func (v Vector2i) SubScalar(s int32) Vector2i {
	return Vector2i{v.X - s, v.Y - s}
}

// SetSub sets this to subtraction with other vector (i.e., -= or minus-equals).
func (v *Vector2i) SetSub(other Vector2i) {
	v.X -= other.X
	v.Y -= other.Y
}

// SetSubScalar sets this to subtraction of scalar.
func (v *Vector2i) SetSubScalar(s int32) {
	v.X -= s
	v.Y -= s
}

// Mul multiplies each component of this vector by the corresponding one from other
// and returns resulting vector.
func (v Vector2i) Mul(other Vector2i) Vector2i {
	return Vector2i{v.X * other.X, v.Y * other.Y}
}

// MulScalar multiplies each component of this vector by the scalar s and returns resulting vector.
func (v Vector2i) MulScalar(s int32) Vector2i {
	return Vector2i{v.X * s, v.Y * s}
}

// SetMul sets this to multiplication with other vector (i.e., *= or times-equals).
func (v *Vector2i) SetMul(other Vector2i) {
	v.X *= other.X
	v.Y *= other.Y
}

// SetMulScalar sets this to multiplication by scalar.
func (v *Vector2i) SetMulScalar(s int32) {
	v.X *= s
	v.Y *= s
}

// Div divides each component of this vector by the corresponding one from other vector
// and returns resulting vector.
func (v Vector2i) Div(other Vector2i) Vector2i {
	return Vector2i{v.X / other.X, v.Y / other.Y}
}

// DivScalar divides each component of this vector by the scalar s and returns resulting vector.
// If scalar is zero, returns zero.
func (v Vector2i) DivScalar(scalar int32) Vector2i {
	if scalar != 0 {
		return v.MulScalar(1 / scalar)
	} else {
		return Vector2i{}
	}
}

// SetDiv sets this to division by other vector (i.e., /= or divide-equals).
func (v *Vector2i) SetDiv(other Vector2i) {
	v.X /= other.X
	v.Y /= other.Y
}

// SetDivScalar sets this to division by scalar.
func (v *Vector2i) SetDivScalar(s int32) {
	if s != 0 {
		v.SetMulScalar(1 / s)
	} else {
		v.SetZero()
	}
}

// Min32i returns the min of the two int32 numbers.
func Min32i(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

// Max32i returns the max of the two int32 numbers.
func Max32i(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}

// Min returns min of this vector components vs. other vector.
func (v Vector2i) Min(other Vector2i) Vector2i {
	return Vector2i{Min32i(v.X, other.X), Min32i(v.Y, other.Y)}
}

// SetMin sets this vector components to the minimum values of itself and other vector.
func (v *Vector2i) SetMin(other Vector2i) {
	v.X = Min32i(v.X, other.X)
	v.Y = Min32i(v.Y, other.Y)
}

// Max returns max of this vector components vs. other vector.
func (v Vector2i) Max(other Vector2i) Vector2i {
	return Vector2i{Max32i(v.X, other.X), Max32i(v.Y, other.Y)}
}

// SetMax sets this vector components to the maximum value of itself and other vector.
func (v *Vector2i) SetMax(other Vector2i) {
	v.X = Max32i(v.X, other.X)
	v.Y = Max32i(v.Y, other.Y)
}

// Clamp sets this vector components to be no less than the corresponding components of min
// and not greater than the corresponding component of max.
// Assumes min < max, if this assumption isn't true it will not operate correctly.
func (v *Vector2i) Clamp(min, max Vector2i) {
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
}

// ClampScalar sets this vector components to be no less than minVal and not greater than maxVal.
func (v *Vector2i) ClampScalar(minVal, maxVal int32) {
	v.Clamp(NewVector2iScalar(minVal), NewVector2iScalar(maxVal))
}

// Negate returns vector with each component negated.
func (v Vector2i) Negate() Vector2i {
	return Vector2i{-v.X, -v.Y}
}

// SetNegate negates each of this vector's components.
func (v *Vector2i) SetNegate() {
	v.X = -v.X
	v.Y = -v.Y
}

// IsEqual returns if this vector is equal to other.
func (v Vector2i) IsEqual(other Vector2i) bool {
	return (other.X == v.X) && (other.Y == v.Y)
}
