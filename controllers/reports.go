package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"production/database"
	"strings"
	"time"
)

// Каждая функция прописана отдельно. Необходимо для системы доступов и ролей.

// Отчеты по продажам

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

// Отчеты по производству

type ProductionRecord struct {
	ProductName    string    `json:"product_name"`
	Quantity       int       `json:"quantity"`
	ProductionDate time.Time `json:"production_date"`
	EmployeeName   string    `json:"employee_name"`
}
type ProductionSummary struct {
	ProductName   string `json:"product_name"`
	TotalQuantity int    `json:"total_quantity"`
}

func ProductionReportHandler(c *gin.Context) {
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

	// Отчет по производству
	var reportRows []ProductionRecord
	if err := database.DB.
		Raw("SELECT * FROM production_report(?, ?)", startDate, endDate).
		Scan(&reportRows).Error; err != nil {
		c.String(http.StatusInternalServerError, "Ошибка отчета по производству")
		return
	}

	// Итоги по производству
	var summaryRows []ProductionSummary
	if err := database.DB.
		Raw("SELECT * FROM production_summary(?, ?)", startDate, endDate).
		Scan(&summaryRows).Error; err != nil {
		c.String(http.StatusInternalServerError, "Ошибка итогов отчета по производству")
		return
	}

	var html strings.Builder

	// Заголовок отчёта
	html.WriteString(fmt.Sprintf(`
	<div class="report-heading">
		<h3>Отчёт по производству за %s – %s</h2>
	</div>`,
		startDate.Format("02.01.2006"), endDate.Format("02.01.2006")))

	// Таблица подробного отчета
	html.WriteString(`
	<div class="table">
		<div class="table-header">
			<div class="header__item">Дата производства</div>
			<div class="header__item">Продукт</div>
			<div class="header__item">Количество</div>
			<div class="header__item">Сотрудник</div>
		</div>
		<div class="table-body">`)

	if len(reportRows) == 0 {
		html.WriteString(`<div class="table-row"><div class="table-data" colspan="4">Данные не найдены.</div></div>`)
	} else {
		for _, row := range reportRows {
			html.WriteString(fmt.Sprintf(`
			<div class="table-row">
				<div class="table-data">%s</div>
				<div class="table-data">%s</div>
				<div class="table-data">%d</div>
				<div class="table-data">%s</div>
			</div>`,
				row.ProductionDate.Format("02.01.2006"),
				row.ProductName,
				row.Quantity,
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
		</div>
		<div class="table-body">`)

	if len(summaryRows) == 0 {
		html.WriteString(`<div class="table-row"><div class="table-data" colspan="2">Итоги не найдены.</div></div>`)
	} else {
		for _, row := range summaryRows {
			html.WriteString(fmt.Sprintf(`
			<div class="table-row">
				<div class="table-data">%s</div>
				<div class="table-data">%d</div>
			</div>`,
				row.ProductName,
				row.TotalQuantity))
		}
	}
	html.WriteString(`</div></div>`)

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html.String()))
}

// Отчеты по закупкам

type PurchaseRecord struct {
	RawMaterialName string    `json:"raw_material_name"`
	Quantity        float64   `json:"quantity"`
	TotalAmount     float64   `json:"total_amount"`
	PurchaseDate    time.Time `json:"purchase_date"`
	EmployeeName    string    `json:"employee_name"`
}
type PurchaseSummary struct {
	RawMaterialName string  `json:"raw_material_name"`
	TotalQuantity   float64 `json:"total_quantity"`
	TotalAmount     float64 `json:"total_amount"`
}

func PurchaseReportHandler(c *gin.Context) {
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

	var reportRows []PurchaseRecord
	if err := database.DB.
		Raw("SELECT * FROM raw_material_purchase_report(?, ?)", startDate, endDate).
		Scan(&reportRows).Error; err != nil {
		c.String(http.StatusInternalServerError, "Ошибка отчета закупок")
		return
	}

	var summaryRows []PurchaseSummary
	if err := database.DB.
		Raw("SELECT * FROM raw_material_purchase_summary(?, ?)", startDate, endDate).
		Scan(&summaryRows).Error; err != nil {
		c.String(http.StatusInternalServerError, "Ошибка итогов отчета")
		return
	}

	var html strings.Builder

	html.WriteString(fmt.Sprintf(`
	<div class="report-heading">
		<h3>Отчёт по закупкам сырья за %s – %s</h3>
	</div>`,
		startDate.Format("02.01.2006"), endDate.Format("02.01.2006")))

	html.WriteString(`
	<div class="table">
		<div class="table-header">
			<div class="header__item">Дата</div>
			<div class="header__item">Сырьё</div>
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
				<div class="table-data">%.2f</div>
				<div class="table-data">%.2f</div>
				<div class="table-data">%s</div>
			</div>`,
				row.PurchaseDate.Format("02.01.2006"),
				row.RawMaterialName,
				row.Quantity,
				row.TotalAmount,
				row.EmployeeName))
		}
	}
	html.WriteString(`</div></div>`)

	html.WriteString(`
	<div class="report-heading">
		<h3>Общая сводка</h3>
	</div>`)

	html.WriteString(`
	<div class="table">
		<div class="table-header">
			<div class="header__item">Сырьё</div>
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
				<div class="table-data">%.2f</div>
				<div class="table-data">%.2f</div>
			</div>`,
				row.RawMaterialName,
				row.TotalQuantity,
				row.TotalAmount))
		}
	}
	html.WriteString(`</div></div>`)

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html.String()))
}

// Отчеты по зарплатам

type SalaryRecord struct {
	EmployeeName       string  `json:"employee_name"`
	Year               int     `json:"year"`
	Month              int     `json:"month"`
	PurchaseCount      int     `json:"purchase_count"`
	ProductionCount    int     `json:"production_count"`
	SaleCount          int     `json:"sale_count"`
	TotalParticipation int     `json:"total_participation"`
	Bonus              float64 `json:"bonus"`
	TotalSalary        float64 `json:"total_salary"`
	IsPaid             bool    `json:"is_paid"`
}
type SalarySummary struct {
	EmployeeName string  `json:"employee_name"`
	MonthsCount  int     `json:"months_count"`
	TotalBonus   float64 `json:"total_bonus"`
	TotalSalary  float64 `json:"total_salary"`
	AvgBonus     float64 `json:"avg_bonus"`
	AvgSalary    float64 `json:"avg_salary"`
	PaidCount    int     `json:"paid_count"`
	UnpaidCount  int     `json:"unpaid_count"`
}

func SalaryReportHandler(c *gin.Context) {
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

	var reportRows []SalaryRecord
	if err := database.DB.
		Raw("SELECT * FROM salary_report(?, ?)", startDate, endDate).
		Scan(&reportRows).Error; err != nil {
		c.String(http.StatusInternalServerError, "Ошибка при получении данных зарплат")
		return
	}

	var summaryRows []SalarySummary
	if err := database.DB.
		Raw("SELECT * FROM salary_summary(?, ?)", startDate, endDate).
		Scan(&summaryRows).Error; err != nil {
		c.String(http.StatusInternalServerError, "Ошибка при получении сводки зарплат")
		return
	}

	var html strings.Builder

	html.WriteString(fmt.Sprintf(`
	<div class="report-heading">
		<h3>Отчёт по зарплатам за %s – %s</h3>
	</div>`,
		startDate.Format("02.01.2006"), endDate.Format("02.01.2006")))

	html.WriteString(`
	<div class="table">
		<div class="table-header">
			<div class="header__item">Сотрудник</div>
			<div class="header__item">Год</div>
			<div class="header__item">Месяц</div>
			<div class="header__item">Закупки</div>
			<div class="header__item">Производство</div>
			<div class="header__item">Продажи</div>
			<div class="header__item">Участия</div>
			<div class="header__item">Бонус</div>
			<div class="header__item">Зарплата</div>
			<div class="header__item">Выплачено</div>
		</div>
		<div class="table-body">`)

	if len(reportRows) == 0 {
		html.WriteString(`<div class="table-row"><div class="table-data" colspan="10">Данные не найдены.</div></div>`)
	} else {
		for _, row := range reportRows {
			html.WriteString(fmt.Sprintf(`
			<div class="table-row">
				<div class="table-data">%s</div>
				<div class="table-data">%d</div>
				<div class="table-data">%02d</div>
				<div class="table-data">%d</div>
				<div class="table-data">%d</div>
				<div class="table-data">%d</div>
				<div class="table-data">%d</div>
				<div class="table-data">%.2f</div>
				<div class="table-data">%.2f</div>
				<div class="table-data">%s</div>
			</div>`,
				row.EmployeeName,
				row.Year,
				row.Month,
				row.PurchaseCount,
				row.ProductionCount,
				row.SaleCount,
				row.TotalParticipation,
				row.Bonus,
				row.TotalSalary,
				map[bool]string{true: "Да", false: "Нет"}[row.IsPaid]))
		}
	}
	html.WriteString(`</div></div>`)

	html.WriteString(`
	<div class="report-heading">
		<h3>Сводка по зарплатам</h3>
	</div>`)

	html.WriteString(`
	<div class="table">
		<div class="table-header">
			<div class="header__item">Сотрудник</div>
			<div class="header__item">Месяцев</div>
			<div class="header__item">Сумма бонусов</div>
			<div class="header__item">Сумма зарплат</div>
			<div class="header__item">Средний бонус</div>
			<div class="header__item">Средняя зарплата</div>
			<div class="header__item">Выплачено</div>
			<div class="header__item">Не выплачено</div>
		</div>
		<div class="table-body">`)

	if len(summaryRows) == 0 {
		html.WriteString(`<div class="table-row"><div class="table-data" colspan="8">Сводка не найдена.</div></div>`)
	} else {
		for _, row := range summaryRows {
			html.WriteString(fmt.Sprintf(`
			<div class="table-row">
				<div class="table-data">%s</div>
				<div class="table-data">%d</div>
				<div class="table-data">%.2f</div>
				<div class="table-data">%.2f</div>
				<div class="table-data">%.2f</div>
				<div class="table-data">%.2f</div>
				<div class="table-data">%d</div>
				<div class="table-data">%d</div>
			</div>`,
				row.EmployeeName,
				row.MonthsCount,
				row.TotalBonus,
				row.TotalSalary,
				row.AvgBonus,
				row.AvgSalary,
				row.PaidCount,
				row.UnpaidCount))
		}
	}
	html.WriteString(`</div></div>`)

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html.String()))
}

// Отчеты по выплатам по кредиту

type CreditPaymentRecord struct {
	CreditID         int64     `json:"credit_id"`
	PaymentDate      time.Time `json:"payment_date"`
	MonthNumber      int       `json:"month_number"`
	PrincipalPart    float64   `json:"principal_part"`
	InterestPart     float64   `json:"interest_part"`
	PenaltyAmount    float64   `json:"penalty_amount"`
	TotalPayment     float64   `json:"total_payment"`
	RemainingDebt    float64   `json:"remaining_debt"`
	DaysOverdue      int       `json:"days_overdue"`
	TotalWithPenalty float64   `json:"total_with_penalty"`
}
type CreditSummary struct {
	CreditID       int64   `json:"credit_id"`
	PaymentsCount  int     `json:"payments_count"`
	TotalPrincipal float64 `json:"total_principal"`
	TotalInterest  float64 `json:"total_interest"`
	TotalPenalty   float64 `json:"total_penalty"`
	TotalPaid      float64 `json:"total_paid"`
	AvgPayment     float64 `json:"avg_payment"`
	TotalOverdue   int     `json:"total_overdue_days"`
}

func CreditReportHandler(c *gin.Context) {
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

	var reportRows []CreditPaymentRecord
	if err := database.DB.
		Raw("SELECT * FROM credit_report(?, ?) ORDER BY credit_id, payment_date", startDate, endDate).
		Scan(&reportRows).Error; err != nil {
		c.String(http.StatusInternalServerError, "Ошибка при получении выплат")
		return
	}

	var summaryRows []CreditSummary
	if err := database.DB.
		Raw("SELECT * FROM credit_summary(?, ?)", startDate, endDate).
		Scan(&summaryRows).Error; err != nil {
		c.String(http.StatusInternalServerError, "Ошибка при получении сводки")
		return
	}

	var html strings.Builder
	html.WriteString(fmt.Sprintf(`<div class="report-heading">
		<h3>Отчёт по выплатам по кредитам за %s – %s</h3>
	</div>`,
		startDate.Format("02.01.2006"), endDate.Format("02.01.2006")))

	// Таблица выплат
	html.WriteString(`<div class="table">
		<div class="table-header">
			<div class="header__item">Кредит</div>
			<div class="header__item">Месяц</div>
			<div class="header__item">Дата</div>
			<div class="header__item">Основной долг</div>
			<div class="header__item">Проценты</div>
			<div class="header__item">Пени</div>
			<div class="header__item">Общая выплата</div>
			<div class="header__item">Остаток</div>
			<div class="header__item">Просрочка (дн.)</div>
		</div><div class="table-body">`)

	if len(reportRows) == 0 {
		html.WriteString(`<div class="table-row"><div class="table-data" colspan="9">Данные не найдены.</div></div>`)
	} else {
		for _, row := range reportRows {
			html.WriteString(fmt.Sprintf(`
			<div class="table-row">
				<div class="table-data">%d</div>
				<div class="table-data">%d</div>
				<div class="table-data">%s</div>
				<div class="table-data">%.2f</div>
				<div class="table-data">%.2f</div>
				<div class="table-data">%.2f</div>
				<div class="table-data">%.2f</div>
				<div class="table-data">%.2f</div>
				<div class="table-data">%d</div>
			</div>`,
				row.CreditID,
				row.MonthNumber,
				row.PaymentDate.Format("02.01.2006"),
				row.PrincipalPart,
				row.InterestPart,
				row.PenaltyAmount,
				row.TotalPayment,
				row.RemainingDebt,
				row.DaysOverdue))
		}
	}
	html.WriteString(`</div></div>`)

	// Сводка
	html.WriteString(`<div class="report-heading">
		<h3>Сводка по кредитам</h3>
	</div>
	<div class="table">
		<div class="table-header">
			<div class="header__item">Кредит</div>
			<div class="header__item">Платежей</div>
			<div class="header__item">Всего основной долг</div>
			<div class="header__item">Всего проценты</div>
			<div class="header__item">Всего пени</div>
			<div class="header__item">Всего выплачено</div>
			<div class="header__item">Средняя выплата</div>
			<div class="header__item">Всего дней просрочки</div>
		</div><div class="table-body">`)

	if len(summaryRows) == 0 {
		html.WriteString(`<div class="table-row"><div class="table-data" colspan="8">Сводка не найдена.</div></div>`)
	} else {
		for _, row := range summaryRows {
			html.WriteString(fmt.Sprintf(`
			<div class="table-row">
				<div class="table-data">%d</div>
				<div class="table-data">%d</div>
				<div class="table-data">%.2f</div>
				<div class="table-data">%.2f</div>
				<div class="table-data">%.2f</div>
				<div class="table-data">%.2f</div>
				<div class="table-data">%.2f</div>
				<div class="table-data">%d</div>
			</div>`,
				row.CreditID,
				row.PaymentsCount,
				row.TotalPrincipal,
				row.TotalInterest,
				row.TotalPenalty,
				row.TotalPaid,
				row.AvgPayment,
				row.TotalOverdue))
		}
	}
	html.WriteString(`</div></div>`)

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html.String()))
}
