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
	"math"
)

// niceNum returns a "nice" number approximately equal to x. The number is
// rounded if round is true, converted to its ceiling otherwise.
func niceNum(val float64, round bool) float64 {
	var nf float64

	exp := int(math.Floor(math.Log10(val)))
	f := val / math.Pow10(exp)
	if round {
		switch {
		case f < 1.5:
			nf = 1
		case f < 3.0:
			nf = 2
		case f < 7.0:
			nf = 5
		default:
			nf = 10
		}
	} else {
		switch {
		case f <= 1:
			nf = 1
		case f <= 2.0:
			nf = 2
		case f <= 5.0:
			nf = 5
		default:
			nf = 10
		}
	}
	return nf * math.Pow10(exp)
}

// TickmarkPrecision returns an appropriate precision value for label
// formatting.
func TickmarkPrecision(div float64) int {
	return int(math.Max(-math.Floor(math.Log10(div)), 0))
}

// Tickmarks returns a slice of tickmarks appropriate for a chart axis and an
// appropriate precision for formatting purposes. The values min and max will
// be contained within the tickmark range.
func Tickmarks(min, max float64) (list []float64, precision int) {
	if max > min {
		spread := niceNum(max-min, false)
		d := niceNum((spread / 4), true)
		graphMin := math.Floor(min/d) * d
		graphMax := math.Ceil(max/d) * d
		precision = TickmarkPrecision(d)
		for x := graphMin; x < graphMax+0.5*d; x += d {
			list = append(list, x)
		}
	}
	return
}
