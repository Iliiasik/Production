package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"production/database"
	"production/models"
	"time"
)

// Получение задач для текущего сотрудника
func GetMyTasks(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		cookie, err := c.Cookie("token")
		if err != nil {
			c.JSON(401, gin.H{"error": "Необходим токен авторизации"})
			return
		}
		tokenString = cookie
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})

	if err != nil || !token.Valid {
		c.JSON(401, gin.H{"error": "Недействительный токен"})
		return
	}

	// Теперь EmployeeID берется из claims
	userID := claims.EmployeeID

	var tasks []models.Task
	if err := database.DB.Where("assigned_to = ?", userID).Find(&tasks).Error; err != nil {
		c.JSON(500, gin.H{"error": "Не удалось получить задачи"})
		return
	}

	c.JSON(200, gin.H{"tasks": tasks})
}

// Назначение задачи директором
func AssignTask(c *gin.Context) {
	var input struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		AssignedTo  uint   `json:"assigned_to"`
		DueDate     string `json:"due_date"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Неверный формат запроса: " + err.Error()})
		return
	}

	// Парсим дату
	dueDate, err := time.Parse("2006-01-02", input.DueDate)
	if err != nil {
		c.JSON(400, gin.H{"error": "Некорректный формат даты. Ожидается YYYY-MM-DD"})
		return
	}

	task := models.Task{
		Title:       input.Title,
		Description: input.Description,
		AssignedTo:  input.AssignedTo,
		DueDate:     dueDate,
		Status:      "Новая",
		CreatedAt:   time.Now(),
	}

	if err := database.DB.Create(&task).Error; err != nil {
		c.JSON(500, gin.H{"error": "Ошибка базы данных: " + err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"task_id": task.ID,
	})
}
func GetTasks(c *gin.Context) {
	var tasks []models.Task

	// Загружаем задачи с предзагрузкой связанных сотрудников
	if err := database.DB.Preload("Employee").Find(&tasks).Error; err != nil {
		c.JSON(500, gin.H{"error": "Не удалось загрузить задачи: " + err.Error()})
		return
	}

	// Упрощенная структура для ответа
	type TaskResponse struct {
		ID          uint   `json:"id"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Status      string `json:"status"`
		Employee    struct {
			FullName string `json:"full_name"`
		} `json:"employee"`
		DueDate time.Time `json:"due_date"`
	}

	var response []TaskResponse
	for _, task := range tasks {
		item := TaskResponse{
			ID:          task.ID,
			Title:       task.Title,
			Description: task.Description, // Добавляем описание
			Status:      task.Status,
			DueDate:     task.DueDate,
		}
		item.Employee.FullName = task.Employee.FullName
		response = append(response, item)
	}

	c.JSON(200, gin.H{
		"success": true,
		"tasks":   response,
	})
}
