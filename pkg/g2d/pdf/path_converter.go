package pdf

import (
	"math"

	d2d "github.com/bhojpur/render/pkg/g2d/draw"
)

const deg = 180 / math.Pi

// ConvertPath converts a paths to the pdf api
func ConvertPath(path *d2d.Path, pdf Vectorizer) {
	var startX, startY float64 = 0, 0
	i := 0
	for _, cmp := range path.Components {
		switch cmp {
		case d2d.MoveToCmp:
			startX, startY = path.Points[i], path.Points[i+1]
			pdf.MoveTo(startX, startY)
			i += 2
		case d2d.LineToCmp:
			pdf.LineTo(path.Points[i], path.Points[i+1])
			i += 2
		case d2d.QuadCurveToCmp:
			pdf.CurveTo(path.Points[i], path.Points[i+1], path.Points[i+2], path.Points[i+3])
			i += 4
		case d2d.CubicCurveToCmp:
			pdf.CurveBezierCubicTo(path.Points[i], path.Points[i+1], path.Points[i+2], path.Points[i+3], path.Points[i+4], path.Points[i+5])
			i += 6
		case d2d.ArcToCmp:
			pdf.ArcTo(path.Points[i], path.Points[i+1], path.Points[i+2], path.Points[i+3],
				0,                    // degRotate
				path.Points[i+4]*deg, // degStart = startAngle
				(path.Points[i+4]-path.Points[i+5])*deg) // degEnd = startAngle-angle
			i += 6
		case d2d.CloseCmp:
			pdf.LineTo(startX, startY)
			pdf.ClosePath()
		}
	}
}
