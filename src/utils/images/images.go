package images

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/jung-kurt/gofpdf"
	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
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

func ConvertMusescoreSVGsToPDF(svgPaths []string, outputPDF string) (string, error) {
	// Crear nuevo PDF
	pdf := gofpdf.New("P", "mm", "A4", "")

	for _, svgPath := range svgPaths {
		// 1. Leer archivo SVG
		svgData, err := os.ReadFile(svgPath)
		if err != nil {
			return "", fmt.Errorf("error leyendo SVG %s: %w", svgPath, err)
		}

		// 2. Parsear SVG con oksvg (más tolerante)
		icon, err := oksvg.ReadIconStream(bytes.NewReader(svgData))
		if err != nil {
			return "", fmt.Errorf("error parseando SVG %s: %w", svgPath, err)
		}

		// 3. Configurar dimensiones (similar a tu ajuste de -800 en Node)
		width := icon.ViewBox.W
		height := icon.ViewBox.H
		if width < 0 {
			width = icon.ViewBox.W
		}
		if height < 0 {
			height = icon.ViewBox.H
		}

		// 4. Renderizar SVG a imagen PNG en memoria
		img := image.NewRGBA(image.Rect(0, 0, int(width), int(height)))
		scannerGV := rasterx.NewScannerGV(int(width), int(height), img, img.Bounds())
		raster := rasterx.NewDasher(int(width), int(height), scannerGV)
		icon.Draw(raster, 1.0)

		// 5. Guardar PNG temporal en memoria
		var pngBuf bytes.Buffer
		if err := png.Encode(&pngBuf, img); err != nil {
			return "", fmt.Errorf("error convirtiendo a PNG: %w", err)
		}

		// 6. Añadir al PDF
		pdf.AddPageFormat("P", gofpdf.SizeType{
			Wd: width / 3.78, // Convertir px a mm (96dpi -> 3.78px/mm)
			Ht: height / 3.78,
		})

		// Registrar la imagen PNG en el PDF
		opt := gofpdf.ImageOptions{
			ImageType: "PNG",
			ReadDpi:   true,
		}
		imageName := filepath.Base(svgPath)
		pdf.RegisterImageOptionsReader(imageName, opt, &pngBuf)
		pdf.Image(imageName, 0, 0, width/3.78, height/3.78, false, "", 0, "")
	}

	// Crear directorio de salida si no existe
	if err := os.MkdirAll(filepath.Dir(outputPDF), 0755); err != nil {
		return "", fmt.Errorf("error creando directorio: %w", err)
	}

	// Guardar PDF
	if err := pdf.OutputFileAndClose(outputPDF + ".pdf"); err != nil {
		return "", fmt.Errorf("error guardando PDF: %w", err)
	}

	// Devolver ruta absoluta
	absPath, err := filepath.Abs(outputPDF)
	if err != nil {
		return outputPDF, nil
	}
	return absPath, nil
}
