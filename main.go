package main

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/trugamr/go-spotify-oauth/initializers"
)

func main() {
	// Load environment variables from .env file
	env, err := initializers.LoadEnv()
	if err != nil {
		fmt.Errorf("fatal error loading env file: %s", err)
	}

	app := gin.Default()

	// Load template files
	app.LoadHTMLGlob("templates/*")

	app.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{})
	})

	app.GET("/api/auth/login/spotify", func(c *gin.Context) {
		u := url.URL{
			Scheme: "https",
			Host:   "accounts.spotify.com",
			Path:   "/authorize",
		}

		q := url.Values{}
		q.Add("response_type", "code")
		q.Add("client_id", env.SpotifyClientID)
		q.Add("redirect_uri", env.SpotifyOAuthRedirectURI)
		q.Add("scope", "user-read-private user-read-email")

		u.RawQuery = q.Encode()

		c.Redirect(http.StatusTemporaryRedirect, u.String())
	})

	app.GET("/api/auth/spotify/callback", func(c *gin.Context) {
		c.JSON(http.StatusNotImplemented, gin.H{
			"error": "not implemented",
		})
	})

	app.Run(":8080")
}
