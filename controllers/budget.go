package controllers

import (
	"github.com/gin-gonic/gin"
	"log"
	"production/database"
	"production/models"
)

func GetBudget(c *gin.Context) {
	var budget models.Budget
	if err := database.DB.First(&budget).Error; err != nil {
		log.Printf("Ошибка при получении бюджета: %v", err)
		c.JSON(500, gin.H{"success": false, "error": "Не удалось получить бюджет"})
		return
	}

	c.JSON(200, gin.H{"success": true, "total_amount": budget.TotalAmount})
}

func BudgetList(c *gin.Context) {
	// Fetch все записи
	var budget []models.Budget
	if err := database.DB.Find(&budget).Error; err != nil {
		log.Printf("Ошибка при получении бюджета: %v", err)
		c.JSON(500, gin.H{"error": "Не удалось получить бюджет"})
		return
	}

	// Передача данных в шаблон
	c.HTML(200, "budget.html", gin.H{
		"budget": budget,
	})
}
