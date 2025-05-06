package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"production/database"
	"strings"
	"time"
)

type SaleRecord struct {
	ProductName  string    `json:"product_name"`
	Quantity     int       `json:"quantity"`
	TotalAmount  float64   `json:"total_amount"`
	SaleDate     time.Time `json:"sale_date"`
	EmployeeName string    `json:"employee_name"`
}

type SaleSummary struct {
	ProductName   string  `json:"product_name"`
	TotalQuantity int     `json:"total_quantity"`
	TotalSales    float64 `json:"total_sales"`
}

func SalesReportHandler(c *gin.Context) {
	type DateRange struct {
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
	}

	var dateRange DateRange
	if err := c.BindJSON(&dateRange); err != nil {
		c.String(http.StatusBadRequest, "Неверный формат тела запроса")
		return
	}

	startDate, err := time.Parse("2006-01-02", dateRange.StartDate)
	endDate, err2 := time.Parse("2006-01-02", dateRange.EndDate)
	if err != nil || err2 != nil {
		c.String(http.StatusBadRequest, "Неверный формат дат")
		return
	}

	var reportRows []SaleRecord
	if err := database.DB.
		Raw("SELECT * FROM sales_report(?, ?)", startDate, endDate).
		Scan(&reportRows).Error; err != nil {
		c.String(http.StatusInternalServerError, "Ошибка отчета продаж")
		return
	}

	var summaryRows []SaleSummary
	if err := database.DB.
		Raw("SELECT * FROM sales_summary(?, ?)", startDate, endDate).
		Scan(&summaryRows).Error; err != nil {
		c.String(http.StatusInternalServerError, "Ошибка итогов отчета")
		return
	}

	var html strings.Builder

	// Заголовок отчёта
	html.WriteString(fmt.Sprintf(`
	<div class="report-heading">
		<h3>Отчёт по продажам за %s – %s</h2>
	</div>`,
		startDate.Format("02.01.2006"), endDate.Format("02.01.2006")))

	// Таблица подробного отчета
	html.WriteString(`
	<div class="table">
		<div class="table-header">
			<div class="header__item">Дата</div>
			<div class="header__item">Продукт</div>
			<div class="header__item">Количество</div>
			<div class="header__item">Сумма</div>
			<div class="header__item">Сотрудник</div>
		</div>
		<div class="table-body">`)

	if len(reportRows) == 0 {
		html.WriteString(`<div class="table-row"><div class="table-data" colspan="5">Данные не найдены.</div></div>`)
	} else {
		for _, row := range reportRows {
			html.WriteString(fmt.Sprintf(`
			<div class="table-row">
				<div class="table-data">%s</div>
				<div class="table-data">%s</div>
				<div class="table-data">%d</div>
				<div class="table-data">%.2f</div>
				<div class="table-data">%s</div>
			</div>`,
				row.SaleDate.Format("02.01.2006"),
				row.ProductName,
				row.Quantity,
				row.TotalAmount,
				row.EmployeeName))
		}
	}
	html.WriteString(`</div></div>`)

	// Подзаголовок общей сводки
	html.WriteString(`
	<div class="report-heading">
		<h3>Общая сводка</h3>
	</div>`)

	// Таблица итогов
	html.WriteString(`
	<div class="table">
		<div class="table-header">
			<div class="header__item">Продукт</div>
			<div class="header__item">Общее количество</div>
			<div class="header__item">Общая сумма</div>
		</div>
		<div class="table-body">`)

	if len(summaryRows) == 0 {
		html.WriteString(`<div class="table-row"><div class="table-data" colspan="3">Итоги не найдены.</div></div>`)
	} else {
		for _, row := range summaryRows {
			html.WriteString(fmt.Sprintf(`
			<div class="table-row">
				<div class="table-data">%s</div>
				<div class="table-data">%d</div>
				<div class="table-data">%.2f</div>
			</div>`,
				row.ProductName,
				row.TotalQuantity,
				row.TotalSales))
		}
	}
	html.WriteString(`</div></div>`)

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html.String()))
}
