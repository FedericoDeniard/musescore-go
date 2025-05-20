package images

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/mskrha/svg2png"
)

func ConvertSvgToPng(svgPath string) string {

	var input []byte

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
