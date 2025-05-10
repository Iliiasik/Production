package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
	"production/controllers"
	"production/database"
	"production/routes"
)

func main() {
	// Загружаем .env
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Устанавливаем JWT ключ
	controllers.JwtKey = []byte(os.Getenv("JWT_SECRET"))
	if len(controllers.JwtKey) == 0 {
		log.Fatal("JWT_SECRET is not set in .env")
	}

	// Инициализируем БД
	database.InitDB()

	// Настраиваем роутер
	r := gin.Default()
	routes.RegisterRoutes(r)

	r.LoadHTMLGlob("templates/*")
	r.Static("/assets", "./public")

	err = r.Run(":8080")
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
