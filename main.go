package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	app := gin.Default()

	// Load template files
	app.LoadHTMLGlob("templates/*")

	app.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{})
	})

	app.POST("/api/auth/login/spotify", func(c *gin.Context) {
		c.JSON(http.StatusNotImplemented, gin.H{
			"error": "not implemented",
		})
	})

	app.GET("/api/auth/spotify/callback", func(c *gin.Context) {
		c.JSON(http.StatusNotImplemented, gin.H{
			"error": "not implemented",
		})
	})

	app.Run(":8080")
}
