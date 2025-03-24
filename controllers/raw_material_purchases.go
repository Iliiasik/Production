package controllers

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"production/database"
	"production/models"
	"time"
)

func ListRawMaterialPurchases(c *gin.Context) {
	// Получаем все сырьевые материалы
	var purchases []models.RawMaterialPurchase
	if err := database.DB.Find(&purchases).Error; err != nil {
		log.Printf("Ошибка при получении закупок сырья: %v", err)
		c.JSON(500, gin.H{"error": "Не удалось получить закупленное сырье"})
		return
	}
	// Получаем список сырья
	var rawmaterials []models.RawMaterial
	if err := database.DB.Find(&rawmaterials).Error; err != nil {
		log.Printf("Ошибка при получении сырья: %v", err)
		c.JSON(500, gin.H{"error": "Не удалось получить сырье"})
		return
	}
	// Создаем мапу rawmaterial_id -> rawmaterial_name
	rawmaterialMap := make(map[uint]string)
	for _, rawmaterial := range rawmaterials {
		rawmaterialMap[rawmaterial.ID] = rawmaterial.Name
	}

	// Список сотрудников
	var employees []models.Employee
	if err := database.DB.Find(&employees).Error; err != nil {
		log.Printf("Ошибка при получении списка сотрудников: %v", err)
		c.JSON(500, gin.H{"error": "Не удалось получить список сотрудников"})
		return
	}
	// Создаем мапу employee_id -> employee_name
	employeeMap := make(map[uint]string)
	for _, employee := range employees {
		employeeMap[employee.ID] = employee.FullName
	}

	type PurchasesWithDetails struct {
		ID           uint
		Material     string
		Quantity     float64
		TotalAmount  float64
		PurchaseDate time.Time
		Employee     string
	}

	var purchasesWithDetails []PurchasesWithDetails
	for _, prc := range purchases {
		purchasesWithDetails = append(purchasesWithDetails, PurchasesWithDetails{
			ID:           prc.ID,
			Material:     rawmaterialMap[prc.RawMaterialID],
			Quantity:     prc.Quantity,
			TotalAmount:  prc.TotalAmount,
			PurchaseDate: prc.PurchaseDate,
			Employee:     employeeMap[prc.EmployeeID],
		})
	}

	// Передаем данные в шаблон
	c.HTML(200, "raw-material-purchases.html", gin.H{
		"purchases": purchasesWithDetails,
	})
}

func DeletePurchase(c *gin.Context) {
	id := c.Param("id")
	log.Printf("Deleting purchase with ID: %s", id)

	// Находим запись закупки перед удалением
	var purchase models.RawMaterialPurchase
	if err := database.DB.First(&purchase, id).Error; err != nil {
		log.Printf("Ошибка при получении закупки с ID %s: %v", id, err)
		c.JSON(404, gin.H{"error": "Закупка не найдена"})
		return
	}

	// Находим соответствующее сырье
	var rawMaterial models.RawMaterial
	if err := database.DB.First(&rawMaterial, purchase.RawMaterialID).Error; err != nil {
		log.Printf("Ошибка при получении сырья с ID %d: %v", purchase.RawMaterialID, err)
		c.JSON(500, gin.H{"error": "Ошибка сервера"})
		return
	}

	// Восстанавливаем количество и сумму сырья
	rawMaterial.Quantity -= purchase.Quantity
	rawMaterial.TotalAmount -= purchase.TotalAmount
	if rawMaterial.Quantity < 0 || rawMaterial.TotalAmount < 0 {
		log.Println("Ошибка: количество или сумма сырья ушли в отрицательное значение")
		c.JSON(500, gin.H{"error": "Ошибка при восстановлении данных"})
		return
	}
	database.DB.Save(&rawMaterial)

	// Восстанавливаем бюджет
	var budget models.Budget
	if err := database.DB.First(&budget).Error; err != nil {
		log.Println("Ошибка получения бюджета:", err)
		c.JSON(500, gin.H{"error": "Ошибка сервера"})
		return
	}
	budget.TotalAmount += purchase.TotalAmount
	database.DB.Save(&budget)

	// Удаляем запись о закупке
	if err := database.DB.Delete(&purchase).Error; err != nil {
		log.Printf("Ошибка при удалении закупки с ID %s: %v", id, err)
		c.JSON(500, gin.H{"error": "Не удалось удалить запись"})
		return
	}

	log.Printf("Purchase with ID %s deleted successfully", id)
	c.JSON(200, gin.H{"success": true})
}

func AddPurchase(c *gin.Context) {

	// создаем переменную purchase, которая будет хранить все данные о закупке (структура по модельке rawmaterialpurchase)
	var purchase models.RawMaterialPurchase

	// парсим json в структуру purchase
	// проверяем валидность данных (заполнение всех полей, типы данных)
	// при nil отправляем ошибку
	if err := c.ShouldBindJSON(&purchase); err != nil {
		log.Println("Ошибка валидации данных:", err)
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Неверные данные"})
		return
	}

	// Проверяем бюджет
	// бюджет - первая запись в таблице (т.к я не знаю как будет дальше, может у нас будет несколько записей бюджета, тогда поменять first на find и id конкретного бюджета)
	// SELECT * FROM budgets ORDER BY id LIMIT 1; - данный запрос осуществляется при DB.First(*переменная*)
	var budget models.Budget
	if err := database.DB.First(&budget).Error; err != nil {
		log.Println("Ошибка получения бюджета:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Ошибка сервера"})
		return
	}
	// проверка средств (проверка происходит на клиентской стороне, но в случае ошибки, или атаки в бд, на сервере тоже проверяется)
	if budget.TotalAmount < purchase.TotalAmount {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Недостаточно средств в бюджете"})
		return
	}

	// Обновляем бюджет
	budget.TotalAmount -= purchase.TotalAmount
	database.DB.Save(&budget)

	// Обновляем количество сырья
	// находим запись сырья с rawmaterialid

	// Почему first, если ищем по ID
	// SELECT * FROM raw_materials WHERE id = ? LIMIT 1;
	// этот запрос эквивалентен gorm запросу DB.First
	// Если использовать Find, то результат будет пустой rawmaterial (при условии, что такой записи нет)
	// Если использовать First, то вернется ошибка при том же условии.
	// Если сырья с таким ID нет — значит ошибка (нельзя делать закупку для несуществующего сырья)
	var rawMaterial models.RawMaterial
	if err := database.DB.First(&rawMaterial, purchase.RawMaterialID).Error; err != nil {
		log.Println("Ошибка получения сырья:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Ошибка сервера"})
		return
	}
	// добавляем к сумме и количеству сырья
	rawMaterial.Quantity += purchase.Quantity
	rawMaterial.TotalAmount += purchase.TotalAmount
	database.DB.Save(&rawMaterial)

	// Устанавливаем дату покупки (текущая) и сохраняем
	database.DB.Create(&purchase)

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Закупка успешно добавлена"})
}
