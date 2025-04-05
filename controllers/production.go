package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"production/database"
	"production/models"
	"strconv"
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

// Запрос на производство
type ProduceProductRequest struct {
	Quantity       float64 `json:"quantity" binding:"required"`
	ProductionDate string  `json:"production_date" binding:"required"`
	EmployeeID     uint    `json:"employee_id" binding:"required"`
}

func ProduceProduct(c *gin.Context) {
	// Получаем ID продукта из параметра запроса
	productID, err := strconv.Atoi(c.Param("product_id"))
	if err != nil {
		log.Println("Ошибка конвертации product_id:", err)
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Неверный ID продукта"})
		return
	}

	// Получаем продукт. Через First - первая запись с заданным id, не Find - потому, что вернет несколько записей (массив)
	var product models.FinishedGood
	if err := database.DB.First(&product, productID).Error; err != nil {
		log.Println("Ошибка получения продукта:", err)
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "Продукт не найден"})
		return
	}

	// Валидация запроса
	var request ProduceProductRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Println("Ошибка валидации данных:", err)
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Неверные данные"})
		return
	}

	// Получаем ингредиенты
	var ingredients []models.Ingredient
	if err := database.DB.Where("product_id = ?", productID).Find(&ingredients).Error; err != nil {
		log.Println("Ошибка получения ингредиентов:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Ошибка сервера"})
		return
	}

	// Создаем мапу для сырья. Через ингредиенты получаем сырье.
	// Создаем пустую map - словарь, где uint ключ id сырья, make - выделяет память
	// Перебор всех ингредиентов
	rawMaterialsMap := make(map[uint]models.RawMaterial)
	for _, ingredient := range ingredients {
		// Создаем структуру для сырья
		var rawMaterial models.RawMaterial
		// Перебираем все сырье ингредиентов
		if err := database.DB.First(&rawMaterial, ingredient.RawMaterialID).Error; err != nil {
			log.Println("Ошибка получения сырья:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Ошибка сервера"})
			return
		}
		// Добавляем сырье в мапу
		rawMaterialsMap[rawMaterial.ID] = rawMaterial
	}

	// Проверка наличия достаточного количества сырья
	// Переменная хранит недостающее сырье
	insufficientMaterials := []string{}
	// Общая стоимость производства
	totalProductionCost := 0.0
	// Проходим по каждому ингредиенту который требуется для производства
	// Расчитываем сырье
	// Необходимое сырье берем из ранее созданной мапы
	for _, ingredient := range ingredients {
		// получаем необходимое количество
		requiredAmount := ingredient.Quantity * request.Quantity
		rawMaterial := rawMaterialsMap[ingredient.RawMaterialID]
		// записываем недостающее сырье
		if rawMaterial.Quantity < requiredAmount {
			insufficientMaterials = append(insufficientMaterials, fmt.Sprintf("%s (не хватает %.2f)", rawMaterial.Name, requiredAmount-rawMaterial.Quantity))
		}
	}

	if len(insufficientMaterials) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Недостаточно сырья", "details": insufficientMaterials})
		return
	}

	// Начинаем транзакцию
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Обновляем сырье и рассчитываем стоимость производства
	for _, ingredient := range ingredients {
		// перебираем какое сырье необходимо списать
		requiredAmount := ingredient.Quantity * request.Quantity
		// Достаем сырье из мапы
		rawMaterial := rawMaterialsMap[ingredient.RawMaterialID]
		// вычисляем среднюю стоимость сырья
		averagePrice := rawMaterial.TotalAmount / rawMaterial.Quantity
		// Расчитываем стоимость списанного сырья
		usedCost := requiredAmount * averagePrice
		// Увеличиваем общую стоимость производства
		totalProductionCost += usedCost

		rawMaterial.Quantity -= requiredAmount
		rawMaterial.TotalAmount -= usedCost

		if err := tx.Save(&rawMaterial).Error; err != nil {
			tx.Rollback()
			log.Println("Ошибка обновления сырья:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Ошибка сервера"})
			return
		}
	}

	// Обновляем продукт
	product.Quantity += request.Quantity
	product.TotalAmount += totalProductionCost
	// Обновляем запись продукции (сумма, количество)
	if product.Quantity != 0 {
		productPrice := product.TotalAmount / product.Quantity
		product.TotalAmount = productPrice * product.Quantity
	}
	// сохраняем продукт
	if err := tx.Save(&product).Error; err != nil {
		tx.Rollback()
		log.Println("Ошибка обновления продукта:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Ошибка сервера"})
		return
	}

	// Записываем запись о производстве
	productionDate, err := time.Parse("2006-01-02", request.ProductionDate)
	if err != nil {
		tx.Rollback()
		log.Println("Ошибка парсинга даты:", err)
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Неверный формат даты"})
		return
	}

	production := models.ProductProduction{
		ProductID:      uint(productID),
		Quantity:       request.Quantity,
		ProductionDate: productionDate,
		EmployeeID:     request.EmployeeID,
	}

	if err := tx.Create(&production).Error; err != nil {
		tx.Rollback()
		log.Println("Ошибка создания записи производства:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Ошибка сервера"})
		return
	}

	// Фиксируем транзакцию
	if err := tx.Commit().Error; err != nil {
		log.Println("Ошибка фиксации транзакции:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Ошибка сервера"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Производство успешно завершено"})
}

func GetIngredientsList(c *gin.Context) {
	var ingredients []models.Ingredient
	if err := database.DB.Preload("Product").Preload("RawMaterial").Find(&ingredients).Error; err != nil {
		log.Printf("Ошибка при получении ингредиентов: %v", err)
		c.JSON(500, gin.H{"success": false, "error": "Не удалось получить ингредиенты"})
		return
	}

	c.JSON(200, gin.H{"success": true, "ingredients": ingredients})
}
