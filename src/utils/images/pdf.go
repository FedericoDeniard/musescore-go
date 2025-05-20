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
			httpError := customErrors.HttpError{StatusCode: 500, Message: "Error al leer el archivo"}
			fmt.Println(httpError.Error())
			return "", &httpError
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
		httpError := customErrors.HttpError{StatusCode: 500, Message: "Error al crear directorio"}
		fmt.Println(httpError.Error())
		return "", &httpError
	}

	if err := pdf.OutputFileAndClose(outputPDF); err != nil {
		httpError := customErrors.HttpError{StatusCode: 500, Message: "Error al guardar el PDF"}
		fmt.Println(httpError.Error())
		return "", &httpError
	}

	return outputPDF, nil
}
