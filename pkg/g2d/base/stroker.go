package base

import (
	"math"

	d2d "github.com/bhojpur/render/pkg/g2d/draw"
)

type LineStroker struct {
	Flattener     Flattener
	HalfLineWidth float64
	Cap           d2d.LineCap
	Join          d2d.LineJoin
	vertices      []float64
	rewind        []float64
	x, y, nx, ny  float64
}

func NewLineStroker(c d2d.LineCap, j d2d.LineJoin, flattener Flattener) *LineStroker {
	l := new(LineStroker)
	l.Flattener = flattener
	l.HalfLineWidth = 0.5
	l.Cap = c
	l.Join = j
	return l
}

func (l *LineStroker) MoveTo(x, y float64) {
	l.x, l.y = x, y
}

func (l *LineStroker) LineTo(x, y float64) {
	l.line(l.x, l.y, x, y)
}

func (l *LineStroker) LineJoin() {

}

func (l *LineStroker) line(x1, y1, x2, y2 float64) {
	dx := (x2 - x1)
	dy := (y2 - y1)
	d := vectorDistance(dx, dy)
	if d != 0 {
		nx := dy * l.HalfLineWidth / d
		ny := -(dx * l.HalfLineWidth / d)
		l.appendVertex(x1+nx, y1+ny, x2+nx, y2+ny, x1-nx, y1-ny, x2-nx, y2-ny)
		l.x, l.y, l.nx, l.ny = x2, y2, nx, ny
	}
}

func (l *LineStroker) Close() {
	if len(l.vertices) > 1 {
		l.appendVertex(l.vertices[0], l.vertices[1], l.rewind[0], l.rewind[1])
	}
}

func (l *LineStroker) End() {
	if len(l.vertices) > 1 {
		l.Flattener.MoveTo(l.vertices[0], l.vertices[1])
		for i, j := 2, 3; j < len(l.vertices); i, j = i+2, j+2 {
			l.Flattener.LineTo(l.vertices[i], l.vertices[j])
		}
	}
	for i, j := len(l.rewind)-2, len(l.rewind)-1; j > 0; i, j = i-2, j-2 {
		l.Flattener.LineTo(l.rewind[i], l.rewind[j])
	}
	if len(l.vertices) > 1 {
		l.Flattener.LineTo(l.vertices[0], l.vertices[1])
	}
	l.Flattener.End()
	// reinit vertices
	l.vertices = l.vertices[0:0]
	l.rewind = l.rewind[0:0]
	l.x, l.y, l.nx, l.ny = 0, 0, 0, 0

}

func (l *LineStroker) appendVertex(vertices ...float64) {
	s := len(vertices) / 2
	l.vertices = append(l.vertices, vertices[:s]...)
	l.rewind = append(l.rewind, vertices[s:]...)
}

func vectorDistance(dx, dy float64) float64 {
	return float64(math.Sqrt(dx*dx + dy*dy))
}
