package images

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/jung-kurt/gofpdf"
	"github.com/mskrha/svg2png"
)

func DownloadImage(urls []string) ([]string, error) {
	var imagesTempNames []string
	imagesFolder := "src/downloads/images/"
	var extension string

	for _, url := range urls {
		resp, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, err
		}

		extension = getExtensionFromImage(url)
		fmt.Println(extension)
		fileName := uuid.New().String()
		filePath := filepath.Join(imagesFolder, fileName+extension)

		file, err := os.Create(filePath)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		if _, err := io.Copy(file, resp.Body); err != nil {
			return nil, err
		}

		imagesTempNames = append(imagesTempNames, filePath)
	}
	return imagesTempNames, nil
}

func getExtensionFromImage(url string) string {
	if strings.Contains(url, ".svg") || strings.Contains(url, "image/svg+xml") {
		return ".svg"
	} else if strings.Contains(url, ".png") || strings.Contains(url, "image/png") {
		return ".png"
	} else if strings.Contains(url, ".jpg") || strings.Contains(url, "image/jpg") {
		return ".jpg"
	}
	panic("Error al obtener la extension")
}
func ConvertSvgToPng(svgPath string) string {

	var input []byte

	// Fill input with svg data
	input, err := os.ReadFile(svgPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	converter := svg2png.New()
	output, err := converter.Convert(input)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	pngPath := filepath.Join(filepath.Dir(svgPath), filepath.Base(svgPath[:len(svgPath)-len(filepath.Ext(svgPath))])+".png")
	err = os.WriteFile(pngPath, output, 0644)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error escribiendo el archivo PNG:", err)
	}
	return pngPath
}

func ConvertMultipleSvgToPng(svgPaths ...string) ([]string, error) {
	var wg sync.WaitGroup
	outputPaths := make([]string, len(svgPaths))
	for i, svgPath := range svgPaths {
		wg.Add(1)
		go func(path string, index int) {
			defer wg.Done()
			outputPaths[index] = ConvertSvgToPng(path)
		}(svgPath, i)
	}
	wg.Wait()
	return outputPaths, nil
}

func ConvertPngToPdf(pngPaths ...string) (string, error) {
	// Crear nuevo PDF
	pdf := gofpdf.New("P", "mm", "A4", "")
	for _, pngPath := range pngPaths {
		fmt.Println(pngPath)
	}

	for _, pngPath := range pngPaths {
		// 1. Leer archivo PNG
		pngData, err := os.ReadFile(pngPath)
		if err != nil {
			return "", fmt.Errorf("error leyendo PNG %s: %w", pngPath, err)
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
		imageName := filepath.Base(pngPath)
		pdf.RegisterImageOptionsReader(imageName, opt, bytes.NewReader(pngData))
		pdf.Image(imageName, 0, 0, 210, 297, false, "", 0, "")
	}

	outputPDF := filepath.Join(filepath.Dir(pngPaths[0]), filepath.Base(pngPaths[0][:len(pngPaths[0])-len(filepath.Ext(pngPaths[0]))]))

	// Crear directorio de salida si no existe
	if err := os.MkdirAll(filepath.Dir(outputPDF), 0755); err != nil {
		return "", fmt.Errorf("error creando directorio: %w", err)
	}

	// Guardar PDF
	if err := pdf.OutputFileAndClose(outputPDF + ".pdf"); err != nil {
		return "", fmt.Errorf("error guardando PDF: %w", err)
	}

	return outputPDF, nil
}
