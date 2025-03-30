package models

import "time"

type Instances struct {
	Id          string    `gorm:"primaryKey;column:id"`
	Name        string    `gorm:"column:name"`
	Token       string    `gorm:"column:token"`
	BuyerId     string    `gorm:"column:buyerId"`
	Plan        string    `gorm:"column:plan"`
	ContainerId string    `gorm:"column:containerId"`
	CreatedAt   time.Time `gorm:"column:createdAt"`
}
