package db

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitPsqlDB() (*gorm.DB, error) {
	dns := os.Getenv("POSTGRES_URI")
	log.Println("dns: ", dns)
	dialector := postgres.Open(dns)
	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}