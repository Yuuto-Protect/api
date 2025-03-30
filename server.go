package main

import (
	"api/routes"
	"api/utils"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func auth(c *gin.Context) {
	token, err := c.Cookie("token")
	// If there is no token or there was an error getting it
	if token == "" || err != nil {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}

	var user utils.User
	claims := jwt.MapClaims{}
	_, err = jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			return nil, errors.New("failed to get \"JWT_SECRET\" environment variable")
		}
		return []byte(secret), nil
	})

	// If there was an error during JWT parsing
	if err != nil {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}

	user.Id = claims["id"].(string)
	user.Email = claims["email"].(string)
	user.Username = claims["username"].(string)

	c.Set("user", user)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Failed to load .env file")
	}

	r := gin.Default()

	r.GET("/discord/callback", routes.DiscordCallback)

	// Auth middleware
	r.Use(auth)
	{
		r.POST("/subscription/checkout", routes.SubscriptionCheckout)
	}

	r.TrustedPlatform = gin.PlatformCloudflare
	r.Run()
}
