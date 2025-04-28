package controllers

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/gin-gonic/gin"
	"log"
	"production/database"
	"production/models"
)

const (
	budgetFetchErrorLog   = "Ошибка при получении бюджета: %v"
	budgetFetchErrorReply = "Не удалось получить бюджет"
)

func GetBudgetRow(c *gin.Context) {
	id := c.Param("id")
	var budget models.Budget

	if err := database.DB.First(&budget, id).Error; err != nil {
		log.Printf("Ошибка при получении бюджета с ID %s: %v", id, err)
		c.JSON(500, gin.H{"error": "Не удалось получить запись бюджета"})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"budget":  budget,
	})
}

func GetBudget(c *gin.Context) {
	var budget models.Budget
	if err := database.DB.First(&budget).Error; err != nil {
		log.Printf(budgetFetchErrorLog, err)
		c.JSON(500, gin.H{"success": false, "error": budgetFetchErrorReply})
		return
	}

	c.JSON(200, gin.H{"success": true, "total_amount": budget.TotalAmount})
}

func GetMarkup(c *gin.Context) {
	var budget models.Budget
	if err := database.DB.First(&budget).Error; err != nil {
		log.Printf(budgetFetchErrorLog, err)
		c.JSON(500, gin.H{"success": false, "error": budgetFetchErrorReply})
		return
	}

	c.JSON(200, gin.H{"success": true, "markup": budget.Markup})
}

func BudgetList(c *gin.Context) {
	var budgets []models.Budget
	if err := database.DB.Find(&budgets).Error; err != nil {
		log.Printf("Ошибка при получении бюджета: %v", err)
		c.JSON(500, gin.H{"error": "Не удалось получить бюджет"})
		return
	}

	type BudgetView struct {
		ID                   uint
		TotalAmount          float64
		TotalAmountFormatted string
		Markup               float64
		SalaryBonus          float64
	}

	var viewData []BudgetView
	for _, b := range budgets {
		// Пример: "KGS 123 456.79"
		formatted := fmt.Sprintf("KGS %s", humanize.FormatFloat("# ###.##", b.TotalAmount))

		viewData = append(viewData, BudgetView{
			ID:                   b.ID,
			TotalAmount:          b.TotalAmount,
			TotalAmountFormatted: formatted,
			Markup:               b.Markup,
			SalaryBonus:          b.SalaryBonus,
		})
	}

	c.HTML(200, "budget.html", gin.H{
		"budget": viewData,
	})
}
func UpdateBudget(c *gin.Context) {
	id := c.Param("id")
	var budget models.Budget

	// Получаем данные из тела запроса
	if err := c.ShouldBindJSON(&budget); err != nil {
		c.JSON(400, gin.H{"error": "Некорректные данные"})
		return
	}

	// Ищем запись по ID
	var existingBudget models.Budget
	if err := database.DB.First(&existingBudget, id).Error; err != nil {
		log.Printf("Ошибка при получении бюджета с ID %s: %v", id, err)
		c.JSON(404, gin.H{"error": "Не удалось найти запись для обновления"})
		return
	}

	// Логируем полученные данные
	log.Printf("Обновление бюджета ID=%s: TotalAmount=%.2f, Markup=%.2f, SalaryBonus=%.2f",
		id, budget.TotalAmount, budget.Markup, budget.SalaryBonus)

	// Создаем мапу с обновляемыми значениями
	updateData := map[string]interface{}{
		"total_amount": budget.TotalAmount,
		"markup":       budget.Markup,
		"salary_bonus": budget.SalaryBonus,
	}

	// Выполняем обновление
	if err := database.DB.Model(&existingBudget).Updates(updateData).Error; err != nil {
		log.Printf("Ошибка при обновлении бюджета с ID %s: %v", id, err)
		c.JSON(500, gin.H{"error": "Не удалось обновить запись"})
		return
	}

	// Возвращаем успешный ответ
	c.JSON(200, gin.H{"success": true})
}
