package db

import (
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectPostgres() {
	dsn := "host=localhost user=postgres password=lion dbname=practiceDB port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	if err != nil {
		log.Fatal("Database connection unsuccessful, err")
	}
	if err := db.Exec("SET TIME ZONE 'UTC'").Error; err != nil {
		log.Fatal("Failed to set timezone:", err)
	}
	log.Println("Database connected successfully")
	DB = db
}
