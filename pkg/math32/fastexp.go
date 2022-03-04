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

import (
	"math"
)

// This is a fast version of the natural exponential function, for highly
// time-sensitive uses where precision is less important.

// Based on: N. N. Schraudolph. "A fast, compact approximation of the exponential function." Neural Computation, 11(4), May 1999, pp.853-862.

/*
// FastExpBad is the basic original N.N. Schraudolph version
// which has bad error and is no faster than the better cubic
// and quadratic cases.
func FastExpBad(x float32) float32 {
	i := int32(1512775*x + 1072632447)
	if x <= -88.76731 { // this doesn't add anything and -exp is main use-case anyway
		return 0
	}
	return math.Float32frombits(uint32(i << 32))
}

// FastExp3 is less accurate and no faster than quartic version.

// FastExp3 is a cubic spline approximation to the Exp function, by N.N. Schraudolph
// It does not have any of the sanity checking of a standard method -- returns
// nonsense when arg is out of range.  Runs in .24ns vs. 8.7ns for 64bit which is faster
// than math32.Exp actually.
func FastExp3(x float32) float32 {
	// const (
	// 	Overflow  = 88.43114
	// 	Underflow = -88.76731
	// 	NearZero  = 1.0 / (1 << 28) // 2**-28
	// )

	// special cases
	// switch {
	// these "sanity check" cases cost about 1 ns
	// case IsNaN(x) || IsInf(x, 1): /
	// 	return x
	// case IsInf(x, -1):
	// 	return 0
	// these cases cost about 4+ ns
	// case x >= Overflow:
	// 	return Inf(1)
	// case x <= Underflow:
	// 	return 0
	// case -NearZero < x && x < NearZero:
	// 	return 1 + x
	// }
	if x <= -88.76731 { // this doesn't add anything and -exp is main use-case anyway
		return 0
	}
	i := int32(12102203*x) + 127*(1<<23)
	m := i >> 7 & 0xFFFF // copy mantissa
	i += ((((((((1277 * m) >> 14) + 14825) * m) >> 14) - 79749) * m) >> 11) - 626
	return math.Float32frombits(uint32(i))
}
*/

// FastExp is a quartic spline approximation to the Exp function, by N.N. Schraudolph
// It does not have any of the sanity checking of a standard method -- returns
// nonsense when arg is out of range.  Runs in .24ns vs. 8.7ns for 64bit which is faster
// than math32.Exp actually.
func FastExp(x float32) float32 {
	if x <= -88.76731 { // this doesn't add anything and -exp is main use-case anyway
		return 0
	}
	i := int32(12102203*x) + 127*(1<<23)
	m := i >> 7 & 0xFFFF // copy mantissa
	i += (((((((((((3537 * m) >> 16) + 13668) * m) >> 18) + 15817) * m) >> 14) - 80470) * m) >> 11)
	return math.Float32frombits(uint32(i))
}
