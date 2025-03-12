package controllers

import (
	"github.com/gin-gonic/gin"
	"log"
	"production/database"
	"production/models"
)

func ListUnits(c *gin.Context) {
	// Fetch все записи
	var units []models.Unit
	if err := database.DB.Find(&units).Error; err != nil {
		log.Printf("Ошибка при получении единиц измерения: %v", err)
		c.JSON(500, gin.H{"error": "Не удалось получить единицы измерения"})
		return
	}

	// Передача данных в шаблон
	c.HTML(200, "units.html", gin.H{
		"units": units,
	})
}

func DeleteUnit(c *gin.Context) {
	id := c.Param("id")
	// Логируем ID
	log.Printf("Deleting unit with ID: %s", id)

	// Удаляем запись
	if err := database.DB.Delete(&models.Unit{}, id).Error; err != nil {
		log.Printf("Ошибка при удалении записи с ID %s: %v", id, err)
		c.JSON(500, gin.H{"error": "Не удалось удалить запись"})
		return
	}
	// Логируем успешное удаление
	log.Printf("Unit with ID %s deleted successfully", id)
	c.JSON(200, gin.H{"success": true})
}
func AddUnit(c *gin.Context) {
	var unit models.Unit
	if err := c.ShouldBindJSON(&unit); err != nil {
		c.JSON(400, gin.H{"error": "Некорректные данные"})
		return
	}

	// Создаем новую запись в базе данных
	if err := database.DB.Create(&unit).Error; err != nil {
		c.JSON(500, gin.H{"error": "Не удалось добавить запись"})
		return
	}

	c.JSON(200, gin.H{"success": true})
}

func GetUnit(c *gin.Context) {
	id := c.Param("id")
	var unit models.Unit

	// Ищем запись по ID
	if err := database.DB.First(&unit, id).Error; err != nil {
		log.Printf("Ошибка при получении единицы измерения с ID %s: %v", id, err)
		c.JSON(500, gin.H{"error": "Не удалось получить запись"})
		return
	}

	// Отправляем данные в ответ
	c.JSON(200, gin.H{
		"success": true,
		"unit":    unit,
	})
}
func UpdateUnit(c *gin.Context) {
	id := c.Param("id")
	var unit models.Unit

	// Получаем данные из тела запроса
	if err := c.ShouldBindJSON(&unit); err != nil {
		c.JSON(400, gin.H{"error": "Некорректные данные"})
		return
	}

	// Ищем запись по ID
	var existingUnit models.Unit
	if err := database.DB.First(&existingUnit, id).Error; err != nil {
		log.Printf("Ошибка при получении единицы измерения с ID %s: %v", id, err)
		c.JSON(500, gin.H{"error": "Не удалось найти запись для обновления"})
		return
	}

	// Обновляем поля записи
	if err := database.DB.Model(&existingUnit).Updates(unit).Error; err != nil {
		log.Printf("Ошибка при обновлении единицы измерения с ID %s: %v", id, err)
		c.JSON(500, gin.H{"error": "Не удалось обновить запись"})
		return
	}

	// Возвращаем успешный ответ
	c.JSON(200, gin.H{"success": true})
}

// Для работы с внешними ключами

func GetUnitsList(c *gin.Context) {
	var units []models.Unit
	if err := database.DB.Find(&units).Error; err != nil {
		log.Printf("Ошибка при получении единиц измерения: %v", err)
		c.JSON(500, gin.H{"success": false, "error": "Не удалось получить единицы измерения"})
		return
	}

	c.JSON(200, gin.H{"success": true, "units": units})
}
