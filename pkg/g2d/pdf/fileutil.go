package pdf

import "github.com/bhojpur/render/pkg/document"

// SaveToPdfFile creates and saves a bdf document to a file
func SaveToPdfFile(filePath string, bdf *document.Bdf) error {
	return bdf.OutputFileAndClose(filePath)
}
