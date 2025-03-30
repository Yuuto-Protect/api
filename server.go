package main

import (
	"api/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Failed to load .env file")
	}

	r := gin.Default()
	r.GET("/discord/callback", routes.DiscordCallback)

	r.TrustedPlatform = gin.PlatformCloudflare
	r.Run()
}
