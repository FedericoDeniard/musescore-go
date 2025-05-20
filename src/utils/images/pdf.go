package images

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	customErrors "github.com/FedericoDeniard/musescore-go/src/utils/error"
	"github.com/jung-kurt/gofpdf"
)

func ConvertPngToPdf(pngPaths ...string) (string, *customErrors.HttpError) {
	pdf := gofpdf.New("P", "mm", "A4", "")

	for _, path := range pngPaths {
		data, err := os.ReadFile(path)
		if err != nil {
			return "", &customErrors.HttpError{
				StatusCode: 500,
				Message:    fmt.Sprintf("error leyendo PNG %s: %w", path, err),
			}
		}
		pdf.AddPageFormat("P", gofpdf.SizeType{
			Wd: 210,
			Ht: 297,
		})
		opt := gofpdf.ImageOptions{
			ImageType: "PNG",
			ReadDpi:   true,
		}

		imageName := filepath.Base(path)
		pdf.RegisterImageOptionsReader(imageName, opt, bytes.NewReader(data))
		pdf.Image(imageName, 0, 0, 210, 297, false, "", 0, "")

	}

	outputPDF := filepath.Join(
		filepath.Dir(pngPaths[0]),
		filepath.Base(pngPaths[0][:len(pngPaths[0])-len(filepath.Ext(pngPaths[0]))]),
	) + ".pdf"

	if err := os.MkdirAll(filepath.Dir(outputPDF), 0755); err != nil {
		return "", &customErrors.HttpError{
			StatusCode: 500,
			Message:    fmt.Sprintf("error creando directorio: %w", err),
		}
	}

	if err := pdf.OutputFileAndClose(outputPDF); err != nil {
		return "", &customErrors.HttpError{
			StatusCode: 500,
			Message:    fmt.Sprintf("error guardando PDF: %w", err),
		}
	}

	return outputPDF, nil
}
