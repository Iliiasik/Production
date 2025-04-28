package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"production/controllers"
	"production/database"
	"production/models"
)

// Authorize проверяет, есть ли у текущего пользователя разрешение permissionName.
// Токен берётся из cookie "token", в которой сохраняется JWT после входа.
func Authorize(permissionName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1) достаём JWT из cookie
		tokenString, err := c.Cookie("token")
		if err != nil || tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Пожалуйста, авторизуйтесь"})
			return
		}

		// 2) парсим и валидируем токен
		token, err := jwt.ParseWithClaims(tokenString, &controllers.Claims{}, func(t *jwt.Token) (interface{}, error) {
			return controllers.JwtKey, nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Неверный или просроченный токен"})
			return
		}

		// 3) извлекаем employeeID из claims
		claims, ok := token.Claims.(*controllers.Claims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Неверные данные токена"})
			return
		}
		employeeID := claims.EmployeeID

		// 4) проверяем, что такое разрешение вообще есть в таблице permissions
		var perm models.Permission
		if err := database.DB.Where("name = ?", permissionName).First(&perm).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Разрешение не найдено"})
			return
		}

		// 5) сначала смотрим индивидуальные права пользователя
		var userPermissions []models.UserPermission
		if err := database.DB.
			Where("employee_id = ?", employeeID).
			Find(&userPermissions).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Ошибка загрузки прав пользователя"})
			return
		}

		// Обрабатываем индивидуальные разрешения пользователя
		var userHasPermission bool
		for _, up := range userPermissions {
			if up.PermissionID == perm.ID {
				userHasPermission = true
				break
			}
		}

		if userHasPermission {
			c.Next()
			return
		}

		// 6) затем — права по должности
		var emp models.Employee
		if err := database.DB.First(&emp, employeeID).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Пользователь не найден"})
			return
		}

		var positionPermissions []models.PositionPermission
		if err := database.DB.
			Where("position_id = ? AND permission_id = ?", emp.PositionID, perm.ID).
			Find(&positionPermissions).Error; err == nil && len(positionPermissions) > 0 {
			c.Next()
			return
		}

		// 7) если ни по индивидуальным, ни по должности не нашли — запретим
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "У вас нет прав для выполнения действия"})
	}
}
