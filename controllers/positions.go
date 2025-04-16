package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"production/database"
	"production/models"
)

func ListPositions(c *gin.Context) {
	// Fetch все записи
	var positions []models.Position
	if err := database.DB.Find(&positions).Error; err != nil {
		log.Printf("Ошибка при получении списка должностей: %v", err)
		c.JSON(500, gin.H{"error": "Не удалось получить список должностей"})
		return
	}

	// Передача данных в шаблон
	c.HTML(200, "positions.html", gin.H{
		"positions": positions,
	})
}

func GetPosition(c *gin.Context) {
	positionID := c.Param("id")
	var position models.Position

	// Ищем запись в базе данных
	if err := database.DB.First(&position, positionID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "Должность не найдена"})
		return
	}

	// Возвращаем данные о должности
	c.JSON(http.StatusOK, gin.H{"success": true, "position": position})
}

func DeletePosition(c *gin.Context) {
	positionID := c.Param("id")

	// Удаляем запись из базы данных
	if err := database.DB.Delete(&models.Position{}, positionID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Ошибка при удалении должности"})
		return
	}

	// Возвращаем успешный ответ
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func EditPosition(c *gin.Context) {
	positionID := c.Param("id")
	var updatedPosition models.Position

	// Парсим JSON из тела запроса
	if err := json.NewDecoder(c.Request.Body).Decode(&updatedPosition); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Неверные данные"})
		return
	}

	// Проверяем, что название не пустое
	if updatedPosition.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Название не может быть пустым"})
		return
	}

	// Обновляем запись в базе данных
	if err := database.DB.Model(&models.Position{}).Where("id = ?", positionID).Updates(updatedPosition).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Ошибка при редактировании должности"})
		return
	}

	// Возвращаем успешный ответ
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func AddPosition(c *gin.Context) {
	var position models.Position

	// Парсим JSON из тела запроса
	if err := json.NewDecoder(c.Request.Body).Decode(&position); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Неверные данные"})
		return
	}

	// Проверяем, что название не пустое
	if position.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Название не может быть пустым"})
		return
	}

	// Добавляем запись в базу данных
	if err := database.DB.Create(&position).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Ошибка при добавлении должности"})
		return
	}

	// Возвращаем успешный ответ
	c.JSON(http.StatusOK, gin.H{"success": true})
}
