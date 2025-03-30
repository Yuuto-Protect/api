package routes

import (
	"api/models"
	"api/utils"
	"github.com/gin-gonic/gin"
)

type AutoRolesBody struct {
	Members []string `json:"members"`
	Bots    []string `json:"bots"`
}

func ManageGuildAutoRoles(c *gin.Context) {
	var body AutoRolesBody
	err := c.BindJSON(&body)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	if body.Members == nil || body.Bots == nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	guildId := c.Param("guildId")
	db := utils.DbConnect()
	welcome := models.Welcome{}
	db.
		Model(&welcome).
		Where(models.Welcome{GuildId: guildId}).
		Updates(models.Welcome{
			AutoRoleBots:    body.Bots,
			AutoRoleMembers: body.Members,
		})

	c.JSON(200, gin.H{
		"members": welcome.AutoRoleMembers,
		"bots":    welcome.AutoRoleBots,
	})
}
