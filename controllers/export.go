package controllers

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/fumiama/go-docx"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jung-kurt/gofpdf"
	"github.com/xuri/excelize/v2"
	"net/http"
	"os"
	"production/database"
	"production/models"
	"strings"
	"time"
)

// ExportHandler обрабатывает запросы на экспорт отчетов
func ExportReport(c *gin.Context) {
	type ExportRequest struct {
		Content string `json:"content"`
		Title   string `json:"title"`
		Period  string `json:"period"`
		Format  string `json:"format"`
	}

	var req ExportRequest
	if err := c.BindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, "Неверный формат запроса")
		return
	}

	tokenString, err := c.Cookie("token")
	if err != nil || tokenString == "" {
		c.String(http.StatusUnauthorized, "Пожалуйста, авторизуйтесь")
		return
	}

	// Парсим токен с нашими Claims
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})
	if err != nil {
		c.String(http.StatusUnauthorized, "Неверный токен: "+err.Error())
		return
	}

	// Проверяем валидность токена
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		c.String(http.StatusUnauthorized, "Неверный токен")
		return
	}

	// Получаем данные сотрудника из БД
	var employee models.Employee
	if err := database.DB.Preload("Position").First(&employee, claims.EmployeeID).Error; err != nil {
		c.String(http.StatusInternalServerError, "Не удалось получить данные сотрудника")
		return
	}

	// Получаем данные директора
	var director models.Employee
	if err := database.DB.Joins("JOIN positions ON positions.id = employees.position_id").
		Where("positions.name = ?", "Директор").
		First(&director).Error; err != nil {
		c.String(http.StatusInternalServerError, "Не удалось получить данные директора")
		return
	}

	// Формируем метаданные отчета
	metadata := struct {
		Organization  string
		ReportTitle   string
		Period        string
		CreationDate  string
		Responsible   string
		Director      string
		ReportContent string
	}{
		Organization:  "Спортивный инвентарь",
		ReportTitle:   req.Title,
		Period:        req.Period,
		CreationDate:  time.Now().Format("02.01.2006"),
		Responsible:   employee.FullName,
		Director:      director.FullName,
		ReportContent: req.Content,
	}

	// Генерируем файл в нужном формате
	var fileData []byte
	var contentType string
	var fileExt string

	switch req.Format {
	case "docx":
		fileData, err = generateWordReport(metadata)
		contentType = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
		fileExt = "docx"
	case "xlsx":
		fileData, err = generateExcelReport(metadata)
		contentType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
		fileExt = "xlsx"
	case "pdf":
		fileData, err = generatePdfReport(metadata)
		contentType = "application/pdf"
		fileExt = "pdf"
	default:
		c.String(http.StatusBadRequest, "Неверный формат экспорта")
		return
	}

	if err != nil {
		c.String(http.StatusInternalServerError, "Ошибка генерации отчета: "+err.Error())
		return
	}

	// Формируем имя файла
	fileName := fmt.Sprintf("%s_%s.%s",
		strings.ReplaceAll(req.Title, " ", "_"),
		time.Now().Format("2006-01-02"),
		fileExt)

	// Отправляем файл клиенту
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	c.Data(http.StatusOK, contentType, fileData)
}

// ReportMetadata содержит метаданные для генерации отчета
type ReportMetadata struct {
	Organization  string
	ReportTitle   string
	Period        string
	CreationDate  string
	Responsible   string
	Director      string
	ReportContent string
}

// generateExcelReport создает строгий Excel-файл из HTML-таблицы
func generateExcelReport(meta ReportMetadata) ([]byte, error) {
	f := excelize.NewFile()
	sheetName := "Sheet1"

	// Создаем стили
	titleStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 16, Family: "Times New Roman"},
		Alignment: &excelize.Alignment{Horizontal: "left"},
	})

	subtitleStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 14, Family: "Times New Roman"},
		Alignment: &excelize.Alignment{Horizontal: "left"},
	})

	infoStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Size: 11, Family: "Times New Roman"},
	})

	headingStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 13, Family: "Times New Roman"},
		Alignment: &excelize.Alignment{Horizontal: "left"},
	})

	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 10, Family: "Times New Roman"},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})

	cellStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Size: 10, Family: "Times New Roman"},
		Alignment: &excelize.Alignment{Horizontal: "left", Vertical: "center", WrapText: true},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})

	signatureStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Italic: true, Size: 11, Family: "Times New Roman"},
		Alignment: &excelize.Alignment{Horizontal: "right"},
	})

	// Устанавливаем заголовок организации
	f.SetCellValue(sheetName, "A1", meta.Organization)
	f.MergeCell(sheetName, "A1", "Z1")
	f.SetCellStyle(sheetName, "A1", "A1", titleStyle)

	// Название отчета
	f.SetCellValue(sheetName, "A2", meta.ReportTitle)
	f.MergeCell(sheetName, "A2", "Z2")
	f.SetCellStyle(sheetName, "A2", "A2", subtitleStyle)

	// Метаданные
	currentRow := 3
	f.SetCellValue(sheetName, fmt.Sprintf("A%d", currentRow+1), fmt.Sprintf("Дата составления: %s", meta.CreationDate))
	f.SetCellValue(sheetName, fmt.Sprintf("A%d", currentRow+2), fmt.Sprintf("Ответственный: %s", meta.Responsible))
	f.SetCellStyle(sheetName, fmt.Sprintf("A%d", currentRow+1), fmt.Sprintf("A%d", currentRow+2), infoStyle)
	currentRow += 3

	// Разделитель (пустая строка)
	currentRow++

	// Парсинг HTML содержимого
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(meta.ReportContent))
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга HTML: %v", err)
	}

	// Обработка каждого раздела отчета
	doc.Find(".report-heading").Each(func(i int, heading *goquery.Selection) {
		headingText := strings.TrimSpace(heading.Find("h3").Text())

		// Заголовок раздела
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", currentRow), headingText)
		f.SetCellStyle(sheetName, fmt.Sprintf("A%d", currentRow), fmt.Sprintf("A%d", currentRow), headingStyle)
		currentRow++

		nextTable := heading.NextFiltered(".table")
		if nextTable.Length() == 0 {
			return
		}

		// Извлечение заголовков таблицы
		var headers []string
		nextTable.Find(".table-header .header__item").Each(func(i int, header *goquery.Selection) {
			headers = append(headers, strings.TrimSpace(header.Text()))
		})

		// Извлечение данных таблицы
		var rows [][]string
		nextTable.Find(".table-body .table-row").Each(func(i int, row *goquery.Selection) {
			var rowData []string
			row.Find(".table-data").Each(func(j int, cell *goquery.Selection) {
				rowData = append(rowData, strings.TrimSpace(cell.Text()))
			})
			rows = append(rows, rowData)
		})

		colCount := len(headers)
		if colCount == 0 && len(rows) > 0 {
			colCount = len(rows[0])
		}
		if colCount == 0 {
			return
		}

		// Записываем заголовки таблицы
		startCol := 'A'
		for i, header := range headers {
			if i >= 26 { // Ограничение на количество столбцов (A-Z)
				break
			}
			col := string(startCol + rune(i))
			f.SetCellValue(sheetName, fmt.Sprintf("%s%d", col, currentRow), header)
			f.SetCellStyle(sheetName, fmt.Sprintf("%s%d", col, currentRow), fmt.Sprintf("%s%d", col, currentRow), headerStyle)
		}

		// Записываем данные таблицы
		for _, row := range rows {
			currentRow++
			for i, cell := range row {
				if i >= 26 { // Ограничение на количество столбцов (A-Z)
					break
				}
				col := string(startCol + rune(i))
				f.SetCellValue(sheetName, fmt.Sprintf("%s%d", col, currentRow), cell)
				f.SetCellStyle(sheetName, fmt.Sprintf("%s%d", col, currentRow), fmt.Sprintf("%s%d", col, currentRow), cellStyle)
			}
		}

		// Автонастройка ширины столбцов
		for i := 0; i < colCount && i < 26; i++ {
			col := string(startCol + rune(i))
			maxWidth := 10.0 // минимальная ширина

			// Проверяем ширину заголовка
			if i < len(headers) {
				headerWidth := float64(len(headers[i])) * 1.2
				if headerWidth > maxWidth {
					maxWidth = headerWidth
				}
			}

			// Проверяем ширину данных
			for _, row := range rows {
				if i < len(row) {
					cellWidth := float64(len(row[i])) * 1.1
					if cellWidth > maxWidth {
						maxWidth = cellWidth
					}
				}
			}

			// Ограничиваем максимальную ширину
			if maxWidth > 50 {
				maxWidth = 50
			}

			f.SetColWidth(sheetName, col, col, maxWidth)
		}

		currentRow += 2 // Отступ после таблицы
	})

	// Подпись директора
	f.SetCellValue(sheetName, fmt.Sprintf("A%d", currentRow), fmt.Sprintf("Директор: ___________________ %s", meta.Director))
	f.SetCellStyle(sheetName, fmt.Sprintf("A%d", currentRow), fmt.Sprintf("A%d", currentRow), signatureStyle)

	// Сохраняем файл в буфер
	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, fmt.Errorf("ошибка сохранения Excel: %v", err)
	}

	return buf.Bytes(), nil
}

// generatePdfReport создает отчет в формате PDF
func generatePdfReport(meta ReportMetadata) ([]byte, error) {
	// Инициализация PDF документа
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(10, 10, 10)
	pdf.SetAutoPageBreak(true, 10)

	// Установка футера с номером страницы
	pdf.SetFooterFunc(func() {
		pdf.SetY(-10)
		pdf.SetFont("TildaSans", "I", 8)
		pdf.CellFormat(0, 10, fmt.Sprintf("Страница %d", pdf.PageNo()), "", 0, "C", false, 0, "")
	})

	pdf.AddPage()

	leftMargin, _, _, _ := pdf.GetMargins()

	// Добавление шрифтов
	fontPath := "public/fonts/TildaSans-VF.ttf"
	pdf.AddUTF8Font("TildaSans", "", fontPath)
	pdf.AddUTF8Font("TildaSans", "B", fontPath)
	pdf.AddUTF8Font("TildaSans", "I", fontPath)

	// Добавление логотипа
	logoPath := "public/images/logos/logo.png"
	if _, err := os.Stat(logoPath); err == nil {
		logoWidth := 10.0
		margin := 10.0
		pageWidth, _ := pdf.GetPageSize()
		x := pageWidth - logoWidth - margin
		y := 10.0 // верхний отступ

		pdf.Image(logoPath, x, y, logoWidth, 0, false, "", 0, "")
	}

	// Заголовок отчета
	pdf.SetFont("TildaSans", "B", 16)
	pdf.CellFormat(0, 12, meta.Organization, "", 1, "C", false, 0, "")
	pdf.SetFont("TildaSans", "B", 14)
	pdf.CellFormat(0, 10, meta.ReportTitle, "", 1, "C", false, 0, "")
	pdf.Ln(4)

	// Дата и ответственный
	pdf.SetFont("TildaSans", "", 11)
	pdf.CellFormat(0, 7, fmt.Sprintf("Дата составления: %s", meta.CreationDate), "", 1, "L", false, 0, "")
	pdf.CellFormat(0, 7, fmt.Sprintf("Ответственный: %s", meta.Responsible), "", 1, "L", false, 0, "")
	pdf.Ln(8)

	// Разделитель
	pdf.Line(10, pdf.GetY(), 200, pdf.GetY())
	pdf.Ln(6)

	// Парсинг HTML содержимого
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(meta.ReportContent))
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга HTML: %v", err)
	}

	// Обработка каждого раздела отчета
	doc.Find(".report-heading").Each(func(i int, heading *goquery.Selection) {
		headingText := strings.TrimSpace(heading.Find("h3").Text())
		pdf.SetFont("TildaSans", "B", 13)
		pdf.MultiCell(0, 8, headingText, "", "L", false)
		pdf.Ln(4)

		nextTable := heading.NextFiltered(".table")
		if nextTable.Length() == 0 {
			return
		}

		// Извлечение заголовков таблицы
		var headers []string
		nextTable.Find(".table-header .header__item").Each(func(i int, header *goquery.Selection) {
			headers = append(headers, strings.TrimSpace(header.Text()))
		})

		// Извлечение данных таблицы
		var rows [][]string
		nextTable.Find(".table-body .table-row").Each(func(i int, row *goquery.Selection) {
			var rowData []string
			row.Find(".table-data").Each(func(j int, cell *goquery.Selection) {
				rowData = append(rowData, strings.TrimSpace(cell.Text()))
			})
			rows = append(rows, rowData)
		})

		colCount := len(headers)
		if colCount == 0 && len(rows) > 0 {
			colCount = len(rows[0])
		}
		if colCount == 0 {
			return
		}

		// Расчет оптимальной ширины столбцов
		colWidths := make([]float64, colCount)
		pageWidth, _ := pdf.GetPageSize()
		maxTableWidth := pageWidth - leftMargin*2
		minColWidth := 15.0 // Минимальная ширина столбца
		maxColWidth := 60.0 // Максимальная ширина столбца

		// Функция для расчета ширины текста
		calcTextWidth := func(text string, fontStyle string, fontSize float64) float64 {
			pdf.SetFont("TildaSans", fontStyle, fontSize)
			return pdf.GetStringWidth(text) + 6 // Добавляем небольшой отступ
		}

		// Расчет ширины для заголовков
		pdf.SetFont("TildaSans", "B", 10)
		for i := 0; i < colCount; i++ {
			headerText := ""
			if i < len(headers) {
				headerText = headers[i]
			}
			width := calcTextWidth(headerText, "B", 10)

			if width > maxColWidth {
				smallWidth := calcTextWidth(headerText, "B", 8)
				if smallWidth < maxColWidth {
					width = smallWidth
				} else {
					width = maxColWidth
				}
			}

			if width < minColWidth {
				width = minColWidth
			}
			colWidths[i] = width
		}

		// Расчет ширины для содержимого
		pdf.SetFont("TildaSans", "", 10)
		for _, row := range rows {
			for i := 0; i < colCount; i++ {
				cellText := ""
				if i < len(row) {
					cellText = row[i]
				}
				width := calcTextWidth(cellText, "", 10)
				if width > colWidths[i] && width <= maxColWidth {
					colWidths[i] = width
				}
			}
		}

		// Проверка общей ширины таблицы
		totalWidth := 0.0
		for _, width := range colWidths {
			totalWidth += width
		}

		// Если таблица не помещается - масштабируем
		if totalWidth > maxTableWidth {
			scale := maxTableWidth / totalWidth
			for i := range colWidths {
				colWidths[i] *= scale
				if colWidths[i] < minColWidth {
					colWidths[i] = minColWidth
				}
			}
		}

		lineHeight := 6.0
		currentY := pdf.GetY()

		// Рисуем заголовки таблицы (без заливки)
		pdf.SetFont("TildaSans", "B", 10)
		headerHeights := make([]float64, colCount)

		for i := 0; i < colCount; i++ {
			headerText := ""
			if i < len(headers) {
				headerText = headers[i]
			}

			x := leftMargin
			for j := 0; j < i; j++ {
				x += colWidths[j]
			}

			if calcTextWidth(headerText, "B", 10) > colWidths[i]-4 {
				pdf.SetFont("TildaSans", "B", 8)
				lines := pdf.SplitText(headerText, colWidths[i]-4)
				headerHeights[i] = lineHeight * float64(len(lines))
				pdf.SetFont("TildaSans", "B", 10)
			} else {
				headerHeights[i] = lineHeight
			}
		}

		// Находим максимальную высоту заголовка
		maxHeaderHeight := lineHeight
		for _, h := range headerHeights {
			if h > maxHeaderHeight {
				maxHeaderHeight = h
			}
		}

		// Рисуем ячейки заголовков (только границы)
		for i := 0; i < colCount; i++ {
			headerText := ""
			if i < len(headers) {
				headerText = headers[i]
			}

			x := leftMargin
			for j := 0; j < i; j++ {
				x += colWidths[j]
			}

			// Рисуем только границу (без заливки)
			pdf.SetDrawColor(0, 0, 0) // Черные границы
			pdf.Rect(x, currentY, colWidths[i], maxHeaderHeight, "D")

			// Устанавливаем текст
			if calcTextWidth(headerText, "B", 10) > colWidths[i]-4 {
				pdf.SetFont("TildaSans", "B", 8)
				pdf.SetXY(x, currentY)
				pdf.MultiCell(colWidths[i], lineHeight, headerText, "", "C", false)
				pdf.SetFont("TildaSans", "B", 10)
			} else {
				pdf.SetXY(x, currentY)
				pdf.MultiCell(colWidths[i], lineHeight, headerText, "", "C", false)
			}
		}

		// Рисуем содержимое таблицы
		pdf.SetFont("TildaSans", "", 10)
		pdf.SetXY(leftMargin, currentY+maxHeaderHeight)

		for _, row := range rows {
			// Рассчитываем высоту строки
			maxLines := 1
			for i := 0; i < colCount; i++ {
				cellText := ""
				if i < len(row) {
					cellText = row[i]
				}
				lines := pdf.SplitText(cellText, colWidths[i]-2)
				if len(lines) > maxLines {
					maxLines = len(lines)
				}
			}
			rowHeight := lineHeight * float64(maxLines)
			currentY := pdf.GetY()

			// Проверяем, нужно ли переносить на новую страницу
			if currentY+rowHeight > 270 {
				pdf.AddPage()
				currentY = pdf.GetY()

				// Повторяем заголовки на новой странице
				for i := 0; i < colCount; i++ {
					headerText := ""
					if i < len(headers) {
						headerText = headers[i]
					}

					x := leftMargin
					for j := 0; j < i; j++ {
						x += colWidths[j]
					}

					// Рисуем границу заголовка
					pdf.Rect(x, currentY, colWidths[i], maxHeaderHeight, "D")

					// Устанавливаем текст заголовка
					if calcTextWidth(headerText, "B", 10) > colWidths[i]-4 {
						pdf.SetFont("TildaSans", "B", 8)
						pdf.SetXY(x, currentY)
						pdf.MultiCell(colWidths[i], lineHeight, headerText, "", "C", false)
						pdf.SetFont("TildaSans", "B", 10)
					} else {
						pdf.SetXY(x, currentY)
						pdf.MultiCell(colWidths[i], lineHeight, headerText, "", "C", false)
					}
				}

				pdf.SetFont("TildaSans", "", 10)
				currentY = currentY + maxHeaderHeight
			}

			// Рисуем строку с данными
			for i := 0; i < colCount; i++ {
				cellText := ""
				if i < len(row) {
					cellText = row[i]
				}

				x := leftMargin
				for j := 0; j < i; j++ {
					x += colWidths[j]
				}

				// Рисуем границу ячейки
				pdf.Rect(x, currentY, colWidths[i], rowHeight, "D")

				// Добавляем текст
				pdf.SetXY(x, currentY)
				pdf.MultiCell(colWidths[i], lineHeight, cellText, "", "L", false)
			}

			pdf.SetXY(leftMargin, currentY+rowHeight)
		}

		pdf.Ln(8)
	})

	// Строка для подписи
	pdf.SetFont("TildaSans", "I", 11)
	pdf.SetY(270)
	pdf.CellFormat(0, 8, fmt.Sprintf("Директор: ___________________ %s", meta.Director), "", 1, "R", false, 0, "")

	// Генерация PDF
	buf := new(bytes.Buffer)
	if err := pdf.Output(buf); err != nil {
		return nil, fmt.Errorf("ошибка генерации PDF: %v", err)
	}
	return buf.Bytes(), nil
}

// generateWordReport создает отчет в формате Word

func generateWordReport(meta ReportMetadata) ([]byte, error) {
	doc := docx.New().WithDefaultTheme()

	// Заголовок отчёта
	doc.AddParagraph().
		Justification("center").
		AddText(meta.Organization).Size("40").Bold()
	doc.AddParagraph().
		Justification("center").
		AddText(meta.ReportTitle).Size("30").Bold()
	doc.AddParagraph().
		AddText(fmt.Sprintf("Дата составления: %s", meta.CreationDate))
	doc.AddParagraph().
		AddText(fmt.Sprintf("Ответственный: %s", meta.Responsible))
	doc.AddParagraph()

	// Парсим HTML содержимого
	htmlReader := strings.NewReader(meta.ReportContent)
	docHTML, err := goquery.NewDocumentFromReader(htmlReader)
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга HTML: %v", err)
	}

	// Обрабатываем каждый раздел
	docHTML.Find(".report-heading").Each(func(i int, heading *goquery.Selection) {
		// Заголовок секции
		sectionTitle := strings.TrimSpace(heading.Find("h3").Text())
		doc.AddParagraph().
			AddText(sectionTitle).Bold()

		// Ищем следующий элемент с классом .table
		tblSel := heading.NextFiltered(".table")
		if tblSel.Length() == 0 {
			return
		}

		// Собираем заголовки столбцов
		var headers []string
		tblSel.Find(".table-header .header__item").Each(func(_ int, h *goquery.Selection) {
			headers = append(headers, strings.TrimSpace(h.Text()))
		})

		// Собираем данные строк
		var rows [][]string
		tblSel.Find(".table-body .table-row").Each(func(_ int, rowSel *goquery.Selection) {
			var rowData []string
			rowSel.Find(".table-data").Each(func(_ int, cellSel *goquery.Selection) {
				rowData = append(rowData, strings.TrimSpace(cellSel.Text()))
			})
			rows = append(rows, rowData)
		})

		// Инициализируем таблицу: 1 строка заголовка + len(rows) строк данных
		tbl := doc.AddTable(len(rows)+1, len(headers), 0, nil) // :contentReference[oaicite:0]{index=0}

		// Заполняем строку заголовков
		for ci, h := range headers {
			cell := tbl.TableRows[0].TableCells[ci] // :contentReference[oaicite:1]{index=1}
			para := cell.AddParagraph()             // :contentReference[oaicite:2]{index=2}
			para.AddText(h).Bold()
		}

		// Заполняем строки с данными
		for ri, rowData := range rows {
			for ci, txt := range rowData {
				cell := tbl.TableRows[ri+1].TableCells[ci]
				para := cell.AddParagraph()
				para.AddText(txt)
			}
		}

		doc.AddParagraph() // отступ после таблицы
	})

	// Подпись
	doc.AddParagraph().
		Justification("end").
		AddText(fmt.Sprintf("Директор: ___________________ %s", meta.Director))

	// Пишем в буфер
	var buf bytes.Buffer
	if _, err := doc.WriteTo(&buf); err != nil {
		return nil, fmt.Errorf("ошибка генерации DOCX: %v", err)
	}
	return buf.Bytes(), nil
}
