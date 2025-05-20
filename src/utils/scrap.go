package scrap

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	customErrors "github.com/FedericoDeniard/musescore-go/src/utils/error"
	"github.com/FedericoDeniard/musescore-go/src/utils/images"
	"github.com/go-rod/rod"
)

func Scrap(browser *rod.Browser, url string) (string, *customErrors.HttpError) {
	fmt.Println("Scraping process started...")
	defer browser.MustClose()

	page := browser.MustPage(url)
	fmt.Println("Page created")
	defer page.MustClose()
	fmt.Println("Page loaded")
	page.MustSetViewport(1920, 1080, 1, false)
	fmt.Println("Viewport set")

	scrollerComponent, err := page.Timeout(10 * time.Second).Element("#jmuse-scroller-component")
	fmt.Println(scrollerComponent)
	if err != nil || scrollerComponent == nil {
		return "", &customErrors.HttpError{
			StatusCode: 400,
			Message:    "No se encontró el componente jmuse-scroller-component",
		}
	}
	scrollerComponent, _ = page.Element("#jmuse-scroller-component")

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

	fmt.Println(imagesPath)
	fmt.Printf("Partituras descargadas: %d \n%v", len(imagesPath), imagesPath)

	imagesExtensions, httpError := images.GetExtensionFromImage(imagesPath[0])
	if httpError != nil {
		return "", httpError
	}

	var convertedImages []string

	if imagesExtensions == ".svg" {
		convertedImages, httpError = images.ConvertMultipleSvgToPng(imagesPath...)
		if httpError != nil {
			return "", httpError
		}
	} else if imagesExtensions == ".png" {
		convertedImages = imagesPath
	} else {
		return "", &customErrors.HttpError{
			StatusCode: 501,
			Message:    "Extension no soportada",
		}
	}

	pdfPath, httpError := images.ConvertPngToPdf(convertedImages...)
	if httpError != nil {
		return "", httpError
	}

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
	sheets := page.MustElements(".EEnGW")
	fmt.Println("Sheets found: ", len(sheets))

	for i, sheet := range sheets {
		fmt.Printf("Procesando hoja %d...\n", i+1)

		sheet.MustEval(`() => this.scrollIntoView({ behavior: "smooth", block: "start", inline: "nearest" })`)

		err := page.Timeout(10 * time.Second).Wait(&rod.EvalOptions{
			ThisObj: sheet.Object,
			JS: `() => {
		const img = this.querySelector("img");
		return img && img.complete && img.naturalHeight !== 0;
	}`,
		})

		if err != nil {
			fmt.Printf("Imagen %d no cargó a tiempo: %v\n", i, err)
			continue
		}

		if img, err := sheet.Element("img"); err == nil && img != nil {
			if srcAttr, _ := img.Attribute("src"); srcAttr != nil && *srcAttr != "" {
				fmt.Println("Imagen " + strconv.Itoa(i+1) + " procesada")
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
