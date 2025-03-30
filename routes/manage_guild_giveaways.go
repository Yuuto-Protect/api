package routes

import (
	"api/models"
	"api/utils"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

func DeleteGiveaway(c *gin.Context) {
	giveawayId := c.Param("giveawayId")
	guildId := c.Param("guildId")

	db := utils.DbConnect()
	giveaway := db.Where(models.Giveaways{MessageId: giveawayId, GuildId: guildId})
	if errors.Is(giveaway.Error, gorm.ErrRecordNotFound) {
		c.JSON(404, gin.H{"error": "Unknown giveaway"})
	}

	c.Writer.WriteHeader(http.StatusNoContent)
}
