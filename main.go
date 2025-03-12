package main

import (
	"github.com/gin-gonic/gin"
	"production/database"
	"production/routes"
)

func main() {
	database.InitDB()
	r := gin.Default()
	routes.RegisterRoutes(r)

	r.LoadHTMLGlob("templates/*")
	r.Static("/assets", "./public")
	err := r.Run(":8080")
	if err != nil {
		return
	}
}
