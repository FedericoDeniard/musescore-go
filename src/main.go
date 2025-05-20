package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	// scrap.Scrap()
	router := gin.Default()

	router.Static("/assets", "./src/static/frontend/dist/assets")
	router.GET("/", func(c *gin.Context) {
		c.File("./src/static/frontend/dist/index.html")
	})

	// Iniciar el servidor
	router.Run(":8080")
}
