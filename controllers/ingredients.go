package controllers

import (
	"github.com/gin-gonic/gin"
	"log"
	"production/database"
	"production/models"
)

func ListIngredients(c *gin.Context) {
	// Получаем список продукции
	var products []models.FinishedGood
	if err := database.DB.Find(&products).Error; err != nil {
		log.Printf("Ошибка при получении продукции: %v", err)
		c.JSON(500, gin.H{"error": "Не удалось получить продукцию"})
		return
	}

	// Передаем только список продукции
	c.HTML(200, "ingredients.html", gin.H{
		"products": products,
	})
}

func DeleteIngredient(c *gin.Context) {
	id := c.Param("id")
	log.Printf("Deleting ingredient with ID: %s", id)

	if err := database.DB.Delete(&models.Ingredient{}, id).Error; err != nil {
		log.Printf("Ошибка при удалении записи с ID %s: %v", id, err)
		c.JSON(500, gin.H{"error": "Не удалось удалить запись"})
		return
	}

	log.Printf("Ingredient with ID %s deleted successfully", id)
	c.JSON(200, gin.H{"success": true})
}

func AddIngredient(c *gin.Context) {
	var ingredient models.Ingredient
	if err := c.ShouldBindJSON(&ingredient); err != nil {
		c.JSON(400, gin.H{"error": "Некорректные данные"})
		return
	}

	if err := database.DB.Create(&ingredient).Error; err != nil {
		c.JSON(500, gin.H{"error": "Не удалось добавить запись"})
		return
	}

	c.JSON(200, gin.H{"success": true})
}

func GetIngredient(c *gin.Context) {
	id := c.Param("id")
	var ingredient models.Ingredient

	if err := database.DB.First(&ingredient, id).Error; err != nil {
		log.Printf("Ошибка при получении ингредиента с ID %s: %v", id, err)
		c.JSON(500, gin.H{"error": "Не удалось получить запись"})
		return
	}

	c.JSON(200, gin.H{
		"success":    true,
		"ingredient": ingredient,
	})
}

func UpdateIngredient(c *gin.Context) {
	id := c.Param("id")
	var ingredient models.Ingredient

	if err := c.ShouldBindJSON(&ingredient); err != nil {
		c.JSON(400, gin.H{"error": "Некорректные данные"})
		return
	}

	var existingIngredient models.Ingredient
	if err := database.DB.First(&existingIngredient, id).Error; err != nil {
		log.Printf("Ошибка при получении ингредиента с ID %s: %v", id, err)
		c.JSON(500, gin.H{"error": "Не удалось найти запись для обновления"})
		return
	}

	if err := database.DB.Model(&existingIngredient).Updates(ingredient).Error; err != nil {
		log.Printf("Ошибка при обновлении записи с ID %s: %v", id, err)
		c.JSON(500, gin.H{"error": "Не удалось обновить запись"})
		return
	}

	c.JSON(200, gin.H{"success": true})
}

func GetIngredientsByProduct(c *gin.Context) {
	productID := c.Param("product_id")
	var ingredients []models.Ingredient

	if err := database.DB.Where("product_id = ?", productID).Find(&ingredients).Error; err != nil {
		log.Printf("Ошибка при получении ингредиентов: %v", err)
		c.JSON(500, gin.H{"error": "Не удалось получить ингредиенты"})
		return
	}

	// Получаем список сырья для подстановки имен
	var materials []models.RawMaterial
	if err := database.DB.Find(&materials).Error; err != nil {
		c.JSON(500, gin.H{"error": "Не удалось получить сырье"})
		return
	}

	materialMap := make(map[uint]string)
	for _, material := range materials {
		materialMap[material.ID] = material.Name
	}

	// Подставляем данные
	type IngredientResponse struct {
		ID       uint    `json:"id"`
		Material string  `json:"material"`
		Quantity float64 `json:"quantity"`
	}

	var response []IngredientResponse
	for _, ing := range ingredients {
		response = append(response, IngredientResponse{
			ID:       ing.ID,
			Material: materialMap[ing.RawMaterialID],
			Quantity: ing.Quantity,
		})
	}

	c.JSON(200, gin.H{"ingredients": response})
}

func GetUsedRawMaterialsByProduct(c *gin.Context) {
	productID := c.Param("product_id")

	// Загружаем ингредиенты для указанного продукта
	var ingredients []models.Ingredient
	if err := database.DB.Where("product_id = ?", productID).Find(&ingredients).Error; err != nil {
		log.Printf("Ошибка при получении ингредиентов: %v", err)
		c.JSON(500, gin.H{"error": "Не удалось получить данные об ингредиентах"})
		return
	}

	// Если ингредиентов нет, возвращаем пустой массив
	if len(ingredients) == 0 {
		c.JSON(200, gin.H{"success": true, "used_raw_materials": []models.RawMaterial{}})
		return
	}

	// Извлекаем ID используемого сырья
	var usedRawMaterialIDs []uint
	for _, ingredient := range ingredients {
		usedRawMaterialIDs = append(usedRawMaterialIDs, ingredient.RawMaterialID)
	}

	// Получаем полные данные о сырье по ID
	var usedRawMaterials []models.RawMaterial
	if len(usedRawMaterialIDs) > 0 {
		if err := database.DB.Where("id IN (?)", usedRawMaterialIDs).Find(&usedRawMaterials).Error; err != nil {
			log.Printf("Ошибка при получении данных о сырье: %v", err)
			c.JSON(500, gin.H{"error": "Не удалось получить данные о сырье"})
			return
		}
	}

	// Формируем ответ
	type RawMaterialResponse struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	}

	var response []RawMaterialResponse
	for _, material := range usedRawMaterials {
		response = append(response, RawMaterialResponse{
			ID:   material.ID,
			Name: material.Name,
		})
	}

	c.JSON(200, gin.H{"success": true, "used_raw_materials": response})
}
