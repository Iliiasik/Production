package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"production/database"
	"production/models"
	"production/utils"
	"strconv"
	"time"
)

func ListCredits(c *gin.Context) {
	// Fetch все записи
	var credits []models.Credit
	if err := database.DB.Find(&credits).Error; err != nil {
		log.Printf("Ошибка при получении списка кредитов: %v", err)
		c.JSON(500, gin.H{"error": "Не удалось получить список кредитов"})
		return
	}

	// Передача данных в шаблон
	c.HTML(200, "credits.html", gin.H{
		"credits": credits,
	})
}

type createCreditInput struct {
	Amount      float64 `json:"amount"`
	StartDate   string  `json:"start_date"`
	TermYears   int     `json:"term_years"`
	AnnualRate  float64 `json:"annual_rate"`
	PenaltyRate float64 `json:"penalty_rate"`
}

// CreateCredit — HTTP‑хендлер для создания кредита через процедуру
func CreateCredit(c *gin.Context) {
	var input createCreditInput
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("Ошибка при привязке данных для создания кредита: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные"})
		return
	}

	// Парсим дату вручную
	startDate, err := time.Parse("2006-01-02", input.StartDate)
	if err != nil {
		log.Printf("Ошибка парсинга даты: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Неверный формат даты"})
		return
	}

	if err := callCreateCreditProcedure(
		input.Amount,
		startDate,
		input.TermYears,
		input.AnnualRate,
		input.PenaltyRate,
	); err != nil {
		log.Printf("Ошибка при вызове процедуры create_credit_with_budget_update: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось создать кредит"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Кредит успешно создан",
	})
}

// callCreateCreditProcedure вызывает PostgreSQL‑процедуру create_credit
func callCreateCreditProcedure(
	amount float64,
	startDate time.Time,
	termYears int,
	annualRate float64,
	penaltyRate float64,
) error {
	// Форматируем дату как YYYY‑MM‑DD, т.к. процедура ожидает DATE
	query := fmt.Sprintf(
		"CALL create_credit(%f, '%s', %d, %f, %f);",
		amount, startDate.Format("2006-01-02"), termYears, annualRate, penaltyRate,
	)

	// Выполняем запрос
	if err := database.DB.Exec(query).Error; err != nil {
		return fmt.Errorf("не удалось выполнить процедуру create_credit_with_budget_update: %v", err)
	}
	return nil
}

func ShowPaymentsPage(c *gin.Context) {
	id := c.Param("id")

	// — получаем кредит
	var credit models.Credit
	if err := database.DB.First(&credit, id).Error; err != nil {
		c.String(http.StatusInternalServerError, "Ошибка при получении кредита")
		return
	}

	// — получаем список выплат
	var payments []models.CreditPayment
	if err := database.DB.
		Where("credit_id = ?", id).
		Order("month_number ASC").
		Find(&payments).Error; err != nil {
		c.String(http.StatusInternalServerError, "Ошибка при получении выплат")
		return
	}

	// — получаем агрегаты
	var agg models.PaymentsAggregate
	if err := database.DB.
		Raw("SELECT * FROM sum_credit_payments(?)", id).
		Scan(&agg).Error; err != nil {
		log.Printf("Ошибка при суммировании выплат: %v", err)
		c.String(http.StatusInternalServerError, "Ошибка при подсчёте итогов")
		return
	}

	// — отрисовываем шаблон
	c.HTML(http.StatusOK, "credit-payments.html", gin.H{
		"credit":     credit,
		"payments":   payments,
		"aggregates": agg,
	})
}

func PayCredit(c *gin.Context) {
	creditIDStr := c.Param("id")
	creditID, err := strconv.Atoi(creditIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID кредита"})
		return
	}

	var req struct {
		Date string `json:"date"`
	}

	if err := c.BindJSON(&req); err != nil || req.Date == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Дата обязательна"})
		return
	}

	if err := database.DB.Exec("CALL credit_payment(?, ?)", creditID, req.Date).Error; err != nil {
		log.Printf("%s: %v", utils.CallerLocation(1), err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": utils.ParseSQLErrorMessage(err.Error()),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}
