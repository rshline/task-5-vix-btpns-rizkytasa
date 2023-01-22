package database

import  (
	"os"
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
	"gorm.io/driver/mysql"
	"github.com/rshline/task-5-vix-btpns-rizkytasa/models"

)

func SetupDB() *gorm.DB {
	//get environment variables
	godotenv.Load(".env")
	DB_HOST := os.Getenv("DB_HOST")
	DB_USER := os.Getenv("DB_USER")
	DB_PASSWORD := os.Getenv("DB_PASSWORD")
	DB_NAME := os.Getenv("DB_NAME")
	DB_PORT := os.Getenv("DB_PORT")

	// Connect to mysql database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", DB_USER, DB_PASSWORD, DB_HOST, DB_PORT, DB_NAME)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})

	if err != nil {
		panic(err)
	} else {
		fmt.Println("Connected to mysql database!")
	}
	
	// migrate database
	db.AutoMigrate(&models.User{}, &models.Photo{})

	return db
}