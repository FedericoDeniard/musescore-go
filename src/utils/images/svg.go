package images

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	customErrors "github.com/FedericoDeniard/musescore-go/src/utils/error"
	"github.com/mskrha/svg2png"
)

func ConvertSvgToPng(svgPath string) (string, *customErrors.HttpError) {

	var input []byte

	input, err := os.ReadFile(svgPath)
	if err != nil {
		httpError := customErrors.HttpError{StatusCode: 500, Message: "Error al leer el archivo SVG"}
		fmt.Println(httpError.Error())
		return "", &httpError
	}
	converter := svg2png.New()
	output, err := converter.Convert(input)
	if err != nil {
		httpError := customErrors.HttpError{StatusCode: 500, Message: "Error al convertir el archivo SVG a PNG"}
		fmt.Println(httpError.Error())
		return "", &httpError
	}

	pngPath := filepath.Join(filepath.Dir(svgPath), filepath.Base(svgPath[:len(svgPath)-len(filepath.Ext(svgPath))])+".png")
	err = os.WriteFile(pngPath, output, 0644)
	if err != nil {
		httpError := customErrors.HttpError{StatusCode: 500, Message: "Error al escribir el archivo PNG"}
		fmt.Println(httpError.Error())
		return "", &httpError
	}
	return pngPath, nil
}

func ConvertMultipleSvgToPng(svgPaths ...string) ([]string, *customErrors.HttpError) {
	var wg sync.WaitGroup
	outputPaths := make([]string, len(svgPaths))
	errChan := make(chan *customErrors.HttpError, len(svgPaths))
	for i, svgPath := range svgPaths {
		wg.Add(1)
		go func(path string, index int) {
			defer wg.Done()
			out, htppError := ConvertSvgToPng(path)
			if htppError != nil {
				errChan <- htppError
				return
			}
			outputPaths[index] = out
		}(svgPath, i)
	}
	wg.Wait()
	close(errChan)
	for err := range errChan {
		var createdImages []string
		for _, path := range outputPaths {
			if path != "" {
				createdImages = append(createdImages, path)
			}
		}
		if len(createdImages) > 0 {
			DeleteImages(createdImages...)
		}
		return nil, err
	}
	return outputPaths, nil
}
