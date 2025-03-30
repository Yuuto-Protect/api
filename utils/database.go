package utils

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"log"
	"os"
)

func DbConnect() *gorm.DB {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	database := os.Getenv("DB_NAME")

	if host == "" || user == "" || password == "" || database == "" {
		log.Fatal("Cannot connect to database, some environment variables are not set")
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s", host, user, password, database,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			NoLowerCase:   true,
			SingularTable: true,
		},
	})
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to connect to database: %s", err))
	}

	return db
}
