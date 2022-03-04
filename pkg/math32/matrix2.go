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
	"fmt"
	"log"
	"strconv"
	"strings"
	"unicode"
)

// note: golang.org/x/image/math/f64 defines Vec2 as [2]float64
// it is instead very convenient and clear to use .X .Y fields for 2D math
// original gg package used Point2D but Vec2 is more general, e.g., for sizes etc
// in go much better to use fewer types so only using Vec2

type Matrix2 struct {
	XX, YX, XY, YY, X0, Y0 float32
}

func Identity2D() Matrix2 {
	return Matrix2{
		1, 0,
		0, 1,
		0, 0,
	}
}

func (m Matrix2) IsIdentity() bool {
	return m.XX == 1 && m.YX == 0 && m.XY == 0 && m.YY == 1 && m.X0 == 0 && m.Y0 == 0
}

func Translate2D(x, y float32) Matrix2 {
	return Matrix2{
		1, 0,
		0, 1,
		x, y,
	}
}

func Scale2D(x, y float32) Matrix2 {
	return Matrix2{
		x, 0,
		0, y,
		0, 0,
	}
}

func Rotate2D(angle float32) Matrix2 {
	c := float32(Cos(angle))
	s := float32(Sin(angle))
	return Matrix2{
		c, s,
		-s, c,
		0, 0,
	}
}

func Shear2D(x, y float32) Matrix2 {
	return Matrix2{
		1, y,
		x, 1,
		0, 0,
	}
}

func Skew2D(x, y float32) Matrix2 {
	return Matrix2{
		1, Tan(y),
		Tan(x), 1,
		0, 0,
	}
}

func (a Matrix2) Mul(b Matrix2) Matrix2 {
	return Matrix2{
		a.XX*b.XX + a.YX*b.XY,
		a.XX*b.YX + a.YX*b.YY,
		a.XY*b.XX + a.YY*b.XY,
		a.XY*b.YX + a.YY*b.YY,
		a.X0*b.XX + a.Y0*b.XY + b.X0,
		a.X0*b.YX + a.Y0*b.YY + b.Y0,
	}
}

// MulVec2AsVec multiplies the Vec2 as a vector -- does not add translations.
// This is for directional vectors and not points.
func (a Matrix2) MulVec2AsVec(v Vector2) Vector2 {
	tx := a.XX*v.X + a.XY*v.Y
	ty := a.YX*v.X + a.YY*v.Y
	return Vector2{tx, ty}
}

// MulVec2AsPt multiplies the Vec2 as a point, including adding translations.
func (a Matrix2) MulVec2AsPt(v Vector2) Vector2 {
	tx := a.XX*v.X + a.XY*v.Y + a.X0
	ty := a.YX*v.X + a.YY*v.Y + a.Y0
	return Vector2{tx, ty}
}

// MulVec2AsPtCtr multiplies the Vec2 as a point relative to given center-point
// including adding translations.
func (a Matrix2) MulVec2AsPtCtr(v, ctr *Vector2) Vector2 {
	rel := v.Sub(ctr)
	tx := ctr.X + a.XX*rel.X + a.XY*rel.Y + a.X0
	ty := ctr.Y + a.YX*rel.X + a.YY*rel.Y + a.Y0
	return Vector2{tx, ty}
}

// MulCtr multiplies the Mat2, first subtracting given translation center point
// from the translation components, and then adding it back in.
func (a Matrix2) MulCtr(b Matrix2, ctr Vector2) Matrix2 {
	a.X0 -= ctr.X
	a.Y0 -= ctr.Y
	rv := a.Mul(b)
	rv.X0 += ctr.X
	rv.Y0 += ctr.Y
	return rv
}

func (a Matrix2) Translate(x, y float32) Matrix2 {
	return Translate2D(x, y).Mul(a)
}

func (a Matrix2) Scale(x, y float32) Matrix2 {
	return Scale2D(x, y).Mul(a)
}

func (a Matrix2) Rotate(angle float32) Matrix2 {
	return Rotate2D(angle).Mul(a)
}

func (a Matrix2) Shear(x, y float32) Matrix2 {
	return Shear2D(x, y).Mul(a)
}

func (a Matrix2) Skew(x, y float32) Matrix2 {
	return Skew2D(x, y).Mul(a)
}

// ExtractRot extracts the rotation component from a given matrix
func (a Matrix2) ExtractRot() float32 {
	return Atan2(-a.XY, a.XX)
}

// ExtractXYScale extracts the X and Y scale factors after undoing any
// rotation present -- i.e., in the original X, Y coordinates
func (a Matrix2) ExtractScale() (scx, scy float32) {
	rot := a.ExtractRot()
	tx := a.Rotate(-rot)
	scxv := tx.MulVec2AsVec(Vector2{1, 0})
	scyv := tx.MulVec2AsVec(Vector2{0, 1})
	return scxv.X, scyv.Y
}

// Inverse returns inverse of matrix, for inverting transforms
func (a Matrix2) Inverse() Matrix2 {
	// homogenous rep, rc indexes, mapping into Mat3 code
	// XX YX X0   n11 n12 n13    a b x
	// XY YY Y0   n21 n22 n23    c d y
	// 0  0  1    n31 n32 n33    0 0 1

	// t11 := a.YY
	// t12 := -a.YX
	// t13 := a.Y0*a.YX - a.YY*a.X0
	det := a.XX*a.YY - a.XY*a.YX // ad - bc
	detInv := 1 / det

	b := Matrix2{}
	b.XX = a.YY * detInv  // a = d
	b.XY = -a.XY * detInv // c = -c
	b.YX = -a.YX * detInv // b = -b
	b.YY = a.XX * detInv  // d = a
	b.X0 = (a.Y0*a.YX - a.YY*a.X0) * detInv
	b.Y0 = (a.X0*a.XY - a.XX*a.Y0) * detInv
	return b
}

// ParseFloat32 logs any strconv.ParseFloat errors
func ParseFloat32(pstr string) (float32, error) {
	r, err := strconv.ParseFloat(pstr, 32)
	if err != nil {
		log.Printf("gi.ParseFloat32: error parsing float32 number from: %v, %v\n", pstr, err)
		return float32(0.0), err
	}
	return float32(r), nil
}

// ParseAngle32 returns radians angle from string that can specify units (deg,
// grad, rad) -- deg is assumed if not specified
func ParseAngle32(pstr string) (float32, error) {
	units := "deg"
	lstr := strings.ToLower(pstr)
	if strings.Contains(lstr, "deg") {
		units = "deg"
		lstr = strings.TrimSuffix(lstr, "deg")
	} else if strings.Contains(lstr, "grad") {
		units = "grad"
		lstr = strings.TrimSuffix(lstr, "grad")
	} else if strings.Contains(lstr, "rad") {
		units = "rad"
		lstr = strings.TrimSuffix(lstr, "rad")
	}
	r, err := strconv.ParseFloat(lstr, 32)
	if err != nil {
		log.Printf("gi.ParseAngle32: error parsing float32 number from: %v, %v\n", lstr, err)
		return float32(0.0), err
	}
	switch units {
	case "deg":
		return float32(r) * Pi / 180, nil
	case "grad":
		return float32(r) * Pi / 200, nil
	case "rad":
		return float32(r), nil
	}
	return float32(r), nil
}

// ReadPoints reads a set of floating point values from a SVG format number
// string -- returns a slice or nil if there was an error
func ReadPoints(pstr string) []float32 {
	lastIdx := -1
	var pts []float32
	lr := ' '
	for i, r := range pstr {
		if unicode.IsNumber(r) == false && r != '.' && !(r == '-' && lr == 'e') && r != 'e' {
			if lastIdx != -1 {
				s := pstr[lastIdx:i]
				p, err := ParseFloat32(s)
				if err != nil {
					return nil
				}
				pts = append(pts, p)
			}
			if r == '-' {
				lastIdx = i
			} else {
				lastIdx = -1
			}
		} else if lastIdx == -1 {
			lastIdx = i
		}
		lr = r
	}
	if lastIdx != -1 && lastIdx != len(pstr) {
		s := pstr[lastIdx:len(pstr)]
		p, err := ParseFloat32(s)
		if err != nil {
			return nil
		}
		pts = append(pts, p)
	}
	return pts
}

// PointsCheckN checks the number of points read and emits an error if not equal to n
func PointsCheckN(pts []float32, n int, errmsg string) error {
	if len(pts) != n {
		return fmt.Errorf("%v incorrect number of points: %v != %v\n", errmsg, len(pts), n)
	}
	return nil
}

// SetString processes the standard SVG-style transform strings
func (a *Matrix2) SetString(str string) error {
	errmsg := "gi.Mat2 SetString"
	str = strings.ToLower(strings.TrimSpace(str))
	*a = Identity2D()
	if str == "none" {
		*a = Identity2D()
		return nil
	}
	// could have multiple transforms
	for {
		pidx := strings.IndexByte(str, '(')
		if pidx < 0 {
			err := fmt.Errorf("gi.Mat2 SetString: no params for xform: %v\n", str)
			log.Println(err)
			return err
		}
		cmd := str[:pidx]
		vals := str[pidx+1:]
		nxt := ""
		eidx := strings.IndexByte(vals, ')')
		if eidx > 0 {
			nxt = strings.TrimSpace(vals[eidx+1:])
			if strings.HasPrefix(nxt, ";") {
				nxt = strings.TrimSpace(strings.TrimPrefix(nxt, ";"))
			}
			vals = vals[:eidx]
		}
		pts := ReadPoints(vals)
		switch cmd {
		case "matrix":
			if err := PointsCheckN(pts, 6, errmsg); err != nil {
				log.Println(err)
				return err
			}
			*a = Matrix2{pts[0], pts[1], pts[2], pts[3], pts[4], pts[5]}
		case "translate":
			if err := PointsCheckN(pts, 2, errmsg); err != nil {
				log.Println(err)
				return err
			}
			*a = a.Translate(pts[0], pts[1])
		case "translatex":
			if err := PointsCheckN(pts, 1, errmsg); err != nil {
				log.Println(err)
				return err
			}
			*a = a.Translate(pts[0], 0)
		case "translatey":
			if err := PointsCheckN(pts, 1, errmsg); err != nil {
				log.Println(err)
				return err
			}
			*a = a.Translate(0, pts[0])
		case "scale":
			if len(pts) == 1 {
				*a = a.Scale(pts[0], pts[0])
			} else if len(pts) == 2 {
				*a = a.Scale(pts[0], pts[1])
			} else {
				err := fmt.Errorf("%v incorrect number of points: 2 != %v\n", errmsg, len(pts))
				log.Println(err)
			}
		case "scalex":
			if err := PointsCheckN(pts, 1, errmsg); err != nil {
				log.Println(err)
				return err
			}
			*a = a.Scale(pts[0], 1)
		case "scaley":
			if err := PointsCheckN(pts, 1, errmsg); err != nil {
				log.Println(err)
				return err
			}
			*a = a.Scale(1, pts[0])
		case "rotate":
			ang := pts[0] * Pi / 180 // always in degrees in this form
			if len(pts) == 3 {
				*a = a.Translate(pts[1], pts[2]).Rotate(ang).Translate(-pts[1], -pts[2])
			} else if len(pts) == 1 {
				*a = a.Rotate(ang)
			} else {
				return PointsCheckN(pts, 1, errmsg)
			}
		case "skew":
			if err := PointsCheckN(pts, 2, errmsg); err != nil {
				log.Println(err)
				return err
			}
			*a = a.Skew(pts[0], pts[1])
		case "skewx":
			if err := PointsCheckN(pts, 1, errmsg); err != nil {
				log.Println(err)
				return err
			}
			*a = a.Skew(pts[0], 0)
		case "skewy":
			if err := PointsCheckN(pts, 1, errmsg); err != nil {
				log.Println(err)
				return err
			}
			*a = a.Skew(0, pts[0])
		}
		if nxt == "" {
			break
		}
		if !strings.Contains(nxt, "(") {
			break
		}
		str = nxt
	}
	return nil
}

// String returns the XML-based string representation of the transform
func (a *Matrix2) String() string {
	if a.IsIdentity() {
		return "none"
	}
	if a.YX == 0 && a.XY == 0 { // no rotation, emit scale and translate
		str := ""
		if a.XX != 1 || a.YY != 1 {
			str = fmt.Sprintf("scale(%g,%g)", a.XX, a.YY)
		}
		if a.X0 != 0 || a.Y0 != 0 {
			if str != "" {
				str += " "
			}
			str += fmt.Sprintf("translate(%g,%g)", a.X0, a.Y0)
		}
		return str
	}
	// just report the whole matrix
	return fmt.Sprintf("matrix(%g,%g,%g,%g,%g,%g)", a.XX, a.YX, a.XY, a.YY, a.X0, a.Y0)
}
