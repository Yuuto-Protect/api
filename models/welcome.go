package models

type Welcome struct {
	GuildId         string  `gorm:"primaryKey;column:guildId"`
	ChannelId       *string `gorm:"column:channelId"`
	Message         *string
	DmMessage       *string
	SendDm          bool
	SendChannel     bool
	AutoRoleMembers []string `gorm:"type:varchar(255)[];serializer:json;column:autoRoleMembers"`
	AutoRoleBots    []string `gorm:"type:varchar(255)[];serializer:json;column:autoRoleBots"`
}
