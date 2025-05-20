package main

import (
	"strconv"

	"github.com/FedericoDeniard/musescore-go/src/utils/images"
)

func main() {
	// scrap.Scrap()
	// router := gin.Default()

	// router.Static("/assets", "./src/static/frontend/dist/assets")
	// router.GET("/", func(c *gin.Context) {
	// 	c.File("./src/static/frontend/dist/index.html")
	// })

	// // Iniciar el servidor
	// router.Run(":8080")
	var pngPaths []string
	for i := 1; i <= 3; i++ {
		pngPath := images.ConvertSvgToPng("src/downloads/images/" + strconv.Itoa(i) + ".svg")
		pngPaths = append(pngPaths, pngPath)
	}

	images.ConvertPngToPdf(pngPaths...)
}
