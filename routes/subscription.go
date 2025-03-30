package routes

import (
	"api/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/checkout/session"
	"os"
)

type SubscriptionCheckoutBody struct {
	Plan string `json:"plan" binding:"required"`
	// user's Discord bot token
	Token string `json:"token" binding:"required"`
}

func SubscriptionCheckout(c *gin.Context) {
	plans := map[string]string{
		"premium": os.Getenv("STRIPE_PREMIUM_PLAN_ID"),
	}

	var body SubscriptionCheckoutBody
	err := c.BindJSON(&body)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	// Check if the plan is valid
	planID, ok := plans[body.Plan]
	if !ok {
		c.JSON(400, gin.H{"error": fmt.Sprintf("Invalid plan: %s", body.Plan)})
		return
	}

	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
	user := c.MustGet("user").(utils.User)

	// Create a new Stripe Checkout checkoutSession
	params := &stripe.CheckoutSessionParams{
		Mode:       stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		SuccessURL: stripe.String("https://bots.yuuto.dev/checkout/success"),
		CancelURL:  stripe.String("https://bots.yuuto.dev/"),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(planID),
				Quantity: stripe.Int64(1),
			},
		},
		CustomerEmail: stripe.String(user.Email),
		Metadata: map[string]string{
			"discord_token":   body.Token,
			"discord_user_id": user.Id,
		},
	}
	checkoutSession, err := session.New(params)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create a Stripe checkout session"})
		return
	}

	url := checkoutSession.URL
	c.JSON(200, gin.H{"url": url})
}
