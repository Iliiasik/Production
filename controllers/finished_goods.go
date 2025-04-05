package controllers

import (
	"github.com/gin-gonic/gin"
	"log"
	"production/database"
	"production/models"
)

func ListFinishedGoods(c *gin.Context) {
	var finishedGoods []models.FinishedGood
	if err := database.DB.Find(&finishedGoods).Error; err != nil {
		log.Printf("Ошибка при получении готовой продукции: %v", err)
		c.JSON(500, gin.H{"error": "Не удалось получить готовую продукцию"})
		return
	}

	var units []models.Unit
	if err := database.DB.Find(&units).Error; err != nil {
		log.Printf("Ошибка при получении единиц измерения: %v", err)
		c.JSON(500, gin.H{"error": "Не удалось получить единицы измерения"})
		return
	}

	unitMap := make(map[uint]string)
	for _, unit := range units {
		unitMap[unit.ID] = unit.Name
	}

	type FinishedGoodsWithUnit struct {
		ID          uint
		Name        string
		Unit        string
		Quantity    float64
		TotalAmount float64
	}

	var finishedGoodsWithUnits []FinishedGoodsWithUnit
	for _, good := range finishedGoods {
		finishedGoodsWithUnits = append(finishedGoodsWithUnits, FinishedGoodsWithUnit{
			ID:          good.ID,
			Name:        good.Name,
			Unit:        unitMap[good.UnitID],
			Quantity:    good.Quantity,
			TotalAmount: good.TotalAmount,
		})
	}

	c.HTML(200, "finished-goods.html", gin.H{
		"finishedGoods": finishedGoodsWithUnits,
	})
}

func DeleteFinishedGood(c *gin.Context) {
	id := c.Param("id")
	log.Printf("Deleting finished good with ID: %s", id)

	if err := database.DB.Delete(&models.FinishedGood{}, id).Error; err != nil {
		log.Printf("Ошибка при удалении записи с ID %s: %v", id, err)
		c.JSON(500, gin.H{"error": "Не удалось удалить запись"})
		return
	}

	log.Printf("Finished good with ID %s deleted successfully", id)
	c.JSON(200, gin.H{"success": true})
}

func AddFinishedGood(c *gin.Context) {
	var finishedGood models.FinishedGood
	if err := c.ShouldBindJSON(&finishedGood); err != nil {
		c.JSON(400, gin.H{"error": "Некорректные данные"})
		return
	}

	if err := database.DB.Create(&finishedGood).Error; err != nil {
		c.JSON(500, gin.H{"error": "Не удалось добавить запись"})
		return
	}

	c.JSON(200, gin.H{"success": true})
}

func GetFinishedGood(c *gin.Context) {
	id := c.Param("id")
	var finishedGood models.FinishedGood

	if err := database.DB.First(&finishedGood, id).Error; err != nil {
		log.Printf("Ошибка при получении готовой продукции с ID %s: %v", id, err)
		c.JSON(500, gin.H{"error": "Не удалось получить запись"})
		return
	}

	c.JSON(200, gin.H{
		"success":      true,
		"finishedGood": finishedGood,
	})
}

func UpdateFinishedGood(c *gin.Context) {
	id := c.Param("id")
	var finishedGood models.FinishedGood

	// Получаем данные из тела запроса
	if err := c.ShouldBindJSON(&finishedGood); err != nil {
		c.JSON(400, gin.H{"error": "Некорректные данные"})
		return
	}

	// Ищем запись по ID
	var existingFinishedGood models.FinishedGood
	if err := database.DB.First(&existingFinishedGood, id).Error; err != nil {
		log.Printf("Ошибка при получении готовой продукции с ID %s: %v", id, err)
		c.JSON(404, gin.H{"error": "Не удалось найти запись для обновления"})
		return
	}

	// Логируем полученные данные
	log.Printf("Updating finished good ID=%s: Name=%s, UnitID=%d, Quantity=%f, TotalAmount=%f",
		id, finishedGood.Name, finishedGood.UnitID, finishedGood.Quantity, finishedGood.TotalAmount)

	// Создаем мапу с обновляемыми значениями
	updateData := map[string]interface{}{
		"name":         finishedGood.Name,
		"unit_id":      finishedGood.UnitID,
		"quantity":     finishedGood.Quantity,
		"total_amount": finishedGood.TotalAmount,
	}

	// Выполняем обновление через мапу
	if err := database.DB.Model(&existingFinishedGood).Updates(updateData).Error; err != nil {
		log.Printf("Ошибка при обновлении готовой продукции с ID %s: %v", id, err)
		c.JSON(500, gin.H{"error": "Не удалось обновить запись"})
		return
	}

	// Возвращаем успешный ответ
	c.JSON(200, gin.H{"success": true})
}
func GetFinishedGoodsList(c *gin.Context) {
	var finishedgoods []models.FinishedGood
	if err := database.DB.Find(&finishedgoods).Error; err != nil {
		log.Printf("Ошибка при получении продукции: %v", err)
		c.JSON(500, gin.H{"success": false, "error": "Не удалось получить продукцию"})
		return
	}

	c.JSON(200, gin.H{"success": true, "finished_goods": finishedgoods})
}
