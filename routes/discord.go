package routes

import (
	"context"
	"encoding/base64"
	"github.com/carlmjohnson/requests"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"net/url"
	"os"
)

type OAuth2TokenResponse struct {
	AccessToken string `json:"access_token"`
	Message     string `json:"message"`
}

type DiscordUserResponse struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func DiscordCallback(c *gin.Context) {
	hub := sentrygin.GetHubFromContext(c)
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Missing code",
		})
		return
	}
	encodedState := c.DefaultQuery("state", base64.StdEncoding.EncodeToString([]byte("/manage")))
	var state []byte
	var err error
	if encodedState != "" {
		state, err = base64.StdEncoding.DecodeString(encodedState)
		if err != nil {
			// If we fail to decode the state, fallback to "/manage"
			state = []byte("/manage")
		}
	}

	clientId := os.Getenv("DISCORD_CLIENT_ID")
	clientSecret := os.Getenv("DISCORD_CLIENT_SECRET")
	redirectUri := os.Getenv("DISCORD_REDIRECT_URI")
	body := url.Values{}
	body.Set("code", code)
	body.Set("grant_type", "authorization_code")
	body.Set("redirect_uri", redirectUri)

	var tokenResp OAuth2TokenResponse
	err = requests.
		URL("https://discord.com/api/oauth2/token").
		ContentType("application/x-www-form-urlencoded").
		BasicAuth(clientId, clientSecret).
		BodyForm(body).
		ToJSON(&tokenResp).
		Fetch(context.Background())
	if err != nil {
		hub.CaptureException(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get an OAuth2 token"})
		return
	}

	var userResp DiscordUserResponse
	err = requests.
		URL("https://discord.com/api/v10/users/@me").
		Bearer(tokenResp.AccessToken).
		ToJSON(&userResp).
		Fetch(context.Background())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to fetch user info"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       userResp.Id,
		"username": userResp.Username,
		"email":    userResp.Email,
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.SetCookie("token", tokenString, 604800, "/", "", true, true)
	c.Redirect(302, string(state))
}
