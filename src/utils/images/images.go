package images

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

func DownloadImage(url string) string {
	imagesFolder := "src/downloads/images/"
	var extension string

	resp, err := http.Get(url)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ""
	}

	extension = GetExtensionFromImage(url)
	fileName := uuid.New().String()
	filePath := filepath.Join(imagesFolder, fileName+extension)

	file, err := os.Create(filePath)
	if err != nil {
		return ""
	}
	defer file.Close()

	if _, err := io.Copy(file, resp.Body); err != nil {
		return ""
	}

	return filePath
}

func GetExtensionFromImage(url string) string {
	if strings.Contains(url, ".svg") || strings.Contains(url, "image/svg+xml") {
		return ".svg"
	} else if strings.Contains(url, ".png") || strings.Contains(url, "image/png") {
		return ".png"
	} else if strings.Contains(url, ".jpg") || strings.Contains(url, "image/jpg") {
		return ".jpg"
	}
	panic("Error al obtener la extension")
}

type ImgPage struct {
	data []byte
	err  error
}

func DeleteImages(paths ...string) {
	for _, path := range paths {
		os.Remove(path)
	}
}
