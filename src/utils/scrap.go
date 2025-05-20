package scrap

import (
	"fmt"
	"sync"
	"time"

	"github.com/FedericoDeniard/musescore-go/src/utils/images"
	"github.com/go-rod/rod"
)

func Scrap(browser *rod.Browser, url string) string {
	fmt.Println("Scraping process started...")
	// Configurar el navegador
	defer browser.MustClose()

	// url := "https://musescore.com/willdsc/scores/5964065"
	// url := "https://musescore.com/ericfontainejazz/scores/5662210"
	// url := "https://musescore.com/user/2539321/scores/7347764"

	page := browser.MustPage(url)
	fmt.Println("Page created")
	defer page.MustClose()
	fmt.Println("Page loaded")
	page.MustSetViewport(1920, 1080, 1, false)
	fmt.Println("Viewport set")

	var wg sync.WaitGroup
	wg.Add(1)

	// var title string
	sheetsChannel := make(chan string)
	var sheetsSource []string
	var imagesPath []string

	// go func() { defer wg.Done(); title = getTitle(page) }()
	go func() { defer wg.Done(); getSheets(page, sheetsChannel) }()

	for sheet := range sheetsChannel {
		sheetsSource = append(sheetsSource, sheet)
		imagesPath = append(imagesPath, images.DownloadImage(sheet))
	}

	wg.Wait()
	// fmt.Printf("Title: %s \n", title)
	fmt.Printf("Partituras encontradas: %d \n", len(sheetsSource))

	fmt.Println(imagesPath)
	fmt.Printf("Partituras descargadas: %d \n%v", len(imagesPath), imagesPath)

	pngPaths, err := images.ConvertMultipleSvgToPng(imagesPath...)
	if err != nil {
		panic(err)
	}

	pdfPath, err := images.ConvertPngToPdf(pngPaths...)
	if err != nil {
		panic(err)
	}
	filesToDelete := append(imagesPath, pngPaths...)
	images.DeleteImages(filesToDelete...)
	fmt.Println("Process finished")
	return pdfPath
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

func getSheets(page *rod.Page, channel chan<- string) {
	fmt.Println("Getting sheets...")
	defer close(channel)
	jmuseScrollerComponent, err := page.Element("#jmuse-scroller-component")
	if err != nil {
		fmt.Println("No se encontró el componente jmuse-scroller-component:", err)
		return
	}

	jmuseScrollerComponent.MustEval(`() => this.scrollIntoView({ behavior: "smooth", block: "start", inline: "nearest" })`)

	page.MustWaitRequestIdle()
	sheets := page.MustElements(".EEnGW")
	fmt.Println("Sheets found: ", len(sheets))

	for i, sheet := range sheets {
		fmt.Printf("Procesando hoja %d...\n", i)

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
				fmt.Printf("Imagen %d src: %s\n", i, *srcAttr)
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
