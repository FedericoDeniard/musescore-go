package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/FedericoDeniard/musescore-go/src/constants"
	scrap "github.com/FedericoDeniard/musescore-go/src/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

type ScrapRequest struct {
	URL string `json:"url"`
}

func main() {
	router := gin.Default()

	router.Static("/assets", "./src/static/frontend/dist/assets")
	router.GET("/", func(c *gin.Context) {
		c.File("./src/static/frontend/dist/index.html")
	})

	router.POST("/scrap", func(c *gin.Context) {

		var req ScrapRequest

		if err := c.BindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "Invalid JSON"})
			return
		}

		url := req.URL
		fmt.Println("URL recibida:", url)

		chromiumPath := "/usr/bin/chromium-browser"
		if constants.KEYS.ENVIROMENT == "production" {
			chromiumPath = "/usr/bin/chromium"
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		u := launcher.New().Bin(chromiumPath).Headless(true).Set("no-sandbox").MustLaunch()
		browser := rod.New().Context(ctx).ControlURL(u).MustConnect()

		pdfPath := scrap.Scrap(browser, url)
		fmt.Println("PDF Path:", pdfPath)
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filepath.Base(pdfPath)))
		c.Header("Content-Type", "application/pdf")

		c.File(pdfPath)
		os.Remove(pdfPath)

	})

	router.Run(":8000")
}
