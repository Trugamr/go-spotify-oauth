package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
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
		panic(fmt.Errorf("fatal error loading env file: %s", err))
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
		// Get code from query string
		code := c.Query("code")
		if code == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "code is required"})
			return
		}

		// Get access token using code
		u := url.URL{
			Scheme: "https",
			Host:   "accounts.spotify.com",
			Path:   "/api/token",
		}

		q := url.Values{}
		q.Add("grant_type", "authorization_code")
		q.Add("code", code)
		q.Add("redirect_uri", env.SpotifyOAuthRedirectURI)

		body := bytes.NewBufferString(q.Encode())

		// Create request
		req, err := http.NewRequest(http.MethodPost, u.String(), body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Create basic auth header
		credentials := fmt.Sprintf("%s:%s", env.SpotifyClientID, env.SpotifyClientSecret)
		encoded := base64.StdEncoding.EncodeToString([]byte(credentials))

		req.Header = http.Header{
			"Content-Type":  []string{"application/x-www-form-urlencoded"},
			"Authorization": []string{fmt.Sprintf("Basic %s", encoded)},
		}

		// Create client and send request
		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer res.Body.Close()

		// Check response status code
		if res.StatusCode != http.StatusOK {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get access token"})
			return
		}

		var tokens struct {
			AccessToken  string `json:"access_token"`
			ExpiresIn    int    `json:"expires_in"`
			RefreshToken string `json:"refresh_token"`
			Scope        string `json:"scope"`
			TokenType    string `json:"token_type"`
		}
		decoder := json.NewDecoder(res.Body)
		if err := decoder.Decode(&tokens); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Get user profile using access token
		u = url.URL{
			Scheme: "https",
			Host:   "api.spotify.com",
			Path:   "/v1/me",
		}

		// Create request
		req, err = http.NewRequest(http.MethodGet, u.String(), nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		req.Header = http.Header{
			"Authorization": []string{fmt.Sprintf("Bearer %s", tokens.AccessToken)},
		}

		// Send request
		res, err = client.Do(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer res.Body.Close()

		// Check response status code
		if res.StatusCode != http.StatusOK {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get profile"})
			return
		}

		var profile struct {
			DisplayName string `json:"display_name"`
		}
		decoder = json.NewDecoder(res.Body)
		if err := decoder.Decode(&profile); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// TODO: Sign JWT token with user profile and save it in cookie

		c.HTML(http.StatusOK, "callback.tmpl", gin.H{
			"name": profile.DisplayName,
		})
	})

	app.Run(":8080")
}
