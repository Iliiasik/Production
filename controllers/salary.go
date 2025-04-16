package controllers

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"production/database"
	"production/models"
	"production/utils"
	"strconv"
	"time"
)

func ShowSalariesPage(c *gin.Context) {
	currentYear := time.Now().Year()

	var years []int
	for i := 0; i < 10; i++ {
		years = append(years, currentYear-i)
	}

	months := []struct {
		Value int
		Name  string
	}{
		{1, "Январь"}, {2, "Февраль"}, {3, "Март"}, {4, "Апрель"},
		{5, "Май"}, {6, "Июнь"}, {7, "Июль"}, {8, "Август"},
		{9, "Сентябрь"}, {10, "Октябрь"}, {11, "Ноябрь"}, {12, "Декабрь"},
	}

	c.HTML(200, "salaries.html", gin.H{
		"years":  years,
		"months": months,
	})
}

func GetSalaryByDate(c *gin.Context) {
	yearStr := c.Param("year")
	monthStr := c.Param("month")

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный год"})
		return
	}
	month, err := strconv.Atoi(monthStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный месяц"})
		return
	}

	var records []models.SalaryRecord
	if err := database.DB.Preload("Employee").
		Where("year = ? AND month = ?", year, month).
		Find(&records).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении данных"})
		return
	}

	c.JSON(http.StatusOK, records)
}

func CalculateSalary(c *gin.Context) {
	yearStr := c.Param("year")
	monthStr := c.Param("month")

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный год"})
		return
	}
	month, err := strconv.Atoi(monthStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный месяц"})
		return
	}

	if err := database.DB.Exec("CALL calculatesalary(?, ?)", year, month).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при расчёте зарплат"})
		return
	}

	c.Status(http.StatusOK)
}

func EditSalary(c *gin.Context) {
	type RequestBody struct {
		TotalSalary float64 `json:"total_salary"`
	}

	var body RequestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		log.Printf("Ошибка парсинга тела запроса: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Неверный формат запроса"})
		return
	}

	id := c.Param("id")

	var record models.SalaryRecord
	if err := database.DB.Where("id = ? AND is_paid = FALSE", id).First(&record).Error; err != nil {
		log.Printf("Запись с ID %s не найдена или уже выплачена: %v", id, err)
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "Запись не найдена или уже выплачена"})
		return
	}

	if err := database.DB.Model(&record).Update("total_salary", body.TotalSalary).Error; err != nil {
		log.Printf("Ошибка при обновлении зарплаты для ID %s: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Ошибка при обновлении"})
		return
	}

	log.Printf("Зарплата успешно обновлена для ID %s: новая сумма %.2f", id, body.TotalSalary)
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func PaySalaries(c *gin.Context) {
	yearStr := c.Param("year")
	monthStr := c.Param("month")

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный год"})
		return
	}
	month, err := strconv.Atoi(monthStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный месяц"})
		return
	}

	// Выполнение процедуры выплаты зарплат
	if err := database.DB.Exec("CALL pay_salaries(?, ?)", year, month).Error; err != nil {
		// Логируем полную ошибку с файлом и строкой
		log.Printf("%s: %v", utils.CallerLocation(1), err)

		// Возвращаем сообщение из ошибки (например, RAISE EXCEPTION)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": utils.ParseSQLErrorMessage(err.Error()),
		})
		return
	}

	c.Status(http.StatusOK)
}

func GetUnpaidSalariesTotal(c *gin.Context) {
	yearStr := c.Param("year")
	monthStr := c.Param("month")

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный год"})
		return
	}
	month, err := strconv.Atoi(monthStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный месяц"})
		return
	}

	// Получение суммы невыплаченных зарплат
	var total float64
	if err := database.DB.Raw("SELECT get_unpaid_salary_sum(?, ?)", year, month).Scan(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении суммы зарплат"})
		return
	}

	// Отправляем сумму невыплаченных зарплат на фронт
	c.JSON(http.StatusOK, gin.H{"total": total})
}
