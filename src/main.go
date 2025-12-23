package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/FedericoDeniard/musescore-go/src/config"
	"github.com/FedericoDeniard/musescore-go/src/middleware"
	scrap "github.com/FedericoDeniard/musescore-go/src/utils"
	customErrors "github.com/FedericoDeniard/musescore-go/src/utils/error"
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

	protected := router.Group("/api")
	protected.Use(middleware.ValidateJWT())

	protected.POST("/scrap", func(c *gin.Context) {
		user, exists := middleware.GetUserFromContext(c)
		if !exists {
			customErrors.HandleError(c, &customErrors.HttpError{StatusCode: 500, Message: "Error obteniendo usuario"})
			return
		}
		var req ScrapRequest

		if err := c.BindJSON(&req); err != nil {
			customErrors.HandleError(c, &customErrors.HttpError{
				StatusCode: 400,
				Message:    "Formato de solicitud inválido",
			})
			return
		}

		url := req.URL
		fmt.Println("URL recibida:", url)
		fmt.Println("Usuario:", user.Username, user.Email)

		fmt.Println("Starting scrap for URL:", url)

		chromiumPath := "/usr/bin/chromium-browser"
		if config.KEYS.ENVIRONMENT == "production" {
			chromiumPath = "/usr/bin/chromium"
		}

		fmt.Println("Launching Chrome with path:", chromiumPath)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		u := launcher.New().Bin(chromiumPath).Headless(true).Set("no-sandbox").MustLaunch()
		fmt.Println("Chrome launched successfully")
		browser := rod.New().Context(ctx).ControlURL(u).MustConnect()
		fmt.Println("Browser connected")

		pdfPath, httpError := scrap.Scrap(browser, url)
		if httpError != nil {
			customErrors.HandleError(c, httpError)
			return
		}
		fmt.Println("Scrap completed, PDF path:", pdfPath)
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filepath.Base(pdfPath)))
		c.Header("Content-Type", "application/pdf")

		c.File(pdfPath)
		os.Remove(pdfPath)
	})

	router.Run(":8000")
}
