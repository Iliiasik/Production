package controllers

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"production/database"
	"production/models"
	"strconv"
)

func HomePage(c *gin.Context) {
	tokenString, err := c.Cookie("token")
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/")
		return
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})
	if err != nil || !token.Valid {
		c.Redirect(http.StatusSeeOther, "/")
		return
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		c.Redirect(http.StatusSeeOther, "/")
		return
	}

	var emp models.Employee
	if err := database.DB.Preload("Position").First(&emp, claims.EmployeeID).Error; err != nil {
		c.Redirect(http.StatusSeeOther, "/")
		return
	}

	formattedSalary := fmt.Sprintf("KGS %s", humanize.FormatFloat("# ###.##", emp.Salary))

	c.HTML(http.StatusOK, "home.html", gin.H{
		"FullName": emp.FullName,
		"Position": emp.Position.Name,
		"Salary":   formattedSalary,
		"Username": emp.Username,
	})
}

func GetUserPermissions(c *gin.Context) {
	tokenString, err := c.Cookie("token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var emp models.Employee
	if err := database.DB.First(&emp, claims.EmployeeID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load employee"})
		return
	}

	var rolePermissions []models.Permission
	err = database.DB.
		Joins("JOIN position_permissions ON position_permissions.permission_id = permissions.id").
		Where("position_permissions.position_id = ?", emp.PositionID).
		Find(&rolePermissions).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load role permissions"})
		return
	}

	var userPermissions []models.UserPermission
	err = database.DB.
		Where("employee_id = ?", claims.EmployeeID).
		Find(&userPermissions).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load user permissions"})
		return
	}

	permissionsMap := map[uint]models.Permission{}
	for _, p := range rolePermissions {
		permissionsMap[p.ID] = p
	}
	for _, up := range userPermissions {
		var perm models.Permission
		if err := database.DB.Where("id = ? AND visible_to_user = true", up.PermissionID).First(&perm).Error; err == nil {
			permissionsMap[up.PermissionID] = perm
		} else {
			delete(permissionsMap, up.PermissionID)
		}
	}

	var finalPermissions []models.Permission
	for _, p := range permissionsMap {
		if p.VisibleToUser {
			finalPermissions = append(finalPermissions, p)
		}
	}

	c.JSON(http.StatusOK, gin.H{"permissions": finalPermissions})
}

func AllPositions(c *gin.Context) {
	var positions []models.Position

	if err := database.DB.Find(&positions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения списка ролей"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"positions": positions})
}

func AllUsers(c *gin.Context) {
	var users []models.Employee

	// Загружаем всех сотрудников из базы данных
	if err := database.DB.Find(&users).Error; err != nil {
		// Если произошла ошибка, возвращаем ответ с кодом 500 и сообщением об ошибке
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка загрузки списка пользователей"})
		return
	}

	// Если всё прошло успешно, возвращаем список пользователей
	c.JSON(http.StatusOK, gin.H{"users": users})
}

func AllPermissions(c *gin.Context) {
	var permissions []models.Permission

	if err := database.DB.
		Where("visible_to_user = true").
		Find(&permissions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения разрешений"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"permissions": permissions})
}

func GetPositionPermissions(c *gin.Context) {
	positionID := c.Param("id")
	var permissions []models.Permission

	if err := database.DB.Model(&models.PositionPermission{}).
		Select("permissions.*").
		Joins("JOIN permissions ON permissions.id = position_permissions.permission_id").
		Where("position_permissions.position_id = ? AND permissions.visible_to_user = true", positionID).
		Find(&permissions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения разрешений роли"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"permissions": permissions})
}

func GetEmployeePermissionsByID(c *gin.Context) {
	employeeID := c.Param("id")

	var rolePermissions []models.Permission
	if err := database.DB.
		Joins("JOIN position_permissions ON position_permissions.permission_id = permissions.id").
		Joins("JOIN employees ON employees.position_id = position_permissions.position_id").
		Where("employees.id = ? AND permissions.visible_to_user = true", employeeID).
		Find(&rolePermissions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка загрузки прав роли"})
		return
	}

	var userPermissions []models.UserPermission
	if err := database.DB.
		Where("employee_id = ?", employeeID).
		Find(&userPermissions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка загрузки пользовательских прав"})
		return
	}

	permissionsMap := map[uint]models.Permission{}
	for _, p := range rolePermissions {
		permissionsMap[p.ID] = p
	}

	for _, up := range userPermissions {
		var perm models.Permission
		if err := database.DB.Where("id = ? AND visible_to_user = true", up.PermissionID).First(&perm).Error; err == nil {
			permissionsMap[up.PermissionID] = perm
		}
	}

	var finalPermissions []models.Permission
	for _, p := range permissionsMap {
		finalPermissions = append(finalPermissions, p)
	}

	var allPermissions []models.Permission
	if err := database.DB.Where("visible_to_user = true").Find(&allPermissions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка загрузки всех разрешений"})
		return
	}

	// Важно вернуть отдельно все разрешения и idшники разрешённых
	var grantedIDs []uint
	for _, p := range finalPermissions {
		grantedIDs = append(grantedIDs, p.ID)
	}

	c.JSON(http.StatusOK, gin.H{
		"allPermissions":  allPermissions,
		"userPermissions": grantedIDs,
	})
}

var permissionDependencies = map[string][]string{
	"/units/delete/:id": {"/units/get/:id", "/units", "/units/list"},
	"/units/add":        {"/units/get/:id", "/units", "/units/list"},
	"/units/edit/:id":   {"/units/get/:id", "/units", "/units/list"},

	"/raw-materials/delete/:id": {"/raw-materials", "/units/list", "/raw-materials/get/:id", "/raw-materials/list"},
	"/raw-materials/add":        {"/raw-materials", "/units/list", "/raw-materials/get/:id", "/raw-materials/list"},
	"/raw-materials/edit/:id":   {"/raw-materials", "/units/list", "/raw-materials/get/:id", "/raw-materials/list"},

	"/finished-goods/delete/:id": {"/finished-goods", "/raw-materials/list", "/finished-goods/get/:id", "/finished-goods/list"},
	"/finished-goods/add":        {"/finished-goods", "/raw-materials/list", "/finished-goods/get/:id", "/finished-goods/list"},
	"/finished-goods/edit/:id":   {"/finished-goods", "/raw-materials/list", "/finished-goods/get/:id", "/finished-goods/list"},

	"/ingredients/delete/:id": {"/ingredients", "/ingredients/get/:id", "/ingredients/:product_id", "/ingredients/used-raw-materials/:product_id", "/ingredients/list"},
	"/ingredients/add":        {"/ingredients", "/ingredients/get/:id", "/ingredients/:product_id", "/ingredients/used-raw-materials/:product_id", "/ingredients/list"},
	"/ingredients/edit/:id":   {"/ingredients", "/ingredients/get/:id", "/ingredients/:product_id", "/ingredients/used-raw-materials/:product_id", "/ingredients/list"},

	"/purchases/delete/:id": {"/raw-material-purchases", "/raw-materials/list", "/employees/list"},
	"/purchases/add":        {"/budget/get", "/raw-material-purchases", "/raw-materials/list", "/employees/list"},

	"/production/produce/:product_id": {"/production", "/finished-goods/list", "/employees/list", "/ingredients/list"},

	"/sales/add": {"/sales", "/finished-goods/list", "/employees/list", "/markup/get"},

	"/budget":            {"/budget/get-row/:id", "/budget/get", "/markup/get"},
	"/budget/update/:id": {"/budget"},

	"/credits/add":          {"/credits"},
	"/credits/:id/payments": {"/credits"},
	"/credits/pay/:id":      {"/credits/:id/payments"},

	"/employees/add":        {"/employees", "/employees/list", "/positions/list", "/employees/get/:id"},
	"/employees/edit/:id":   {"/employees", "/employees/list", "/positions/list", "/employees/get/:id"},
	"/employees/delete/:id": {"/employees", "/employees/list", "/positions/list", "/employees/get/:id"},

	"/positions/add":        {"/positions", "/positions/get/:id"},
	"/positions/edit/:id":   {"/positions", "/positions/get/:id"},
	"/positions/delete/:id": {"/positions", "/positions/get/:id"},

	"/salaries":                  {"/salaries/:year/:month", "/salaries/calculate/:year/:month", "/salaries/total-unpaid/:year/:month"},
	"/salaries/pay/:year/:month": {"/salaries"},
	"/salaries/edit/:id":         {"/salaries"},
}

type UpdatePermissionsRequest struct {
	PermissionIDs []uint `json:"permission_ids"`
}

func UpdatePositionPermissions(c *gin.Context) {
	positionID := c.Param("id")
	var req UpdatePermissionsRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат запроса"})
		return
	}

	tx := database.DB.Begin()

	// Получаем имена разрешений по их ID
	var permissionNames []string
	var permissions []models.Permission
	if err := tx.Where("id IN ?", req.PermissionIDs).Find(&permissions).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения разрешений"})
		return
	}

	for _, p := range permissions {
		permissionNames = append(permissionNames, p.Name)
	}

	// Удаляем старые разрешения
	if err := tx.Where("position_id = ?", positionID).Delete(&models.PositionPermission{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка удаления старых разрешений"})
		return
	}

	permissionSet := make(map[string]struct{})

	// Обрабатываем зависимости на основе имен
	for _, permName := range permissionNames {
		permissionSet[permName] = struct{}{}
		if deps, ok := permissionDependencies[permName]; ok {
			for _, dep := range deps {
				permissionSet[dep] = struct{}{}
			}
		}
	}

	// Добавляем разрешения
	for permissionName := range permissionSet {
		var permission models.Permission
		if err := tx.Where("name = ?", permissionName).First(&permission).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка поиска разрешения: " + permissionName})
			return
		}

		newPermission := models.PositionPermission{
			PositionID:   parseUint(positionID),
			PermissionID: permission.ID,
		}
		if err := tx.Create(&newPermission).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка добавления разрешений"})
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка фиксации транзакции"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Разрешения роли успешно обновлены"})
}
func UpdateUserPermissions(c *gin.Context) {
	employeeID := c.Param("id")
	var req UpdatePermissionsRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат запроса"})
		return
	}

	tx := database.DB.Begin()

	// Получаем имена разрешений по их ID
	var permissionNames []string
	var permissions []models.Permission
	if err := tx.Where("id IN ?", req.PermissionIDs).Find(&permissions).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения разрешений"})
		return
	}

	for _, p := range permissions {
		permissionNames = append(permissionNames, p.Name)
	}

	// Удаляем старые разрешения
	if err := tx.Where("employee_id = ?", employeeID).Delete(&models.UserPermission{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка удаления старых разрешений"})
		return
	}

	permissionSet := make(map[string]struct{})

	// Обрабатываем зависимости на основе имен
	for _, permName := range permissionNames {
		permissionSet[permName] = struct{}{}
		if deps, ok := permissionDependencies[permName]; ok {
			for _, dep := range deps {
				permissionSet[dep] = struct{}{}
			}
		}
	}

	// Добавляем разрешения
	for permissionName := range permissionSet {
		var permission models.Permission
		if err := tx.Where("name = ?", permissionName).First(&permission).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка поиска разрешения: " + permissionName})
			return
		}

		newPermission := models.UserPermission{
			EmployeeID:   parseUint(employeeID),
			PermissionID: permission.ID,
		}
		if err := tx.Create(&newPermission).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка добавления разрешений"})
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка фиксации транзакции"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Разрешения сотрудника успешно обновлены"})
}

func parseUint(s string) uint {
	id, _ := strconv.ParseUint(s, 10, 64)
	return uint(id)
}
func GetRoleByUserID(c *gin.Context) {
	// Получаем ID пользователя из параметров маршрута
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный ID"})
		return
	}

	// Получаем роль пользователя из базы данных
	var employee models.Employee
	if err := database.DB.Preload("Position").First(&employee, userID).Error; err != nil {
		log.Printf("Ошибка при получении сотрудника: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "не удалось найти сотрудника"})
		return
	}

	// Проверка на наличие должности у сотрудника
	if employee.PositionID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "сотрудник не имеет должности"})
		return
	}

	// Отправляем роль сотрудника
	c.JSON(http.StatusOK, gin.H{
		"id": employee.Position.ID,
	})
}
