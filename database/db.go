package database

import (
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"production/models"
)

var DB *gorm.DB

func InitDB() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")

	if user == "" || password == "" || dbName == "" || host == "" || port == "" {
		log.Fatal("One or more environment variables are missing")
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, dbName, port,
	)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	log.Println("Connection success")

	if DB == nil {
		log.Fatal("Database connection is nil")
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("Error getting DB instance: %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}
	log.Println("Database ping success")

	err = DB.AutoMigrate(
		&models.Unit{},
		&models.RawMaterial{},
		&models.FinishedGood{},
		&models.Position{},
		&models.Employee{},
		&models.SalaryRecord{},
		&models.Ingredient{},
		&models.Budget{},
		&models.RawMaterialPurchase{},
		&models.ProductSale{},
		&models.ProductProduction{},
	)
	if err != nil {
		log.Fatalf("Error running migrations: %v", err)
	}

	log.Println("Migrations completed successfully")
}
