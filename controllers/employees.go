package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"production/database"
	"production/models"
)

func GetEmployeesList(c *gin.Context) {
	var employees []models.Employee
	if err := database.DB.Find(&employees).Error; err != nil {
		log.Printf("Ошибка при получении сотрудников: %v", err)
		c.JSON(500, gin.H{"success": false, "error": "Не удалось получить сотрудников"})
		return
	}

	c.JSON(200, gin.H{"success": true, "employees": employees})
}

func ListEmployees(c *gin.Context) {
	var employees []models.Employee
	if err := database.DB.Find(&employees).Error; err != nil {
		log.Printf("Ошибка при получении списка сотрудников: %v", err)
		c.JSON(500, gin.H{"error": "Не удалось получить список сотрудников"})
		return
	}
	// Получаем список должностей
	var positions []models.Position
	if err := database.DB.Find(&positions).Error; err != nil {
		log.Printf("Ошибка при получении списка должностей: %v", err)
		c.JSON(500, gin.H{"error": "Не удалось получить список должностей"})
		return
	}
	// Создаем мапу
	positionMap := make(map[uint]string)
	for _, position := range positions {
		positionMap[position.ID] = position.Name
	}

	type EmployeesWithDetails struct {
		ID       uint
		FullName string
		Position string
		Salary   float64
		Address  string
		Phone    string
	}

	var employeesWithDetails []EmployeesWithDetails
	for _, emp := range employees {
		employeesWithDetails = append(employeesWithDetails, EmployeesWithDetails{
			ID:       emp.ID,
			FullName: emp.FullName,
			Position: positionMap[emp.PositionID],
			Salary:   emp.Salary,
			Address:  emp.Address,
			Phone:    emp.Phone,
		})
	}

	// Передаем данные в шаблон
	c.HTML(200, "employees.html", gin.H{
		"employees": employeesWithDetails,
	})
}

func GetEmployee(c *gin.Context) {
	id := c.Param("id")
	var employee models.Employee

	if err := database.DB.Preload("Position").First(&employee, id).Error; err != nil {
		log.Printf("Ошибка при получении сотрудника с ID %s: %v", id, err)
		c.JSON(500, gin.H{"error": "Не удалось получить запись"})
		return
	}

	c.JSON(200, gin.H{
		"success":  true,
		"employee": employee,
	})
}
func UpdateEmployee(c *gin.Context) {
	id := c.Param("id")
	var employee models.Employee

	if err := c.ShouldBindJSON(&employee); err != nil {
		c.JSON(400, gin.H{"error": "Некорректные данные"})
		return
	}

	var existingEmployee models.Employee
	if err := database.DB.First(&existingEmployee, id).Error; err != nil {
		log.Printf("Ошибка при получении сотрудника с ID %s: %v", id, err)
		c.JSON(404, gin.H{"error": "Не удалось найти запись для обновления"})
		return
	}

	log.Printf("Обновление сотрудника ID=%s: FullName=%s, PositionID=%d, Salary=%f, Address=%s, Phone=%s",
		id, employee.FullName, employee.PositionID, employee.Salary, employee.Address, employee.Phone)

	updateData := map[string]interface{}{
		"full_name":   employee.FullName,
		"position_id": employee.PositionID,
		"salary":      employee.Salary,
		"address":     employee.Address,
		"phone":       employee.Phone,
	}

	if err := database.DB.Model(&existingEmployee).Updates(updateData).Error; err != nil {
		log.Printf("Ошибка при обновлении сотрудника с ID %s: %v", id, err)
		c.JSON(500, gin.H{"error": "Не удалось обновить запись"})
		return
	}

	c.JSON(200, gin.H{"success": true})
}
func DeleteEmployee(c *gin.Context) {
	id := c.Param("id")
	log.Printf("Удаление сотрудника с ID: %s", id)

	if err := database.DB.Delete(&models.Employee{}, id).Error; err != nil {
		log.Printf("Ошибка при удалении сотрудника с ID %s: %v", id, err)
		c.JSON(500, gin.H{"error": "Не удалось удалить запись"})
		return
	}

	log.Printf("Сотрудник с ID %s успешно удален", id)
	c.JSON(200, gin.H{"success": true})
}
func GetAllPositions(c *gin.Context) {
	var positions []models.Position
	if err := database.DB.Find(&positions).Error; err != nil {
		log.Printf("Ошибка при получении позиций: %v", err)
		c.JSON(500, gin.H{"error": "Не удалось получить список позиций"})
		return
	}

	c.JSON(200, gin.H{"success": true, "positions": positions})
}
func AddEmployee(c *gin.Context) {
	var input struct {
		FullName   string  `json:"full_name"`
		Username   string  `json:"username"`
		PositionID uint    `json:"position_id"`
		Salary     float64 `json:"salary"`
		Address    string  `json:"address"`
		Phone      string  `json:"phone"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Некорректные данные"})
		return
	}

	passwordHash, err := HashPassword(input.Username)
	if err != nil {
		c.JSON(500, gin.H{"error": "Ошибка при создании пароля"})
		return
	}

	employee := models.Employee{
		FullName:     input.FullName,
		Username:     input.Username,
		PasswordHash: passwordHash,
		PositionID:   input.PositionID,
		Salary:       input.Salary,
		Address:      input.Address,
		Phone:        input.Phone,
	}

	if err := database.DB.Create(&employee).Error; err != nil {
		c.JSON(500, gin.H{"error": "Не удалось добавить запись"})
		return
	}

	c.JSON(200, gin.H{"success": true})
}
func GetNextUsername(c *gin.Context) {
	var count int64
	if err := database.DB.Model(&models.Employee{}).Count(&count).Error; err != nil {
		c.JSON(500, gin.H{"error": "Ошибка при подсчёте пользователей"})
		return
	}

	username := fmt.Sprintf("emp%03d", count+1)
	c.JSON(200, gin.H{"username": username})
}
