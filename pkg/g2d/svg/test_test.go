package svg_test

// It gives test coverage with the command:
// go test -cover ./... | grep -v "no test"

import (
	"testing"

	d2d "github.com/bhojpur/render/pkg/g2d/draw"
	svgkit "github.com/bhojpur/render/pkg/g2d/svg"
)

type sample func(gc d2d.GraphicContext, ext string) (string, error)

func test(t *testing.T, draw sample) {
	// Initialize the graphic context on a .pdf document
	dest := svgkit.NewSvg()
	gc := svgkit.NewGraphicContext(dest)
	// Draw sample
	output, err := draw(gc, "svg")
	if err != nil {
		t.Errorf("Drawing %q failed: %v", output, err)
		return
	}
	err = svgkit.SaveToSvgFile(output, dest)
	if err != nil {
		t.Errorf("Saving %q failed: %v", output, err)
	}
}
