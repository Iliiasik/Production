package controllers

import (
	"github.com/gin-gonic/gin"
	"log"
	"production/database"
	"production/models"
)

func GetEmployeesList(c *gin.Context) {
	var employees []models.Employee
	if err := database.DB.Find(&employees).Error; err != nil {
		log.Printf("Ошибка при получении сотрудников: %v", err)
		c.JSON(500, gin.H{"success": false, "error": "Не удалось получить сотрудников"})
		return
	}

	c.JSON(200, gin.H{"success": true, "employees": employees})
}
