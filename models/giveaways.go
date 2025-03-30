package models

import "time"

type Giveaways struct {
	GuildId      string    `gorm:"column:guildId"`
	ChannelId    string    `gorm:"column:channelId"`
	MessageId    string    `gorm:"primaryKey;column:messageId"`
	Winners      int       `gorm:"column:winners"`
	WinnersId    []string  `gorm:"type:varchar(255)[];serializer:json;column:winnersId"`
	Prize        string    `gorm:"column:prize"`
	Participants []string  `gorm:"type:varchar(255)[];serializer:json;column:participants"`
	Duration     int       `gorm:"column:duration"`
	Ended        bool      `gorm:"column:ended"`
	EndedAt      time.Time `gorm:"column:endedAt"`
	CreatedAt    time.Time `gorm:"column:createdAt"`
}
