package scrap

import (
	"fmt"
	"sync"

	"github.com/FedericoDeniard/musescore-go/src/utils/images"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

func Scrap() {
	fmt.Println("Scraping process started...")
	// Configurar el navegador
	url := launcher.New().MustLaunch()
	browser := rod.New().ControlURL(url).MustConnect()
	defer browser.MustClose()

	// link := "https://musescore.com/willdsc/scores/5964065"
	link := "https://musescore.com/ericfontainejazz/scores/5662210"

	page := browser.MustPage(link).MustWaitLoad()
	page.MustSetViewport(1920, 1080, 1, false)

	var wg sync.WaitGroup
	wg.Add(2)

	var title string
	var sheetsSource []string

	go func() { defer wg.Done(); title = getTitle(page) }()
	go func() { defer wg.Done(); sheetsSource = getSheets(page) }()

	wg.Wait()
	fmt.Println(title)
	fmt.Printf("Partituras encontradas: %d \n%v", len(sheetsSource), sheetsSource)

	imagesPath, err := images.DownloadImage(sheetsSource)
	if err != nil {
		panic(err)
	}

	fmt.Println(imagesPath)
	fmt.Printf("Partituras descargadas: %d \n%v", len(imagesPath), imagesPath)

	_, err = images.ConvertMusescoreSVGsToPDF(imagesPath, title)
	if err != nil {
		panic(err)
	}

}

func getTitle(page *rod.Page) string {
	title := "musescore"
	asideContainer := page.MustElement("#aside-container-unique")
	if asideContainer != nil {
		titleElement := asideContainer.MustElement(".nFRPI")
		if titleElement != nil {
			span := titleElement.MustElement("span")
			if span != nil {
				title = span.MustText()
			}
		}
	}
	return title
}

func getSheets(page *rod.Page) []string {
	jmuseScrollerComponent := page.MustElement("#jmuse-scroller-component")
	if jmuseScrollerComponent != nil {
		jmuseScrollerComponent.MustScrollIntoView()
		page.MustWaitRequestIdle()

		sheets := page.MustElements(".EEnGW")
		var sheetsUrl []string

		for _, sheet := range sheets {
			sheet.MustScrollIntoView()
			page.MustWaitRequestIdle()()

			image := sheet.MustElement("img")

			if image != nil {
				if src, _ := image.Attribute("src"); src != nil {
					sheetsUrl = append(sheetsUrl, *src)
				}
			}
		}
		return sheetsUrl
	}
	panic("Error al obtener las hojas")
}
