package controllers

import (
	"github.com/gin-gonic/gin"
	"log"
	"production/database"
	"production/models"
)

func ListRawMaterials(c *gin.Context) {
	// Получаем все сырьевые материалы
	var rawmaterials []models.RawMaterial
	if err := database.DB.Find(&rawmaterials).Error; err != nil {
		log.Printf("Ошибка при получении сырья: %v", err)
		c.JSON(500, gin.H{"error": "Не удалось получить сырье"})
		return
	}
	// Получаем список единиц измерения
	var units []models.Unit
	if err := database.DB.Find(&units).Error; err != nil {
		log.Printf("Ошибка при получении единиц измерения: %v", err)
		c.JSON(500, gin.H{"error": "Не удалось получить единицы измерения"})
		return
	}

	// Создаем мапу unit_id -> unit_name
	unitMap := make(map[uint]string)
	for _, unit := range units {
		unitMap[unit.ID] = unit.Name
	}

	// Обновляем данные сырья, подставляя название единицы измерения
	type RawMaterialWithUnit struct {
		ID          uint
		Name        string
		Unit        string
		Quantity    float64
		TotalAmount float64
	}

	var rawmaterialsWithUnits []RawMaterialWithUnit
	for _, raw := range rawmaterials {
		rawmaterialsWithUnits = append(rawmaterialsWithUnits, RawMaterialWithUnit{
			ID:          raw.ID,
			Name:        raw.Name,
			Unit:        unitMap[raw.UnitID], // Заменяем ID на имя
			Quantity:    raw.Quantity,
			TotalAmount: raw.TotalAmount,
		})
	}

	// Передаем данные в шаблон
	c.HTML(200, "raw-materials.html", gin.H{
		"rawmaterials": rawmaterialsWithUnits,
	})
}

func DeleteRawMaterial(c *gin.Context) {
	id := c.Param("id")
	// Логируем ID
	log.Printf("Deleting raw material with ID: %s", id)

	// Удаляем запись
	if err := database.DB.Delete(&models.RawMaterial{}, id).Error; err != nil {
		log.Printf("Ошибка при удалении записи с ID %s: %v", id, err)
		c.JSON(500, gin.H{"error": "Не удалось удалить запись"})
		return
	}
	// Логируем успешное удаление
	log.Printf("Raw material with ID %s deleted successfully", id)
	c.JSON(200, gin.H{"success": true})
}
func AddRawMaterial(c *gin.Context) {
	var rawmaterial models.RawMaterial
	if err := c.ShouldBindJSON(&rawmaterial); err != nil {
		c.JSON(400, gin.H{"error": "Некорректные данные"})
		return
	}
	var existingMaterial models.RawMaterial
	if err := database.DB.Where("name = ?", rawmaterial.Name).First(&existingMaterial).Error; err == nil {
		c.JSON(400, gin.H{"error": "Сырье с таким названием уже существует"})
		return
	}

	// Создаем новую запись в базе данных
	if err := database.DB.Create(&rawmaterial).Error; err != nil {
		c.JSON(500, gin.H{"error": "Не удалось добавить запись"})
		return
	}

	c.JSON(200, gin.H{"success": true})
}

func GetRawMaterial(c *gin.Context) {
	id := c.Param("id")
	var rawmaterial models.RawMaterial

	// Ищем запись по ID
	if err := database.DB.First(&rawmaterial, id).Error; err != nil {
		log.Printf("Ошибка при получении сырья с ID %s: %v", id, err)
		c.JSON(500, gin.H{"error": "Не удалось получить запись"})
		return
	}

	// Отправляем данные в ответ
	c.JSON(200, gin.H{
		"success":     true,
		"rawmaterial": rawmaterial,
	})
}

func UpdateRawMaterial(c *gin.Context) {
	id := c.Param("id")
	var rawmaterial models.RawMaterial

	// Получаем данные из тела запроса
	if err := c.ShouldBindJSON(&rawmaterial); err != nil {
		c.JSON(400, gin.H{"error": "Некорректные данные"})
		return
	}

	// Ищем запись по ID
	var existingRawMaterial models.RawMaterial
	if err := database.DB.First(&existingRawMaterial, id).Error; err != nil {
		log.Printf("Ошибка при получении сырья с ID %s: %v", id, err)
		c.JSON(404, gin.H{"error": "Не удалось найти запись для обновления"})
		return
	}

	// Логируем полученные данные
	log.Printf("Updating raw material ID=%s: Name=%s, UnitID=%d, Quantity=%f, TotalAmount=%f",
		id, rawmaterial.Name, rawmaterial.UnitID, rawmaterial.Quantity, rawmaterial.TotalAmount)

	// Создаем мапу с обновляемыми значениями
	updateData := map[string]interface{}{
		"name":         rawmaterial.Name,
		"unit_id":      rawmaterial.UnitID,
		"quantity":     rawmaterial.Quantity,
		"total_amount": rawmaterial.TotalAmount,
	}

	// Выполняем обновление через мапу
	if err := database.DB.Model(&existingRawMaterial).Updates(updateData).Error; err != nil {
		log.Printf("Ошибка при обновлении сырья с ID %s: %v", id, err)
		c.JSON(500, gin.H{"error": "Не удалось обновить запись"})
		return
	}

	// Возвращаем успешный ответ
	c.JSON(200, gin.H{"success": true})
}

func GetRawMaterialsList(c *gin.Context) {
	var rawmaterials []models.RawMaterial
	if err := database.DB.Find(&rawmaterials).Error; err != nil {
		log.Printf("Ошибка при получении сырья: %v", err)
		c.JSON(500, gin.H{"success": false, "error": "Не удалось получить сырье"})
		return
	}

	c.JSON(200, gin.H{"success": true, "raw_materials": rawmaterials})
}
