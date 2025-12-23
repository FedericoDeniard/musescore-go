package scrap

import (
	"fmt"
	"sync"
	"time"

	customErrors "github.com/FedericoDeniard/musescore-go/src/utils/error"
	"github.com/FedericoDeniard/musescore-go/src/utils/images"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

func Scrap(browser *rod.Browser, url string) (string, *customErrors.HttpError) {
	fmt.Println("Scraping process started for URL:", url)
	defer browser.MustClose()

	page, err := browser.Page(proto.TargetCreateTarget{
		URL: url,
	})
	if err != nil {
		fmt.Println("Error creating page:", err)
		return "", &customErrors.HttpError{StatusCode: 400, Message: "La dirección web que ingresaste no es válida. Verifica que sea correcta e inténtalo nuevamente."}
	}
	fmt.Println("Page created")
	defer page.MustClose()
	fmt.Println("Page loaded")
	page.MustSetViewport(1920, 1080, 1, false)
	fmt.Println("Viewport set")

	scrollerComponent, err := page.Timeout(10 * time.Second).Element("#jmuse-scroller-component")
	if err != nil || scrollerComponent == nil {
		fmt.Println("Scroller component not found")
		httpError := customErrors.HttpError{StatusCode: 404, Message: "No se encontró el componente jmuse-scroller-component"}
		fmt.Println(httpError.Error())
		return "", &httpError
	}
	scrollerComponent, _ = page.Element("#jmuse-scroller-component")
	fmt.Println("Scroller component found")

	var wg sync.WaitGroup
	wg.Add(1)

	// var title string
	sheetsChannel := make(chan string)
	var sheetsSource []string
	var imagesPath []string

	// go func() { defer wg.Done(); title = getTitle(page) }()
	go func() {
		defer wg.Done()
		getSheets(scrollerComponent, sheetsChannel)
	}()

	for sheet := range sheetsChannel {
		sheetsSource = append(sheetsSource, sheet)
		tempImage, err := images.DownloadImage(sheet)
		if err != nil {
			images.DeleteImages(imagesPath...)
			return "", err
		}
		imagesPath = append(imagesPath, tempImage)
	}

	wg.Wait()
	// fmt.Printf("Title: %s \n", title)
	fmt.Printf("Partituras encontradas: %d \n", len(sheetsSource))
	fmt.Printf("Partituras descargadas: %d \n", len(imagesPath))

	if len(imagesPath) == 0 {
		return "", &customErrors.HttpError{StatusCode: 404, Message: "No se encontraron partituras en la página proporcionada."}
	}

	imagesExtensions, httpError := images.GetExtensionFromImage(imagesPath[0])
	if httpError != nil {
		images.DeleteImages(imagesPath...)
		return "", httpError
	}
	fmt.Println("Extension obtained:", imagesExtensions)

	var convertedImages []string

	if imagesExtensions == ".svg" {
		fmt.Println("Converting SVG to PNG")
		convertedImages, httpError = images.ConvertMultipleSvgToPng(imagesPath...)
		if httpError != nil {
			images.DeleteImages(imagesPath...)
			return "", httpError
		}
		fmt.Println("SVG conversion completed")
	} else if imagesExtensions == ".png" {
		fmt.Println("Images are PNG")
		convertedImages = imagesPath
	} else {
		httpError := customErrors.HttpError{StatusCode: 501, Message: "Extension no soportada"}
		fmt.Println(httpError.Error())
		images.DeleteImages(imagesPath...)
		return "", &httpError
	}

	pdfPath, httpError := images.ConvertPngToPdf(convertedImages...)
	if httpError != nil {
		images.DeleteImages(imagesPath...)
		return "", httpError
	}
	fmt.Println("PDF created")

	filesToDelete := append(imagesPath, convertedImages...)
	images.DeleteImages(filesToDelete...)
	fmt.Println("Process finished")
	return pdfPath, nil
}

// func getTitle(page *rod.Page) string {
// 	title := "musescore"
// 	asideContainer := page.MustElement("#aside-container-unique")
// 	if asideContainer != nil {
// 		titleElement := asideContainer.MustElement(".nFRPI")
// 		if titleElement != nil {
// 			span := titleElement.MustElement("span")
// 			if span != nil {
// 				title = span.MustText()
// 			}
// 		}
// 	}
// 	return title
// }

func getSheets(component *rod.Element, channel chan<- string) {
	fmt.Println("Getting sheets...")
	defer close(channel)

	component.MustEval(`() => this.scrollIntoView({ behavior: "smooth", block: "start", inline: "nearest" })`)
	page := component.Page()

	page.MustWaitRequestIdle()
	sheets := page.MustElements(".A8huy")
	fmt.Println("Sheets found: ", len(sheets))

	for i, sheet := range sheets {
		fmt.Printf("Procesando hoja %d...\n", i+1)

		sheet.MustEval(`() => this.scrollIntoView({ behavior: "smooth", block: "start", inline: "nearest" })`)

		err := page.Timeout(10 * time.Second).Wait(&rod.EvalOptions{
			ThisObj: sheet.Object,
			JS: `() => {
		const img = this.querySelector("img");
		return img && img.src && img.src.trim() !== "";
	}`,
		})

		if err != nil {
			fmt.Printf("Imagen %d no cargó a tiempo: %v\n", i, err)
			continue
		}

		if img, err := sheet.Element("img"); err == nil && img != nil {
			if srcAttr, _ := img.Attribute("src"); srcAttr != nil && *srcAttr != "" {
				channel <- *srcAttr
			} else {
				channel <- ""
				fmt.Printf("Imagen %d no tiene atributo src\n", i)
			}
		} else {
			fmt.Printf("No se encontró <img> en la hoja %d\n", i)
		}
	}
}
