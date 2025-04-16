package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"production/database"
	"production/models"
	"time"
)

func ListSales(c *gin.Context) {
	var sales []models.ProductSale
	if err := database.DB.Find(&sales).Error; err != nil {
		log.Printf("Ошибка при получении проданной продукции: %v", err)
		c.JSON(500, gin.H{"error": "Не удалось получить проданную продукцию"})
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

	type SalesWithDetails struct {
		ID          uint
		Product     string
		Quantity    float64
		TotalAmount float64
		SaleDate    time.Time
		Employee    string
	}

	var saleWithDetails []SalesWithDetails
	for _, sale := range sales {
		saleWithDetails = append(saleWithDetails, SalesWithDetails{
			ID:          sale.ID,
			Product:     finishedGoodsMap[sale.ProductID],
			Quantity:    sale.Quantity,
			TotalAmount: sale.TotalAmount,
			SaleDate:    sale.SaleDate,
			Employee:    employeesMap[sale.EmployeeID],
		})
	}

	c.HTML(200, "sales.html", gin.H{
		"sales": saleWithDetails,
	})
}

func MakeSale(c *gin.Context) {
	var sale models.ProductSale
	if err := c.ShouldBindJSON(&sale); err != nil {
		log.Printf("Ошибка при привязке данных: %v", err)
		c.JSON(400, gin.H{"error": "Неверные данные"})
		return
	}

	productID := sale.ProductID
	quantity := sale.Quantity
	saleDate := sale.SaleDate
	employeeID := sale.EmployeeID

	if err := callMakeSaleProcedure(productID, quantity, saleDate, employeeID); err != nil {
		log.Printf("Ошибка при вызове процедуры make_sale: %v", err)
		c.JSON(500, gin.H{"error": "Не удалось выполнить продажу"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Продажа успешно выполнена"})
}

// Функция для вызова процедуры make_sale
func callMakeSaleProcedure(productID uint, quantity float64, saleDate time.Time, employeeID uint) error {
	// Создаём SQL-запрос для вызова процедуры
	query := fmt.Sprintf("CALL make_sale(%d, %f, '%s', %d);", productID, quantity, saleDate.Format("2006-01-02 15:04:05"), employeeID)

	// Выполняем запрос
	if err := database.DB.Exec(query).Error; err != nil {
		return fmt.Errorf("не удалось выполнить процедуру make_sale: %v", err)
	}
	return nil
}
