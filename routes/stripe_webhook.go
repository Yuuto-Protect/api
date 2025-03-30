package routes

import (
	"api/models"
	"api/utils"
	"context"
	"encoding/json"
	"errors"
	"github.com/carlmjohnson/requests"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/webhook"
	"io"
	"net/http"
	"os"
)

var webhookSecret = os.Getenv("STRIPE_WEBHOOK_SECRET")

func createContainer(cli *client.Client, env []string) (string, error) {
	image := "ghcr.io/yuuto-protect/bot:latest"
	name, _ := uuid.NewRandom()
	resp, err := cli.ContainerCreate(context.Background(), &container.Config{
		Image: image,
		Env:   env,
	}, &container.HostConfig{
		NetworkMode: "yuuto_network",
		Binds:       []string{"/var/run/docker.sock:/var/run/docker.sock"},
		// Auto remove the container when it stops
		AutoRemove: true,
	}, nil, nil, name.String())
	if err == nil {
		return "", errors.New("failed to create the container")
	}

	err = cli.ContainerStart(context.Background(), resp.ID, container.StartOptions{})
	if err == nil {
		return "", errors.New("failed to start the container")
	}

	return resp.ID, nil
}

func StripeWebhook(c *gin.Context) {
	payload, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	signature := c.GetHeader("Stripe-Signature")
	event, err := webhook.ConstructEvent(payload, signature, webhookSecret)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid signature"})
	}

	switch event.Type {
	case "checkout.session.completed":
		var session stripe.CheckoutSession
		err := json.Unmarshal(event.Data.Raw, &session)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session data"})
			return
		}
		metadata := session.Metadata
		discordToken := metadata["discord_token"]
		discordUserId := metadata["discord_user_id"]

		var botInfoResp utils.User
		err = requests.
			URL("https://discord.com/api/v10/users/@me").
			Header("Authorization", "Bot "+discordToken).
			ToJSON(&botInfoResp).
			Fetch(context.Background())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create an instance"})
			return
		}

		db := utils.DbConnect()
		var instance models.Instances
		db.
			Where(models.Instances{Id: botInfoResp.Id}).
			Attrs(models.Instances{
				Token:   discordToken,
				Plan:    "premium",
				BuyerId: discordUserId,
				Name:    botInfoResp.Username,
			}).
			FirstOrCreate(&instance)

		cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		env := []string{
			"BOT_TOKEN=" + discordToken,
			"DB_HOST=" + os.Getenv("DB_HOST"),
			"DB_USER=" + os.Getenv("DB_USERNAME"),
			"DB_PASSWORD=" + os.Getenv("DB_PASSWORD"),
			"DB_NAME=" + os.Getenv("DB_NAME"),
			"GHCR_USERNAME" + os.Getenv("GHCR_USERNAME"),
			"GHCR_TOKEN" + os.Getenv("GHCR_TOKEN"),
		}
		containerId, err := createContainer(cli, env)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create an instance"})
			return
		}
		db.
			Where(&instance).
			Update("containerId", containerId)
	}

	c.Writer.WriteHeader(http.StatusNoContent)
}
