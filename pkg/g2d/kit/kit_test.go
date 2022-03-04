package kit

import (
	"image"
	"image/color"
	"testing"

	imgkit "github.com/bhojpur/render/pkg/g2d/img"
)

func TestCircle(t *testing.T) {
	width := 200
	height := 200
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	gc := imgkit.NewGraphicContext(img)

	gc.SetStrokeColor(color.NRGBA{255, 255, 255, 255})
	gc.SetFillColor(color.NRGBA{255, 255, 255, 255})
	gc.Clear()

	gc.SetStrokeColor(color.NRGBA{255, 0, 0, 255})
	gc.SetLineWidth(1)

	// Draw a circle
	Circle(gc, 100, 100, 50)
	gc.Stroke()

	imgkit.SaveToPngFile("../output/kit/TestCircle.png", img)
}
