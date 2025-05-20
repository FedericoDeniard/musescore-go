package images

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/jung-kurt/gofpdf"
)

func ConvertPngToPdf(pngPaths ...string) (string, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pages := make([]ImgPage, len(pngPaths))

	var wg sync.WaitGroup

	for i, path := range pngPaths {
		wg.Add(1)
		go func(path string, index int) {
			defer wg.Done()
			data, err := os.ReadFile(path)
			pages[index] = ImgPage{data, err}
		}(path, i)
	}

	wg.Wait()

	for i, page := range pages {
		if page.err != nil {
			return "", fmt.Errorf("error leyendo PNG %s: %w", pngPaths[i], page.err)
		}

		// 2. Añadir al PDF
		pdf.AddPageFormat("P", gofpdf.SizeType{
			Wd: 210,
			Ht: 297,
		})

		opt := gofpdf.ImageOptions{
			ImageType: "PNG",
			ReadDpi:   true,
		}
		imageName := filepath.Base(pngPaths[i])
		pdf.RegisterImageOptionsReader(imageName, opt, bytes.NewReader(page.data))
		pdf.Image(imageName, 0, 0, 210, 297, false, "", 0, "")
	}

	outputPDF := filepath.Join(
		filepath.Dir(pngPaths[0]),
		filepath.Base(pngPaths[0][:len(pngPaths[0])-len(filepath.Ext(pngPaths[0]))]),
	) + ".pdf"

	if err := os.MkdirAll(filepath.Dir(outputPDF), 0755); err != nil {
		return "", fmt.Errorf("error creando directorio: %w", err)
	}

	if err := pdf.OutputFileAndClose(outputPDF); err != nil {
		return "", fmt.Errorf("error guardando PDF: %w", err)
	}

	return outputPDF, nil
}
