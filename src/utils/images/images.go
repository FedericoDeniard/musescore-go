package images

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	customErrors "github.com/FedericoDeniard/musescore-go/src/utils/error"
	"github.com/google/uuid"
)

func DownloadImage(url string) (string, *customErrors.HttpError) {
	imagesFolder := "src/downloads/images/"
	var extension string

	err := os.MkdirAll(imagesFolder, 0755)
	if err != nil {
		fmt.Println("Error al crear directorio:", err)
		return "", &customErrors.HttpError{
			StatusCode: 500,
			Message:    "Error al crear directorio",
		}
	}

	resp, err := http.Get(url)
	if err != nil {
		return "", &customErrors.HttpError{
			StatusCode: 500,
			Message:    "Error al descargar la imagen",
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error al descargar la imagen:", resp.StatusCode)
		return "", &customErrors.HttpError{
			StatusCode: 500,
			Message:    "Error al descargar la imagen",
		}
	}

	extension, httpError := GetExtensionFromImage(url)
	if httpError != nil {
		return "", httpError
	}

	fileName := uuid.New().String()
	filePath := filepath.Join(imagesFolder, fileName+extension)

	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Error al crear el archivo:", err)
		return "", &customErrors.HttpError{
			StatusCode: 500,
			Message:    "Error al crear el archivo",
		}
	}
	defer file.Close()

	if _, err := io.Copy(file, resp.Body); err != nil {
		fmt.Println("Error al copiar el archivo:", err)
		return "", &customErrors.HttpError{
			StatusCode: 500,
			Message:    "Error al copiar el archivo",
		}
	}

	return filePath, nil
}

func GetExtensionFromImage(url string) (string, *customErrors.HttpError) {
	if strings.Contains(url, ".svg") || strings.Contains(url, "image/svg+xml") {
		return ".svg", nil
	} else if strings.Contains(url, ".png") || strings.Contains(url, "image/png") {
		return ".png", nil
	} else if strings.Contains(url, ".jpg") || strings.Contains(url, "image/jpg") {
		return ".jpg", nil
	}
	return "", &customErrors.HttpError{
		StatusCode: 500,
		Message:    "Error al obtener la extension",
	}
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
