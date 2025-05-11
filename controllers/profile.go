package controllers

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
	"log"
	"net/http"
	"production/database"
	"production/models"
	"strconv"
)

func HomePage(c *gin.Context) {
	// Пытаемся получить JWT-токен из cookie с именем "token".
	tokenString, err := c.Cookie("token")
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/") // Если не нашли, редирект на страницу входа
		return
	}
	// jwt.ParseWithClaims разбирает токен и проверяет его подпись
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})
	if err != nil || !token.Valid {
		c.Redirect(http.StatusSeeOther, "/")
		return
	}
	// Claims — это структура с полями из payload токена (например, EmployeeID, Username, exp, iat).
	claims, ok := token.Claims.(*Claims)
	if !ok {
		c.Redirect(http.StatusSeeOther, "/")
		return
	}
	// По EmployeeID, полученному из токена, загружается сотрудник из БД.
	//Preload("Position") — это JOIN с таблицей должностей, чтобы можно было обратиться к emp.Position.Name.
	var emp models.Employee
	if err := database.DB.Preload("Position").First(&emp, claims.EmployeeID).Error; err != nil {
		c.Redirect(http.StatusSeeOther, "/")
		return
	}
	// Форматируем денежное поле
	formattedSalary := fmt.Sprintf("KGS %s", humanize.FormatFloat("# ###.##", emp.Salary))

	c.HTML(http.StatusOK, "home.html", gin.H{
		"FullName":          emp.FullName,
		"Position":          emp.Position.Name,
		"Salary":            formattedSalary,
		"Username":          emp.Username,
		"IsPasswordChanged": emp.IsPasswordChanged,
	})
}

func GetUserPermissions(c *gin.Context) {
	tokenString, err := c.Cookie("token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"}) // 401 - неавторизован
		return
	}
	// Парсинг и валидация токена

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
	// Получаем запись сотрудника из базы по ID, чтобы узнать его PositionID
	var emp models.Employee
	if err := database.DB.First(&emp, claims.EmployeeID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load employee"})
		return
	}
	var rolePermissions []models.Permission
	err = database.DB. // получаем все права (permissions), связанные с определённой должностью (position_id)
				Joins("JOIN position_permissions ON position_permissions.permission_id = permissions.id").
				Where("position_permissions.position_id = ?", emp.PositionID).
				Find(&rolePermissions).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load role permissions"})
		return
	}
	// Загружаем персональные права, выданные именно этому сотруднику
	var userPermissions []models.UserPermission
	err = database.DB.
		Where("employee_id = ?", claims.EmployeeID). // claims - данные которые мы получили из JWT
		Find(&userPermissions).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load user permissions"})
		return
	}
	// Используем map[uint]Permission{} как структуру с уникальными правами (без дублирования по ID)
	permissionsMap := map[uint]models.Permission{}
	for _, p := range rolePermissions {
		permissionsMap[p.ID] = p // Сначала добавляются права по роли
	}
	for _, up := range userPermissions { // Затем добавляются права пользователя (если visible_to_user = true)
		var perm models.Permission
		if err := database.DB.Where("id = ? AND visible_to_user = true", up.PermissionID).First(&perm).Error; err == nil {
			permissionsMap[up.PermissionID] = perm
		} else {
			delete(permissionsMap, up.PermissionID)
		}
	}
	// Создаётся финальный список только тех прав, которые можно показывать пользователю (VisibleToUser = true)
	var finalPermissions []models.Permission
	for _, p := range permissionsMap {
		if p.VisibleToUser {
			finalPermissions = append(finalPermissions, p)
		}
	}
	// передаем данные
	c.JSON(http.StatusOK, gin.H{"permissions": finalPermissions})
}

func AllPositions(c *gin.Context) {
	var positions []models.Position

	if err := database.DB.Find(&positions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения списка должностей"})
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
	// Model(&models.PositionPermission{})
	//Говорит GORM, что основной запрос начинается от таблицы position_permissions (связка должностей и прав).
	//
	//Select("permissions.*")
	//Мы хотим получить все поля из таблицы permissions.
	//
	//Joins(...)
	//Объединяем position_permissions с таблицей permissions по полю permission_id.
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
	// Получает все права, которые назначены должности сотрудника (через таблицу position_permissions).
	//
	//Только если permissions.visible_to_user = true.
	//
	//Используется двойной JOIN:
	//
	//permissions ↔ position_permissions ↔ employees
	var rolePermissions []models.Permission
	if err := database.DB.
		Joins("JOIN position_permissions ON position_permissions.permission_id = permissions.id").
		Joins("JOIN employees ON employees.position_id = position_permissions.position_id").
		Where("employees.id = ? AND permissions.visible_to_user = true", employeeID).
		Find(&rolePermissions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка загрузки прав роли"})
		return
	}
	// Получаем индивидуальные (персональные) права, которые назначены только этому сотруднику (не через должность)
	var userPermissions []models.UserPermission
	if err := database.DB.
		Where("employee_id = ?", employeeID).
		Find(&userPermissions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка загрузки пользовательских прав"})
		return
	}
	// Используем map[uint]Permission, чтобы избежать повторов.
	//
	//Добавляем сначала ролевые, потом пользовательские, перезаписывая, если нужно.
	//
	//Каждое право проверяется на visible_to_user = true.
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
	// Формируем итоговый массив прав пользователя
	var finalPermissions []models.Permission
	for _, p := range permissionsMap {
		finalPermissions = append(finalPermissions, p)
	}

	var allPermissions []models.Permission // Загружаются все разрешения, которые можно показать пользователю
	if err := database.DB.Where("visible_to_user = true").Find(&allPermissions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка загрузки всех разрешений"})
		return
	}

	// Важно вернуть отдельно все разрешения и idшники разрешённых
	// Получение ID тех прав, которые уже назначены
	var grantedIDs []uint
	for _, p := range finalPermissions {
		grantedIDs = append(grantedIDs, p.ID)
	}

	c.JSON(http.StatusOK, gin.H{
		"allPermissions":  allPermissions,
		"userPermissions": grantedIDs,
	})
}

// Словарь зависимостей
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

	"/employees":            {"/employees/list"},
	"/employees/add":        {"/employees", "/positions/list", "/employees/get/:id"},
	"/employees/edit/:id":   {"/employees", "/positions/list", "/employees/get/:id"},
	"/employees/delete/:id": {"/employees", "/positions/list", "/employees/get/:id"},

	"/positions/add":        {"/positions", "/positions/get/:id"},
	"/positions/edit/:id":   {"/positions", "/positions/get/:id"},
	"/positions/delete/:id": {"/positions", "/positions/get/:id"},

	"/salaries":                  {"/salaries/:year/:month", "/salaries/calculate/:year/:month", "/salaries/total-unpaid/:year/:month"},
	"/salaries/pay/:year/:month": {"/salaries"},
	"/salaries/edit/:id":         {"/salaries"},

	"/reports/sales":       {"/reports"},
	"/reports/productions": {"/reports"},
	"/reports/salaries":    {"/reports"},
	"/reports/purchases":   {"/reports"},
	"/reports/payments":    {"/reports"},
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

	if err := updatePositionPermissionsTransaction(positionID, req.PermissionIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Разрешения роли успешно обновлены"})
}

func updatePositionPermissionsTransaction(positionID string, permissionIDs []uint) error {
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	permissionNames, err := getPermissionNames(tx, permissionIDs)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("Ошибка получения разрешений: %v", err)
	}

	if err := tx.Where("position_id = ?", positionID).Delete(&models.PositionPermission{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("Ошибка удаления старых разрешений: %v", err)
	}

	permissionSet := buildPermissionSet(permissionNames)

	if err := addNewPositionPermissions(tx, parseUint(positionID), permissionSet); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func addNewPositionPermissions(tx *gorm.DB, positionID uint, permissionSet map[string]struct{}) error {
	for permissionName := range permissionSet {
		var permission models.Permission
		if err := tx.Where("name = ?", permissionName).First(&permission).Error; err != nil {
			return fmt.Errorf("Ошибка поиска разрешения: %s: %v", permissionName, err)
		}

		if err := tx.Create(&models.PositionPermission{
			PositionID:   positionID,
			PermissionID: permission.ID,
		}).Error; err != nil {
			return fmt.Errorf("Ошибка добавления разрешений: %v", err)
		}
	}
	return nil
}

func UpdateUserPermissions(c *gin.Context) {
	employeeID := c.Param("id")
	var req UpdatePermissionsRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат запроса"})
		return
	}

	if err := updateUserPermissionsTransaction(employeeID, req.PermissionIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Разрешения сотрудника успешно обновлены"})
}

func updateUserPermissionsTransaction(employeeID string, permissionIDs []uint) error {
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	permissionNames, err := getPermissionNames(tx, permissionIDs)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("Ошибка получения разрешений: %v", err)
	}

	if err := tx.Where("employee_id = ?", employeeID).Delete(&models.UserPermission{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("Ошибка удаления старых разрешений: %v", err)
	}

	permissionSet := buildPermissionSet(permissionNames)

	if err := addNewPermissions(tx, parseUint(employeeID), permissionSet); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func getPermissionNames(tx *gorm.DB, permissionIDs []uint) ([]string, error) {
	var permissions []models.Permission
	if err := tx.Where("id IN ?", permissionIDs).Find(&permissions).Error; err != nil {
		return nil, err
	}

	names := make([]string, len(permissions))
	for i, p := range permissions {
		names[i] = p.Name
	}
	return names, nil
}

func buildPermissionSet(permissionNames []string) map[string]struct{} {
	permissionSet := make(map[string]struct{})
	for _, permName := range permissionNames {
		permissionSet[permName] = struct{}{}
		if deps, ok := permissionDependencies[permName]; ok {
			for _, dep := range deps {
				permissionSet[dep] = struct{}{}
			}
		}
	}
	return permissionSet
}

func addNewPermissions(tx *gorm.DB, employeeID uint, permissionSet map[string]struct{}) error {
	for permissionName := range permissionSet {
		var permission models.Permission
		if err := tx.Where("name = ?", permissionName).First(&permission).Error; err != nil {
			return fmt.Errorf("Ошибка поиска разрешения: %s: %v", permissionName, err)
		}

		if err := tx.Create(&models.UserPermission{
			EmployeeID:   employeeID,
			PermissionID: permission.ID,
		}).Error; err != nil {
			return fmt.Errorf("Ошибка добавления разрешений: %v", err)
		}
	}
	return nil
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
