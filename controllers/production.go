package controllers

import (
	"github.com/gin-gonic/gin"
	"log"
	"production/database"
	"production/models"
	"time"
)

func ListProductProduction(c *gin.Context) {
	var productions []models.ProductProduction
	if err := database.DB.Find(&productions).Error; err != nil {
		log.Printf("Ошибка при получении произведенной продукции: %v", err)
		c.JSON(500, gin.H{"error": "Не удалось получить произведенную продукцию"})
		return
	}

	var finishedGoods []models.FinishedGood
	if err := database.DB.Find(&finishedGoods).Error; err != nil {
		log.Printf("Ошибка при получении списка готовой продукции: %v", err)
		c.JSON(500, gin.H{"error": "Не удалось получить список готовой продукции"})
		return
	}

	finishedGoodsMap := make(map[uint]string)
	for _, finishedGood := range finishedGoods {
		finishedGoodsMap[finishedGood.ID] = finishedGood.Name
	}

	var employees []models.Employee
	if err := database.DB.Find(&employees).Error; err != nil {
		log.Printf("Ошибка при получении списка сотрудников: %v", err)
		c.JSON(500, gin.H{"error": "Не удалось получить список сотрудников"})
		return
	}

	employeesMap := make(map[uint]string)
	for _, employee := range employees {
		employeesMap[employee.ID] = employee.FullName
	}

	type ProductProductionWithDetails struct {
		ID             uint
		Product        string
		Quantity       float64
		ProductionDate time.Time
		Employee       string
	}

	var productProductionWithDetails []ProductProductionWithDetails
	for _, prod := range productions {
		productProductionWithDetails = append(productProductionWithDetails, ProductProductionWithDetails{
			ID:             prod.ID,
			Product:        finishedGoodsMap[prod.ProductID],
			Quantity:       prod.Quantity,
			ProductionDate: prod.ProductionDate,
			Employee:       employeesMap[prod.EmployeeID],
		})
	}

	c.HTML(200, "product-production.html", gin.H{
		"productProductions": productProductionWithDetails,
	})
}
