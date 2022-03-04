package pdf_test

// It gives test coverage with the command:
// go test -cover ./... | grep -v "no test"

import (
	"testing"

	d2d "github.com/bhojpur/render/pkg/g2d/draw"
	d2pdf "github.com/bhojpur/render/pkg/g2d/pdf"
)

type sample func(gc d2d.GraphicContext, ext string) (string, error)

func test(t *testing.T, draw sample) {
	// Initialize the graphic context on an pdf document
	dest := d2pdf.NewPdf("L", "mm", "A4")
	gc := d2pdf.NewGraphicContext(dest)
	// Draw sample
	output, err := draw(gc, "pdf")
	if err != nil {
		t.Errorf("Drawing %q failed: %v", output, err)
		return
	}
	/*
		// Save to pdf only if it doesn't exist because of git
		if _, err = os.Stat(output); err == nil {
			t.Skipf("Saving %q skipped, as it exists already. (Git would consider it modified.)", output)
			return
		}
	*/
	err = d2pdf.SaveToPdfFile(output, dest)
	if err != nil {
		t.Errorf("Saving %q failed: %v", output, err)
	}
}
