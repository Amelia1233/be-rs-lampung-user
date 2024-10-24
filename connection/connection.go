package connection

import (
	"fmt"
	"log"
	"be-rs-lampung-user/entity"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connection() *gorm.DB {
	errEnv := godotenv.Load()
	if errEnv != nil {
		log.Fatal("Failed Load Env")
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", 
    os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed Connect To Database: %v", err)
	}

	db.AutoMigrate(&entity.Login{})
	return db
}
